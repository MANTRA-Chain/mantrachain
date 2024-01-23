package did

import (
	"strings"

	"github.com/MANTRA-Finance/mantrachain/x/did/keeper"
	"github.com/MANTRA-Finance/mantrachain/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	for _, elem := range genState.DidDocuments {
		account := strings.TrimPrefix(elem.DidDocument.Id, types.DidChainPrefix)
		parts := strings.Split(account, ":")
		key := []byte(parts[len(parts)-1])
		k.SetDidDocument(ctx, key, elem.DidDocument)
		k.SetDidMetadata(ctx, key, elem.DidMetadata)
	}
	// this line is used by starport scaffolding # genesis/module/init
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	didDocuments := []types.GenesisDidDocument{}

	for _, elem := range k.GetAllDidDocuments(ctx) {
		account := strings.TrimPrefix(elem.Id, types.DidChainPrefix)
		parts := strings.Split(account, ":")
		didMetadata, found := k.GetDidMetadata(ctx, []byte(parts[len(parts)-1]))
		if !found {
			panic("cannot retrieve metadata for did document")
		}
		didDocuments = append(didDocuments, types.GenesisDidDocument{
			DidDocument: elem,
			DidMetadata: didMetadata,
		})
	}

	genesis.DidDocuments = didDocuments
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
