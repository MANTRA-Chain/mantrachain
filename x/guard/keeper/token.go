package keeper

import (
	"strings"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/MANTRA-Finance/mantrachain/x/token/types"
)

func (k Keeper) CheckNewRestrictedNftsCollection(ctx sdk.Context, restrictedNfts bool, account string) error {
	conf := k.GetParams(ctx)

	if strings.TrimSpace(conf.AdminAccount) == "" {
		return errors.Wrap(types.ErrInvalidAccount, "missing admin account in params")
	}

	admin := sdk.MustAccAddressFromBech32(conf.AdminAccount)

	isAdmin := admin.Equals(sdk.MustAccAddressFromBech32(account))

	if restrictedNfts && !isAdmin {
		return errors.Wrap(errorstypes.ErrUnauthorized, "not an admin")
	}

	return nil
}

func (k Keeper) CheckRestrictedNftsCollection(ctx sdk.Context, collectionCreator string, collectionId string, account string) error {
	conf := k.GetParams(ctx)

	if strings.TrimSpace(conf.AdminAccount) == "" {
		return errors.Wrap(types.ErrInvalidAccount, "missing admin account in params")
	}

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

	if k.tokenKeeper.HasRestrictedNftsCollection(
		ctx,
		collectionIndex,
	) && !isAdmin {
		return errors.Wrap(errorstypes.ErrUnauthorized, "restricted nfts colection")
	}

	return nil
}

func (k Keeper) GetAccountPrivilegesTokenCollectionCreatorAndCollectionId(ctx sdk.Context) (string, string) {
	conf := k.GetParams(ctx)

	return conf.GetAccountPrivilegesTokenCollectionCreator(), conf.GetAccountPrivilegesTokenCollectionId()
}
