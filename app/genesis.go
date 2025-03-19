package app

import (
	"encoding/json"

	math "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	erc20types "github.com/cosmos/evm/x/erc20/types"
	feemarkettypes "github.com/cosmos/evm/x/feemarket/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	marketmaptypes "github.com/skip-mev/connect/v2/x/marketmap/types"
	oracletypes "github.com/skip-mev/connect/v2/x/oracle/types"
)

// GenesisState of the blockchain is represented here as a map of raw json
// messages key'd by a identifier string.
// The identifier is used to determine which module genesis information belongs
// to so it may be appropriately routed during init chain.
// Within this application default genesis information is retrieved from
// the ModuleBasicManager which populates json from each BasicModule
// object provided to it during init.
type GenesisState map[string]json.RawMessage

var FeeDenom = "uom"

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState(cdc codec.JSONCodec) GenesisState {
	genesisState := module.BasicManager{}.DefaultGenesis(cdc)

	oracleGenesis := oracletypes.DefaultGenesisState()
	oracleGenesisStateBytes, err := json.Marshal(oracleGenesis)
	if err != nil {
		panic("cannot marshal connect genesis state for tests")
	}
	genesisState[oracletypes.ModuleName] = oracleGenesisStateBytes

	marketmapGenesis := marketmaptypes.DefaultGenesisState()
	marketmapGenesisStateBytes, err := json.Marshal(marketmapGenesis)
	if err != nil {
		panic("cannot marshal connect genesis state for tests")
	}
	genesisState[marketmaptypes.ModuleName] = marketmapGenesisStateBytes

	distributionGenesis := distributiontypes.GenesisState{
		Params: distributiontypes.Params{
			CommunityTax: math.LegacyMustNewDecFromStr("0.01"),
			//			McaTax:              math.LegacyMustNewDecFromStr("0.4"),
			//			McaAddress:          "mantra15m77x4pe6w9vtpuqm22qxu0ds7vn4ehzwx8pls",
			WithdrawAddrEnabled: true,
		},
	}
	distributionGenesisStateBytes, err := json.Marshal(distributionGenesis)
	if err != nil {
		panic("cannot marshal distribution genesis state for tests")
	}
	genesisState[distributiontypes.ModuleName] = distributionGenesisStateBytes

	var feeMarketState feemarkettypes.GenesisState
	cdc.MustUnmarshalJSON(genesisState[feemarkettypes.ModuleName], &feeMarketState)
	feeMarketState.Params.NoBaseFee = false
	feeMarketState.Params.BaseFee = math.LegacyMustNewDecFromStr("0.01")
	feeMarketState.Params.MinGasPrice = feeMarketState.Params.BaseFee
	genesisState[feemarkettypes.ModuleName] = cdc.MustMarshalJSON(&feeMarketState)

	var evmState evmtypes.GenesisState
	cdc.MustUnmarshalJSON(genesisState[evmtypes.ModuleName], &evmState)
	evmState.Params.EvmDenom = FeeDenom
	genesisState[evmtypes.ModuleName] = cdc.MustMarshalJSON(&evmState)

	var erc20State erc20types.GenesisState
	cdc.MustUnmarshalJSON(genesisState[erc20types.ModuleName], &erc20State)
	erc20State.TokenPairs[0].Denom = FeeDenom
	genesisState[erc20types.ModuleName] = cdc.MustMarshalJSON(&erc20State)

	return genesisState
}
