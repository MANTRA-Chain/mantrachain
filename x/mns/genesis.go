package mns

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/LimeChain/mantrachain/x/mns/keeper"
	"github.com/LimeChain/mantrachain/x/mns/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
    // Set all the domain
for _, elem := range genState.DomainList {
	k.SetDomain(ctx, elem)
}
// Set all the domainName
for _, elem := range genState.DomainNameList {
	k.SetDomainName(ctx, elem)
}
// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

    genesis.DomainList = k.GetAllDomain(ctx)
genesis.DomainNameList = k.GetAllDomainName(ctx)
// this line is used by starport scaffolding # genesis/module/export

    return genesis
}
