package ante

import sdk "github.com/cosmos/cosmos-sdk/types"

type GuardKeeper interface {
	CheckIsAdmin(ctx sdk.Context, address string) error
	CheckNewRestrictedNftsCollection(ctx sdk.Context, restrictedNftsCollection bool, address string) error
	CheckRestrictedNftsCollection(ctx sdk.Context, collectionCreator string, collectionId string, address string) error
}
