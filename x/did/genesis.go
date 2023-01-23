package did

import (
	"github.com/LimeChain/mantrachain/x/did/keeper"
	"github.com/LimeChain/mantrachain/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	for _, elem := range genState.DidDocumentList {
		k.SetDidDocument(ctx, []byte(elem.Id), elem)
	}
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.DidDocumentList = k.GetAllDidDocuments(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
