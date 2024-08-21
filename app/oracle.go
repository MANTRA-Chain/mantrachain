package app

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	oracleconfig "github.com/skip-mev/slinky/oracle/config"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	oraclepreblock "github.com/skip-mev/slinky/abci/preblock/oracle"
	"github.com/skip-mev/slinky/abci/strategies/aggregator"
	"github.com/skip-mev/slinky/abci/strategies/currencypair"
	"github.com/skip-mev/slinky/abci/ve"
	"github.com/skip-mev/slinky/pkg/math/voteweighted"
	oracleclient "github.com/skip-mev/slinky/service/clients/oracle"
	servicemetrics "github.com/skip-mev/slinky/service/metrics"
	marketmaptypes "github.com/skip-mev/slinky/x/marketmap/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"

	tmtypes "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/skip-mev/slinky/cmd/constants/marketmaps"

	"github.com/skip-mev/slinky/abci/proposals"

	compression "github.com/skip-mev/slinky/abci/strategies/codec"
)

// initializeOracle initializes the oracle client and metrics.
func (app *App) initializeOracle(appOpts types.AppOptions) (oracleclient.OracleClient, servicemetrics.Metrics, error) {
	// Read general config from app-opts, and construct oracle service.
	cfg, err := oracleconfig.ReadConfigFromAppOpts(appOpts)
	if err != nil {
		return nil, nil, err
	}

	// If app level instrumentation is enabled, then wrap the oracle service with a metrics client
	// to get metrics on the oracle service (for ABCI++). This will allow the instrumentation to track
	// latency in VerifyVoteExtension requests and more.
	oracleMetrics, err := servicemetrics.NewMetricsFromConfig(cfg, app.ChainID())
	if err != nil {
		return nil, nil, err
	}

	// Create the oracle service.
	oracleClient, err := oracleclient.NewPriceDaemonClientFromConfig(
		cfg,
		app.Logger().With("client", "oracle"),
		oracleMetrics,
	)
	if err != nil {
		return nil, nil, err
	}

	// Connect to the oracle service (default timeout of 5 seconds).
	go func() {
		app.Logger().Info("attempting to start oracle client...", "address", cfg.OracleAddress)
		if err := oracleClient.Start(context.Background()); err != nil {
			app.Logger().Error("failed to start oracle client", "err", err)
			panic(err)
		}
	}()

	return oracleClient, oracleMetrics, nil
}

func (app *App) initializeABCIExtensions(oracleClient oracleclient.OracleClient, oracleMetrics servicemetrics.Metrics) {
	// Create the proposal handler that will be used to fill proposals with
	// transactions and oracle data.
	proposalHandler := proposals.NewProposalHandler(
		app.Logger(),
		baseapp.NoOpPrepareProposal(),
		baseapp.NoOpProcessProposal(),
		ve.NewDefaultValidateVoteExtensionsFn(app.StakingKeeper),
		compression.NewCompressionVoteExtensionCodec(
			compression.NewDefaultVoteExtensionCodec(),
			compression.NewZLibCompressor(),
		),
		compression.NewCompressionExtendedCommitCodec(
			compression.NewDefaultExtendedCommitCodec(),
			compression.NewZStdCompressor(),
		),
		currencypair.NewDeltaCurrencyPairStrategy(app.OracleKeeper),
		oracleMetrics,
	)
	app.SetPrepareProposal(proposalHandler.PrepareProposalHandler())
	app.SetProcessProposal(proposalHandler.ProcessProposalHandler())

	// Create the aggregation function that will be used to aggregate oracle data
	// from each validator.
	aggregatorFn := voteweighted.MedianFromContext(
		app.Logger(),
		app.StakingKeeper,
		voteweighted.DefaultPowerThreshold,
	)
	veCodec := compression.NewCompressionVoteExtensionCodec(
		compression.NewDefaultVoteExtensionCodec(),
		compression.NewZLibCompressor(),
	)
	ecCodec := compression.NewCompressionExtendedCommitCodec(
		compression.NewDefaultExtendedCommitCodec(),
		compression.NewZStdCompressor(),
	)

	// Create the pre-finalize block hook that will be used to apply oracle data
	// to the state before any transactions are executed (in finalize block).
	oraclePreBlockHandler := oraclepreblock.NewOraclePreBlockHandler(
		app.Logger(),
		aggregatorFn,
		app.OracleKeeper,
		oracleMetrics,
		currencypair.NewDeltaCurrencyPairStrategy(app.OracleKeeper), // IMPORTANT: always construct new currency pair strategy objects when functions require them as arguments.
		veCodec,
		ecCodec,
	)

	app.SetPreBlocker(oraclePreBlockHandler.WrappedPreBlocker(app.ModuleManager))

	// Create the vote extensions handler that will be used to extend and verify
	// vote extensions (i.e. oracle data).
	voteExtensionsHandler := ve.NewVoteExtensionHandler(
		app.Logger(),
		oracleClient,
		time.Second, // timeout
		currencypair.NewDeltaCurrencyPairStrategy(app.OracleKeeper), // IMPORTANT: always construct new currency pair strategy objects when functions require them as arguments.
		veCodec,
		aggregator.NewOraclePriceApplier(
			aggregator.NewDefaultVoteAggregator(
				app.Logger(),
				aggregatorFn,
				// we need a separate price strategy here, so that we can optimistically apply the latest prices
				// and extend our vote based on these prices
				currencypair.NewDeltaCurrencyPairStrategy(app.OracleKeeper), // IMPORTANT: always construct new currency pair strategy objects when functions require them as arguments.
			),
			app.OracleKeeper,
			veCodec,
			ecCodec,
			app.Logger(),
		),
		oracleMetrics,
	)
	app.SetExtendVoteHandler(voteExtensionsHandler.ExtendVoteHandler())
	app.SetVerifyVoteExtensionHandler(voteExtensionsHandler.VerifyVoteExtensionHandler())
}

type AppUpgrade struct {
	Handler      upgradetypes.UpgradeHandler
	StoreUpgrade storetypes.StoreUpgrades
}

// createSlinkyUpgrader returns the upgrade name and an upgrade handler that:
// - runs migrations
// - updates the consensus keeper params with a vote extension enable height. (height of upgrade + 10).
// - adds the core markets to x/marketmap keeper.
// additionally, it returns the required StoreUpgrades needed for the new slinky modules added to this chain.
func createSlinkyUpgrader(app *App) AppUpgrade {
	return AppUpgrade{
		Handler: func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			migrations, err := app.ModuleManager.RunMigrations(ctx, app.Configurator(), fromVM)
			if err != nil {
				return nil, err
			}
			// upgrade consensus params to enable vote extensions
			consensusParams, err := app.ConsensusParamsKeeper.Params(ctx, nil)
			if err != nil {
				return nil, err
			}
			sdkCtx := sdk.UnwrapSDKContext(ctx)
			consensusParams.Params.Abci = &tmtypes.ABCIParams{
				VoteExtensionsEnableHeight: sdkCtx.BlockHeight() + int64(10),
			}
			_, err = app.ConsensusParamsKeeper.UpdateParams(ctx, &consensustypes.MsgUpdateParams{
				Authority: app.ConsensusParamsKeeper.GetAuthority(),
				Block:     consensusParams.Params.Block,
				Evidence:  consensusParams.Params.Evidence,
				Validator: consensusParams.Params.Validator,
				Abci:      consensusParams.Params.Abci,
			})
			if err != nil {
				return nil, err
			}

			// add core markets
			coreMarkets := marketmaps.CoreMarketMap
			markets := coreMarkets.Markets
			keys := make([]string, 0, len(markets))
			for name := range markets {
				keys = append(keys, name)
			}
			slices.Sort(keys)

			// iterates over slice and not map
			for _, marketName := range keys {
				// create market
				market := markets[marketName]
				err = app.MarketMapKeeper.CreateMarket(sdkCtx, market)
				if err != nil {
					return nil, err
				}

				// invoke hooks
				err = app.MarketMapKeeper.Hooks().AfterMarketCreated(sdkCtx, market)
				if err != nil {
					return nil, err
				}
			}

			err = app.MarketMapKeeper.SetParams(
				sdkCtx,
				marketmaptypes.Params{
					Admin: authtypes.NewModuleAddress(govtypes.ModuleName).String(), // governance. allows gov to add or remove market authorities.
					// market authority addresses may add and update markets to the x/marketmap module
					MarketAuthorities: []string{
						authtypes.NewModuleAddress(govtypes.ModuleName).String(), // governance
					}},
			)
			if err != nil {
				return nil, fmt.Errorf("failed to set x/marketmap params: %w", err)
			}
			return migrations, nil
		},
		StoreUpgrade: storetypes.StoreUpgrades{
			Added: []string{
				marketmaptypes.ModuleName,
				oracletypes.ModuleName,
			},
		},
	}
}
