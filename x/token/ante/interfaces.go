package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TokenKeeper interface {
	CheckCanTransferNft(ctx sdk.Context, tokenKeeper TokenKeeper, classId string) (bool, error)
	CheckIsSoulBondedCollection(ctx sdk.Context, tokenKeeper TokenKeeper, collectionCreator string, collectionId string) (bool, error)
}
