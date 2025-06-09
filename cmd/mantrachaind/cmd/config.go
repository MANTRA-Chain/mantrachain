package cmd

import (
	"time"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/MANTRA-Chain/mantrachain/v5/app"
	cmtcfg "github.com/cometbft/cometbft/config"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	cosmosevmserverconfig "github.com/cosmos/evm/server/config"
	oracleconfig "github.com/skip-mev/connect/v2/oracle/config"
)

// initCometBFTConfig helps to override default CometBFT Config values.
// return cmtcfg.DefaultConfig if no custom configuration is required for the application.
func initCometBFTConfig() *cmtcfg.Config {
	cfg := cmtcfg.DefaultConfig()

	// increase the number of inbound and outbound peers
	cfg.P2P.MaxNumInboundPeers = 100
	cfg.P2P.MaxNumOutboundPeers = 40
	cfg.Consensus.TimeoutCommit = 2 * time.Second

	return cfg
}

// initAppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func initAppConfig() (string, interface{}) {
	// The following code snippet is just for reference.
	type CustomAppConfig struct {
		serverconfig.Config

		EVM     cosmosevmserverconfig.EVMConfig
		JSONRPC cosmosevmserverconfig.JSONRPCConfig
		TLS     cosmosevmserverconfig.TLSConfig
		Wasm    wasmtypes.NodeConfig   `mapstructure:"wasm"`
		Oracle  oracleconfig.AppConfig `mapstructure:"oracle" json:"oracle"`
	}

	// Optionally allow the chain developer to overwrite the SDK's default
	// server config.
	srvCfg := serverconfig.DefaultConfig()

	oracleCfg := oracleconfig.AppConfig{
		Enabled:        false,
		OracleAddress:  "localhost:8080",
		ClientTimeout:  time.Second * 2,
		MetricsEnabled: false,
	}

	evmConfig := cosmosevmserverconfig.DefaultEVMConfig()
	evmConfig.EVMChainID = app.GetMANTRAEVMChainId()

	customAppConfig := CustomAppConfig{
		Config:  *srvCfg,
		EVM:     *evmConfig,
		JSONRPC: *cosmosevmserverconfig.DefaultJSONRPCConfig(),
		TLS:     *cosmosevmserverconfig.DefaultTLSConfig(),
		Wasm:    wasmtypes.DefaultNodeConfig(),
		Oracle:  oracleCfg,
	}
	// The SDK's default minimum gas price is set to "" (empty value) inside
	// app.toml. If left empty by validators, the node will halt on startup.
	// However, the chain developer can set a default app.toml value for their
	// validators here.
	//
	// In summary:
	// - if you leave srvCfg.MinGasPrices = "", all validators MUST tweak their
	//   own app.toml config,
	// - if you set srvCfg.MinGasPrices non-empty, validators CAN tweak their
	//   own app.toml to override, or use this default value.
	//
	// In this example application, we set the min gas prices to 0.
	srvCfg.MinGasPrices = "0uom"

	customAppTemplate := serverconfig.DefaultConfigTemplate +
		cosmosevmserverconfig.DefaultEVMConfigTemplate +
		wasmtypes.DefaultConfigTemplate() +
		oracleconfig.DefaultConfigTemplate

	return customAppTemplate, customAppConfig
}
