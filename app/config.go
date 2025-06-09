package app

import (
	"fmt"
	"strconv"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
)

// MANTRAChainID is this chain's EVM Chain ID, to be overridden in makefile
var MANTRAChainID string

// GetMANTRAEVMChainId gets the EVM Chain ID for MANTRA chain
func GetMANTRAEVMChainId() uint64 {
	parsedChainId, err := strconv.ParseUint(MANTRAChainID, 10, 64)
	if err != nil {
		panic(err)
	}
	if parsedChainId == 0 {
		panic("Invalid chain ID")
	}
	return parsedChainId
}

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
var ChainsCoinInfo = map[uint64]evmtypes.EvmCoinInfo{
	GetMANTRAEVMChainId(): {
		Denom:         "uom",
		ExtendedDenom: "aom",
		DisplayDenom:  "om",
		Decimals:      evmtypes.SixDecimals,
	},
}

// EvmAppOptions allows to setup the global configuration
// for the Cosmos EVM chain.
func EvmAppOptions(chainID uint64) error {
	if sealed {
		return nil
	}

	coinInfo, found := ChainsCoinInfo[chainID]
	if !found {
		return fmt.Errorf("unknown chain id: %d", chainID)
	}

	// set the denom info for the chain
	if err := setBaseDenom(coinInfo); err != nil {
		return err
	}

	ethCfg := evmtypes.DefaultChainConfig(chainID)

	err := evmtypes.NewEVMConfigurator().
		WithExtendedEips(cosmosEVMActivators).
		WithChainConfig(ethCfg).
		WithEVMCoinInfo(coinInfo).
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
