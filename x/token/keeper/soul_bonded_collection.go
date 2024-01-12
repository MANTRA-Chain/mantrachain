package keeper

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/AumegaChain/aumega/x/token/types"
)

func (k Keeper) CheckSoulBondedNftsCollection(ctx sdk.Context, collectionCreator string, collectionId string) error {
	if strings.TrimSpace(collectionId) == "" {
		return errors.Wrap(types.ErrInvalidNftCollectionId, "nft collection id should not be empty")
	}

	creator, err := sdk.AccAddressFromBech32(collectionCreator)

	if err != nil {
		return status.Error(codes.InvalidArgument, "invalid collection creator")
	}

	collectionIndex := types.GetNftCollectionIndex(creator, collectionId)

	if k.HasSoulBondedNftsCollection(
		ctx,
		collectionIndex,
	) {
		return errors.Wrap(types.ErrSoulBondedNftCollectionOperationDisabled, "soul bonded nft collection operation disabled")
	}

	return nil
}
