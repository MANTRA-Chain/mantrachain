package cmd

import (
	"errors"
	"io"
	"os"

	"cosmossdk.io/log"
	confixcmd "cosmossdk.io/tools/confix/cmd"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmcli "github.com/CosmWasm/wasmd/x/wasm/client/cli"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/MANTRA-Chain/mantrachain/v6/app"
	cmtcli "github.com/cometbft/cometbft/libs/cli"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/pruning"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/snapshot"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	ethermintclient "github.com/cosmos/evm/client"
	evmdebug "github.com/cosmos/evm/client/debug"
	cosmosevmserver "github.com/cosmos/evm/server"
	cosmosevmserverflags "github.com/cosmos/evm/server/flags"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func initRootCmd(
	rootCmd *cobra.Command,
	txConfig client.TxConfig,
	basicManager module.BasicManager,
) {
	cfg := sdk.GetConfig()
	cfg.Seal()
	sdkAppCreator := func(l log.Logger, d dbm.DB, w io.Writer, ao servertypes.AppOptions) servertypes.Application {
		return newApp(l, d, w, ao)
	}

	rootCmd.AddCommand(
		genutilcli.InitCmd(basicManager, app.DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(
			banktypes.GenesisBalancesIterator{},
			app.DefaultNodeHome,
			genutiltypes.DefaultMessageValidator,
			txConfig.SigningContext().ValidatorAddressCodec(),
		),
		genutilcli.GenTxCmd(
			basicManager, txConfig,
			banktypes.GenesisBalancesIterator{},
			app.DefaultNodeHome,
			txConfig.SigningContext().ValidatorAddressCodec(),
		),
		genutilcli.ValidateGenesisCmd(basicManager),
		AddGenesisAccountCmd(app.DefaultNodeHome),
		cmtcli.NewCompletionCmd(rootCmd, true),
		evmdebug.Cmd(),
		pruning.Cmd(sdkAppCreator, app.DefaultNodeHome),
		confixcmd.ConfigCommand(),
		snapshot.Cmd(sdkAppCreator),
	)

	// add server commands
	cosmosevmserver.AddCommands(
		rootCmd,
		cosmosevmserver.NewDefaultStartOptions(newApp, app.DefaultNodeHome),
		appExport,
		addModuleInitFlags,
	)

	wasmcli.ExtendUnsafeResetAllCmd(rootCmd)

	// add keybase, auxiliary RPC, query, genesis, and tx child commands
	rootCmd.AddCommand(
		server.StatusCommand(),
		genesisCommand(txConfig, basicManager),
		queryCommand(),
		txCommand(),
		ethermintclient.KeyCommands(app.DefaultNodeHome, true),
	)

	_, err := cosmosevmserverflags.AddTxFlags(rootCmd)
	if err != nil {
		panic(err)
	}
}

func addModuleInitFlags(startCmd *cobra.Command) {
	crisis.AddModuleInitFlags(startCmd) //nolint:staticcheck
	wasm.AddModuleInitFlags(startCmd)
}

// genesisCommand builds genesis-related `mantrachaind genesis` command. Users may provide application specific commands as a parameter
func genesisCommand(txConfig client.TxConfig, basicManager module.BasicManager, cmds ...*cobra.Command) *cobra.Command {
	cmd := genutilcli.Commands(txConfig, basicManager, app.DefaultNodeHome)

	for _, subCmd := range cmds {
		cmd.AddCommand(subCmd)
	}
	return cmd
}

func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		rpc.QueryEventForTxCmd(),
		rpc.ValidatorCommand(),
		server.QueryBlockCmd(),
		server.QueryBlocksCmd(),
		server.QueryBlockResultsCmd(),
		authcmd.QueryTxsByEventsCmd(),
		authcmd.QueryTxCmd(),
	)

	return cmd
}

func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetMultiSignBatchCmd(),
		authcmd.GetValidateSignaturesCommand(),
		flags.LineBreak,
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
		authcmd.GetSimulateCmd(),
	)

	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

// newApp creates the application
func newApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	appOpts servertypes.AppOptions,
) cosmosevmserver.Application {
	baseappOptions := server.DefaultBaseappOptions(appOpts)
	var wasmOpts []wasmkeeper.Option
	if cast.ToBool(appOpts.Get("telemetry.enabled")) {
		wasmOpts = append(wasmOpts, wasmkeeper.WithVMCacheMetrics(prometheus.DefaultRegisterer))
	}

	return app.New(
		logger, db, traceStore, true,
		appOpts,
		wasmOpts,
		baseappOptions...,
	)
}

// appExport creates a new wasm app (optionally at a given height) and exports state.
func appExport(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	height int64,
	forZeroHeight bool,
	jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
	modulesToExport []string,
) (servertypes.ExportedApp, error) {
	var wasmApp *app.App
	// this check is necessary as we use the flag in x/upgrade.
	// we can exit more gracefully by checking the flag here.
	homePath, ok := appOpts.Get(flags.FlagHome).(string)
	if !ok || homePath == "" {
		return servertypes.ExportedApp{}, errors.New("application home is not set")
	}

	viperAppOpts, ok := appOpts.(*viper.Viper)
	if !ok {
		return servertypes.ExportedApp{}, errors.New("appOpts is not viper.Viper")
	}

	// overwrite the FlagInvCheckPeriod
	viperAppOpts.Set(server.FlagInvCheckPeriod, 1)
	appOpts = viperAppOpts

	var emptyWasmOpts []wasmkeeper.Option
	wasmApp = app.New(
		logger,
		db,
		traceStore,
		height == -1,
		appOpts,
		emptyWasmOpts,
	)

	if height != -1 {
		if err := wasmApp.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	}

	return wasmApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs, modulesToExport)
}

var tempDir = func() string {
	dir, err := os.MkdirTemp("", "wasmd")
	if err != nil {
		panic("failed to create temp dir: " + err.Error())
	}
	defer os.RemoveAll(dir)

	return dir
}
