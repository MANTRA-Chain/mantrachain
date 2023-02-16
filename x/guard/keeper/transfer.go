package keeper

import (
	"strings"

	ante "github.com/LimeChain/mantrachain/x/guard/ante"
	tokentypes "github.com/LimeChain/mantrachain/x/token/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/LimeChain/mantrachain/x/guard/types"
	utils "github.com/LimeChain/mantrachain/x/guard/utils"
)

func (k Keeper) CheckCanTransfer(ctx sdk.Context, tokenKeeper ante.TokenKeeper, nftKeeper ante.NFTKeeper, addresses []sdk.AccAddress, amount sdk.Coins) (bool, error) {
	conf := k.GetParams(ctx)

	guardTansfer, found := k.GetGuardTransfer(ctx)

	if !found {
		return false, sdkerrors.Wrap(types.ErrGuardTransferNotFound, "guard transfer not found")
	}

	// Check if guard transfer is enabled
	if !guardTansfer.Enabled {
		return true, nil
	}

	collectionCreator := conf.TokenCollectionCreator
	collectionId := conf.TokenCollectionId

	if strings.TrimSpace(collectionId) == "" {
		return false, sdkerrors.Wrap(types.ErrInvalidTokenCollectionId, "token collection id should not be empty")
	}

	creator, err := sdk.AccAddressFromBech32(collectionCreator)

	if err != nil {
		return false, sdkerrors.Wrap(types.ErrInvalidTokenCollectionCreator, "token collection creator should not be empty")
	}

	collectionIndex := tokentypes.GetNftCollectionIndex(creator, collectionId)

	for _, address := range addresses {
		index := tokentypes.GetNftIndex(collectionIndex, address.String())

		owner := nftKeeper.GetOwner(ctx, string(collectionIndex), string(index))

		if owner.Empty() || !address.Equals(owner) {
			return false, sdkerrors.Wrapf(types.ErrIncorrectNftOwner, "incorrect nft owner, address %s", address.String())
		}

		priviliges := conf.DefaultPriviliges

		accPerm, found := k.GetAccPerm(ctx, address.String())
		if found {
			priviliges = accPerm.Priviliges
		}

		if !utils.CheckPriviliges(priviliges, types.PRIVILIGE_TRANSFER) {
			return false, sdkerrors.Wrapf(types.ErrInsufficientPriviliges, "insufficient priviliges, address %s", address.String())
		}
	}

	return true, nil
}
