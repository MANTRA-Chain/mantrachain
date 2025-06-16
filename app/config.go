package app

import (
	"os"
	"path/filepath"

	clienthelpers "cosmossdk.io/client/v2/helpers"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/spf13/viper"
)

// EVMOptionsFn defines a function type for setting app options specifically for
// the Cosmos EVM app. The function should receive the chainID and return an error if
// any.
type EVMOptionsFn func(uint64) error

var sealed = false

func NoOpEvmAppOptions(_ uint64) error {
	return nil
}

// ChainsCoinInfo is a map of the chain id and its corresponding EvmCoinInfo
// that allows initializing the app with different coin info based on the
// chain id
var ChainCoinInfo = evmtypes.EvmCoinInfo{
	Denom:         "uom",
	ExtendedDenom: "aom",
	DisplayDenom:  "om",
	Decimals:      evmtypes.SixDecimals,
}

// EvmAppOptions allows to setup the global configuration
// for the Cosmos EVM chain.
func EvmAppOptions(chainID uint64) error {
	if sealed {
		return nil
	}

	// set the denom info for the chain
	if err := setBaseDenom(ChainCoinInfo); err != nil {
		return err
	}

	ethCfg := evmtypes.DefaultChainConfig(chainID)

	err := evmtypes.NewEVMConfigurator().
		WithExtendedEips(cosmosEVMActivators).
		WithChainConfig(ethCfg).
		WithEVMCoinInfo(ChainCoinInfo).
		Configure()
	if err != nil {
		return err
	}

	sealed = true
	return nil
}

// setBaseDenom registers the display denom and base denom and sets the
// base denom for the chain.
func setBaseDenom(ci evmtypes.EvmCoinInfo) error {
	if err := sdk.RegisterDenom(ci.DisplayDenom, math.LegacyOneDec()); err != nil {
		return err
	}

	// sdk.RegisterDenom will automatically overwrite the base denom when the
	// new setBaseDenom() are lower than the current base denom's units.
	return sdk.RegisterDenom(ci.Denom, math.LegacyNewDecWithPrec(1, int64(ci.Decimals)))
}

var (
	EVMChainIDMap = map[string]uint64{
		"mantra-1":            5888, // mainnet Chain ID
		"mantra-dukong-1":     5887, // testnet Chain ID
		"mantra-canary-net-1": 5887, // devnet Chain ID
	}

	MANTRAChainID uint64 = 262144 // default Chain ID
)

// init initializes the MANTRAChainID variable by reading the chain ID from the
// genesis file or app.toml file in the node's home directory.
// If the genesis file exists, it reads the Cosmos chain ID from there and finds the EVM Chain ID
// against the EVMChainIDMap; otherwise, it checks the app.toml file for the EVM chain ID.
// If neither file exists or the chain ID is not found, it defaults to the MANTRA Chain ID (262144).
func init() {
	nodeHome, err := clienthelpers.GetNodeHomeDirectory(NodeDir)
	if err != nil {
		panic(err)
	}

	// check if the genesis file exists and read the chain ID from it
	genesisFilePath := filepath.Join(nodeHome, "config", "genesis.json")
	if _, err = os.Stat(genesisFilePath); err == nil {
		// File exists, read the genesis file to get the chain ID
		reader, err := os.Open(genesisFilePath)
		if err == nil {
			chainID, err := genutiltypes.ParseChainIDFromGenesis(reader)
			if err == nil && chainID != "" {
				evmChainID, found := EVMChainIDMap[chainID]
				if found {
					MANTRAChainID = evmChainID
					return
				}
			}
			defer reader.Close()
		}
	}
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	// If genesis file does not exist or chain ID is not found, check app.toml
	// to get the EVM chain ID
	appTomlPath := filepath.Join(nodeHome, "config", "app.toml")
	if _, err = os.Stat(appTomlPath); err == nil {
		// File exists
		v := viper.New()
		v.SetConfigFile(appTomlPath)
		v.SetConfigType("toml")

		if err = v.ReadInConfig(); err == nil {
			evmChainIDKey := "evm.evm-chain-id"
			if v.IsSet(evmChainIDKey) {
				evmChainID := v.GetUint64(evmChainIDKey)
				MANTRAChainID = evmChainID
			}
		}
	}
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
}
