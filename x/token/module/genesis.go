package token

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/MANTRA-Finance/mantrachain/x/token/keeper"
	"github.com/MANTRA-Finance/mantrachain/x/token/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}

	for _, elem := range genState.NftCollectionList {
		k.SetNftCollection(ctx, elem)
	}
	for _, elem := range genState.NftList {
		k.SetNft(ctx, elem)
	}
	// Set all the soulBondedNftsCollection
	for _, elem := range genState.SoulBondedNftsCollectionList {
		k.SetSoulBondedNftsCollection(ctx, elem)
	}
	// Set all the restrictedNftsCollection
	for _, elem := range genState.RestrictedNftsCollectionList {
		k.SetRestrictedNftsCollection(ctx, elem)
	}
	// Set all the openedNftsCollection
	for _, elem := range genState.OpenedNftsCollectionList {
		k.SetOpenedNftsCollection(ctx, elem)
	}
	// Set all the nftCollectionOwner
	for _, elem := range genState.NftCollectionOwnerList {
		k.SetNftCollectionOwner(ctx, elem.Index, elem.Owner)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.NftCollectionList = k.GetAllNftCollection(ctx)
	genesis.NftList = k.GetAllNft(ctx)

	genesis.SoulBondedNftsCollectionList = k.GetAllSoulBondedNftsCollection(ctx)
	genesis.RestrictedNftsCollectionList = k.GetAllRestrictedNftsCollection(ctx)
	genesis.OpenedNftsCollectionList = k.GetAllOpenedNftsCollection(ctx)
	genesis.NftCollectionOwnerList = k.GetAllNftCollectionOwner(ctx)

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}