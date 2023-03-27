package keeper

import (
	"strings"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/MANTRA-Finance/mantrachain/x/token/types"
)

func (k Keeper) CheckNewRestrictedNftsCollection(ctx sdk.Context, restrictedNfts bool, account string) error {
	conf := k.GetParams(ctx)
	admin := sdk.MustAccAddressFromBech32(conf.AdminAccount)

	isAdmin := admin.Equals(sdk.MustAccAddressFromBech32(account))

	if restrictedNfts && !isAdmin {
		return errors.Wrap(sdkerrors.ErrUnauthorized, "not an admin")
	}

	return nil
}

func (k Keeper) CheckRestrictedNftsCollection(ctx sdk.Context, collectionCreator string, collectionId string, account string) error {
	conf := k.GetParams(ctx)
	admin := sdk.MustAccAddressFromBech32(conf.AdminAccount)

	if strings.TrimSpace(collectionId) == "" {
		return errors.Wrap(types.ErrInvalidNftCollectionId, "nft collection id should not be empty")
	}

	creator, err := sdk.AccAddressFromBech32(collectionCreator)

	if err != nil {
		return status.Error(codes.InvalidArgument, "invalid collection creator")
	}

	collectionIndex := types.GetNftCollectionIndex(creator, collectionId)

	isAdmin := admin.Equals(sdk.MustAccAddressFromBech32(account))

	if k.tk.HasRestrictedNftsCollection(
		ctx,
		collectionIndex,
	) && !isAdmin {
		return errors.Wrap(sdkerrors.ErrUnauthorized, "restricted nfts colection")
	}

	return nil
}
