package app

import (
	"os"
	"path/filepath"

	clienthelpers "cosmossdk.io/client/v2/helpers"
	"github.com/spf13/viper"
)

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
	if _, err := os.Stat(genesisFilePath); err == nil {
		// File exists, read the genesis file to get the chain ID
		v := viper.New()
		v.SetConfigFile(genesisFilePath)
		v.SetConfigType("json")
		if err := v.ReadInConfig(); err == nil {
			chainIDKey := "chain_id"
			if v.IsSet(chainIDKey) {
				chainID := v.GetString(chainIDKey)
				evmChainID, found := EVMChainIDMap[chainID]
				if found {
					MANTRAChainID = evmChainID
					return
				}
			}
		}
	}

	// If genesis file does not exist or chain ID is not found, check app.toml
	// to get the EVM chain ID
	appTomlPath := filepath.Join(nodeHome, "config", "app.toml")
	if _, err := os.Stat(appTomlPath); err == nil {
		// File exists
		v := viper.New()
		v.SetConfigFile(appTomlPath)
		v.SetConfigType("toml")

		if err := v.ReadInConfig(); err == nil {
			evmChainIDKey := "evm.evm-chain-id"
			if v.IsSet(evmChainIDKey) {
				evmChainID := v.GetUint64(evmChainIDKey)
				MANTRAChainID = evmChainID
			}
		}
	}
}
