package v4

import (
	"context"
	"fmt"
	"slices"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	guardkeeper "github.com/MANTRA-Finance/mantrachain/x/guard/keeper"
	tmtypes "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	consensuskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channelkeeper "github.com/cosmos/ibc-go/v8/modules/core/04-channel/keeper"
	"github.com/skip-mev/slinky/cmd/constants/marketmaps"
	marketmapkeeper "github.com/skip-mev/slinky/x/marketmap/keeper"
	marketmaptypes "github.com/skip-mev/slinky/x/marketmap/types"
)

// CreateUpgradeHandler creates an SDK upgrade handler for v3.0.0
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	govkeeper govkeeper.Keeper,
	guardkeeper guardkeeper.Keeper,
	channelkeeper channelkeeper.Keeper,
	consensuskeeper consensuskeeper.Keeper,
	marketmapkeeper marketmapkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(goCtx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		migrations, err := mm.RunMigrations(goCtx, configurator, fromVM)
		if err != nil {
			return nil, err
		}

		// set the min deposit for expedited proposals to 10x the value of the current expedited min deposit
		ctx := sdk.UnwrapSDKContext(goCtx)
		govParams, err := govkeeper.Params.Get(ctx)
		if err != nil {
			return nil, err
		}
		expeditedMinDeposit := govParams.GetMinDeposit()
		for i := range expeditedMinDeposit {
			expeditedMinDeposit[i].Amount = expeditedMinDeposit[i].Amount.MulRaw(10)
		}
		govParams.ExpeditedMinDeposit = expeditedMinDeposit
		if err := govkeeper.Params.Set(ctx, govParams); err != nil {
			return nil, err
		}

		// whitelist all ibc escrow accounts
		var escrowAccounts []sdk.AccAddress
		for _, channel := range channelkeeper.GetAllChannels(ctx) {
			escrowAccount := transfertypes.GetEscrowAddress(channel.PortId, channel.ChannelId)
			escrowAccounts = append(escrowAccounts, escrowAccount)
		}
		guardkeeper.AddTransferAccAddressesWhitelist(ctx, escrowAccounts)

		// upgrade consensus params to enable vote extensions
		consensusParams, err := consensuskeeper.Params(ctx, nil)
		if err != nil {
			return nil, err
		}
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		consensusParams.Params.Abci = &tmtypes.ABCIParams{
			VoteExtensionsEnableHeight: sdkCtx.BlockHeight() + int64(10),
		}
		_, err = consensuskeeper.UpdateParams(ctx, &consensustypes.MsgUpdateParams{
			Authority: consensuskeeper.GetAuthority(),
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
			err = marketmapkeeper.CreateMarket(sdkCtx, market)
			if err != nil {
				return nil, err
			}

			// invoke hooks
			err = marketmapkeeper.Hooks().AfterMarketCreated(sdkCtx, market)
			if err != nil {
				return nil, err
			}
		}

		err = marketmapkeeper.SetParams(
			sdkCtx,
			marketmaptypes.Params{
				Admin: authtypes.NewModuleAddress(govtypes.ModuleName).String(), // governance. allows gov to add or remove market authorities.
				// market authority addresses may add and update markets to the x/marketmap module
				MarketAuthorities: []string{
					authtypes.NewModuleAddress(govtypes.ModuleName).String(), // governance
				},
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to set x/marketmap params: %w", err)
		}
		return migrations, nil
	}
}
