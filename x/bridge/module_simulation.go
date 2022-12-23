package bridge

import (
	"math/rand"

	"github.com/LimeChain/mantrachain/testutil/sample"
	bridgesimulation "github.com/LimeChain/mantrachain/x/bridge/simulation"
	"github.com/LimeChain/mantrachain/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = bridgesimulation.FindAccount
	_ = simappparams.StakePerAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	opWeightMsgCreateCw20Contract = "op_weight_msg_cw_20_contract"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateCw20Contract int = 100

	opWeightMsgUpdateCw20Contract = "op_weight_msg_cw_20_contract"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateCw20Contract int = 100

	opWeightMsgDeleteCw20Contract = "op_weight_msg_cw_20_contract"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteCw20Contract int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	bridgeGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&bridgeGenesis)
}

// ProposalContents doesn't return any content functions for governance proposals
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreateCw20Contract int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreateCw20Contract, &weightMsgCreateCw20Contract, nil,
		func(_ *rand.Rand) {
			weightMsgCreateCw20Contract = defaultWeightMsgCreateCw20Contract
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateCw20Contract,
		bridgesimulation.SimulateMsgCreateCw20Contract(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateCw20Contract int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateCw20Contract, &weightMsgUpdateCw20Contract, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateCw20Contract = defaultWeightMsgUpdateCw20Contract
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateCw20Contract,
		bridgesimulation.SimulateMsgUpdateCw20Contract(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteCw20Contract int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgDeleteCw20Contract, &weightMsgDeleteCw20Contract, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteCw20Contract = defaultWeightMsgDeleteCw20Contract
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteCw20Contract,
		bridgesimulation.SimulateMsgDeleteCw20Contract(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
