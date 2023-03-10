package keeper

import (
	"strings"

	"cosmossdk.io/errors"
	ante "github.com/LimeChain/mantrachain/x/token/ante"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/LimeChain/mantrachain/x/token/types"
)

func (k Keeper) CheckIsSoulBondedCollection(ctx sdk.Context, _ ante.TokenKeeper, collectionCreator string, collectionId string) (bool, error) {
	if strings.TrimSpace(collectionId) == "" {
		return false, errors.Wrap(types.ErrInvalidNftCollectionId, "nft collection id should not be empty")
	}

	creator, err := sdk.AccAddressFromBech32(collectionCreator)

	if err != nil {
		return false, status.Error(codes.InvalidArgument, "invalid collection creator")
	}

	collectionIndex := types.GetNftCollectionIndex(creator, collectionId)

	if k.HasSoulBondedNftsCollection(
		ctx,
		collectionIndex,
	) {
		return true, errors.Wrap(types.ErrSoulBondedNftCollectionOperationDisabled, "soul bonded nft collection operation disabled")
	}

	return false, nil
}
