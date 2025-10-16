package app

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	clienthelpers "cosmossdk.io/client/v2/helpers"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/spf13/viper"
)

// ChainsCoinInfo is a map of the chain id and its corresponding EvmCoinInfo
// that allows initializing the app with different coin info based on the
// chain id
var ChainCoinInfo = evmtypes.EvmCoinInfo{
	Denom:         "uom",
	ExtendedDenom: "aom",
	DisplayDenom:  "om",
	Decimals:      evmtypes.SixDecimals.Uint32(),
}

var (
	EVMChainIDMap = map[string]uint64{
		"mantra-1":            5888, // mainnet Chain ID
		"mantra-dukong-1":     5887, // testnet Chain ID
		"mantra-canary-net-1": 7888, // devnet Chain ID
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
	var chainID string
	if _, err = os.Stat(genesisFilePath); err == nil {
		// File exists, read the genesis file to get the chain ID
		var reader *os.File
		reader, err = os.Open(genesisFilePath)
		if err == nil {
			chainID, err = genutiltypes.ParseChainIDFromGenesis(reader)
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
				return
			}
		}
	}

	if chainID != "" {
		// If chain ID not found in map, try parsing it
		evmChainID, err := ParseChainID(chainID)
		if err != nil {
			panic(err)
		}
		MANTRAChainID = evmChainID
		return
	}

	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
}

var (
	regexChainID         = `[a-z]{1,}`
	regexEIP155Separator = `_{1}`
	regexEIP155          = `[1-9][0-9]*`
	regexEpochSeparator  = `-{1}`
	regexEpoch           = `[1-9][0-9]*`
	evmosChainID         = regexp.MustCompile(fmt.Sprintf(`^(%s)%s(%s)%s(%s)$`,
		regexChainID,
		regexEIP155Separator,
		regexEIP155,
		regexEpochSeparator,
		regexEpoch))
)

// ParseChainID parses a string chain identifier's epoch to an Ethereum-compatible
// chain-id in *big.Int format. The function returns an error if the chain-id has an invalid format
func ParseChainID(chainID string) (uint64, error) {
	chainID = strings.TrimSpace(chainID)
	if len(chainID) > 48 {
		return 0, fmt.Errorf("chain-id '%s' cannot exceed 48 chars", chainID)
	}

	matches := evmosChainID.FindStringSubmatch(chainID)
	if matches == nil || len(matches) != 4 || matches[1] == "" {
		return 0, fmt.Errorf("%s: %v", chainID, matches)
	}

	// verify that the chain-id entered is a base 10 integer
	chainIDInt, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, fmt.Errorf("epoch %s must be base-10 integer format", matches[2])
	}

	return uint64(chainIDInt), nil
}
