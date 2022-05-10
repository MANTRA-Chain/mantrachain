package mns

import (
	"math/rand"

	"github.com/LimeChain/mantrachain/testutil/sample"
	mnssimulation "github.com/LimeChain/mantrachain/x/mns/simulation"
	"github.com/LimeChain/mantrachain/x/mns/types"
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
	_ = mnssimulation.FindAccount
	_ = simappparams.StakePerAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
    opWeightMsgCreateDomain = "op_weight_msg_create_domain"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateDomain int = 100

	opWeightMsgCreateDomainName = "op_weight_msg_create_domain_name"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateDomainName int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	mnsGenesis := types.GenesisState{
		Params:	types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&mnsGenesis)
}

// ProposalContents doesn't return any content functions for governance proposals
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized  param changes for the simulator
func (am AppModule) RandomizedParams(_ *rand.Rand) []simtypes.ParamChange {
	
	return []simtypes.ParamChange{
	}
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreateDomain int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreateDomain, &weightMsgCreateDomain, nil,
		func(_ *rand.Rand) {
			weightMsgCreateDomain = defaultWeightMsgCreateDomain
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateDomain,
		mnssimulation.SimulateMsgCreateDomain(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgCreateDomainName int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreateDomainName, &weightMsgCreateDomainName, nil,
		func(_ *rand.Rand) {
			weightMsgCreateDomainName = defaultWeightMsgCreateDomainName
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateDomainName,
		mnssimulation.SimulateMsgCreateDomainName(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
