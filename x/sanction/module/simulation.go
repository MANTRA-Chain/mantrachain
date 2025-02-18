package sanction

import (
	"math/rand"

	"github.com/MANTRA-Chain/mantrachain/v2/testutil/sample"
	sanctionsimulation "github.com/MANTRA-Chain/mantrachain/v2/x/sanction/simulation"
	"github.com/MANTRA-Chain/mantrachain/v2/x/sanction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// avoid unused import issue
var (
	_ = sanctionsimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgUpsertFeeDenom = "op_weight_msg_upsert_fee_denom"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpsertFeeDenom int = 100

	opWeightMsgRemoveFeeDenom = "op_weight_msg_remove_fee_denom"
	// TODO: Determine the simulation weight value
	defaultWeightMsgRemoveFeeDenom int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	sanctionGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&sanctionGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgUpsertFeeDenom int
	simState.AppParams.GetOrGenerate(opWeightMsgUpsertFeeDenom, &weightMsgUpsertFeeDenom, nil,
		func(_ *rand.Rand) {
			weightMsgUpsertFeeDenom = defaultWeightMsgUpsertFeeDenom
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpsertFeeDenom,
		sanctionsimulation.SimulateMsgAddBlacklistAccount(am.keeper),
	))

	var weightMsgRemoveFeeDenom int
	simState.AppParams.GetOrGenerate(opWeightMsgRemoveFeeDenom, &weightMsgRemoveFeeDenom, nil,
		func(_ *rand.Rand) {
			weightMsgRemoveFeeDenom = defaultWeightMsgRemoveFeeDenom
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRemoveFeeDenom,
		sanctionsimulation.SimulateMsgRemoveBlacklistAccount(am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpsertFeeDenom,
			defaultWeightMsgUpsertFeeDenom,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				sanctionsimulation.SimulateMsgAddBlacklistAccount(am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgRemoveFeeDenom,
			defaultWeightMsgRemoveFeeDenom,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				sanctionsimulation.SimulateMsgRemoveBlacklistAccount(am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
