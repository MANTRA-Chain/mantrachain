package sanction

import (
	"math/rand"

	"github.com/MANTRA-Chain/mantrachain/v3/testutil/sample"
	sanctionsimulation "github.com/MANTRA-Chain/mantrachain/v3/x/sanction/simulation"
	"github.com/MANTRA-Chain/mantrachain/v3/x/sanction/types"
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
	opWeightMsgAddBlacklistAccounts = "op_weight_msg_add_blacklist_account"

	defaultWeightMsgAddBlacklistAccounts int = 5

	opWeightMsgRemoveBlacklistAccounts = "op_weight_msg_remove_blacklist_account"

	defaultWeightMsgRemoveBlacklistAccounts int = 5

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

	var weightMsgAddBlacklistAccount int
	simState.AppParams.GetOrGenerate(opWeightMsgAddBlacklistAccounts, &weightMsgAddBlacklistAccount, nil,
		func(_ *rand.Rand) {
			weightMsgAddBlacklistAccount = defaultWeightMsgAddBlacklistAccounts
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddBlacklistAccount,
		sanctionsimulation.SimulateMsgAddBlacklistAccounts(am.keeper),
	))

	var weightMsgRemoveBlacklistAccount int
	simState.AppParams.GetOrGenerate(opWeightMsgRemoveBlacklistAccounts, &weightMsgRemoveBlacklistAccount, nil,
		func(_ *rand.Rand) {
			weightMsgRemoveBlacklistAccount = defaultWeightMsgRemoveBlacklistAccounts
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRemoveBlacklistAccount,
		sanctionsimulation.SimulateMsgRemoveBlacklistAccounts(am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgAddBlacklistAccounts,
			defaultWeightMsgAddBlacklistAccounts,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				sanctionsimulation.SimulateMsgAddBlacklistAccounts(am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgRemoveBlacklistAccounts,
			defaultWeightMsgRemoveBlacklistAccounts,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				sanctionsimulation.SimulateMsgRemoveBlacklistAccounts(am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
