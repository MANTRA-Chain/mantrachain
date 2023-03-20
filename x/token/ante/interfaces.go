package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TokenKeeper interface {
	CheckCanTransferNft(ctx sdk.Context, classId string) (bool, error)
	CheckSoulBondedNftsCollection(ctx sdk.Context, collectionCreator string, collectionId string) error
}
