package xfeemarket

import (
	"math/rand"

	"github.com/MANTRA-Chain/mantrachain/testutil/sample"
	xfeemarketsimulation "github.com/MANTRA-Chain/mantrachain/x/xfeemarket/simulation"
	"github.com/MANTRA-Chain/mantrachain/x/xfeemarket/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// avoid unused import issue
var (
	_ = xfeemarketsimulation.FindAccount
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
	xfeemarketGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&xfeemarketGenesis)
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
		xfeemarketsimulation.SimulateMsgUpsertFeeDenom(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgRemoveFeeDenom int
	simState.AppParams.GetOrGenerate(opWeightMsgRemoveFeeDenom, &weightMsgRemoveFeeDenom, nil,
		func(_ *rand.Rand) {
			weightMsgRemoveFeeDenom = defaultWeightMsgRemoveFeeDenom
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRemoveFeeDenom,
		xfeemarketsimulation.SimulateMsgRemoveFeeDenom(am.accountKeeper, am.bankKeeper, am.keeper),
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
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) govtypes.Content {
				return xfeemarketsimulation.RandomUpsertFeeDenomProposalContent(r, ctx, accs)
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgRemoveFeeDenom,
			defaultWeightMsgRemoveFeeDenom,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) govtypes.Content {
				return xfeemarketsimulation.RandomRemoveFeeDenomProposalContent(r, ctx, accs)
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}

// RandomUpsertFeeDenomProposalContent generates random proposal content for UpsertFeeDenom
func RandomUpsertFeeDenomProposalContent(r *rand.Rand, _ sdk.Context, _ []simtypes.Account) govtypes.Content {
	title := simtypes.RandStringOfLength(r, 10)
	description := simtypes.RandStringOfLength(r, 50)
	feeDenom := simtypes.RandStringOfLength(r, 3)

	content := &types.UpsertFeeDenomProposal{
		Title:       title,
		Description: description,
		FeeDenom:    feeDenom,
	}

	return content
}

// RandomRemoveFeeDenomProposalContent generates random proposal content for RemoveFeeDenom
func RandomRemoveFeeDenomProposalContent(r *rand.Rand, _ sdk.Context, _ []simtypes.Account) govtypes.Content {
	title := simtypes.RandStringOfLength(r, 10)
	description := simtypes.RandStringOfLength(r, 50)
	feeDenom := simtypes.RandStringOfLength(r, 3)

	content := &types.RemoveFeeDenomProposal{
		Title:       title,
		Description: description,
		FeeDenom:    feeDenom,
	}

	return content
}
