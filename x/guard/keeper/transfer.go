package keeper

import (
	"strings"

	"golang.org/x/exp/slices"

	ante "github.com/LimeChain/mantrachain/x/guard/ante"
	tokentypes "github.com/LimeChain/mantrachain/x/token/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/LimeChain/mantrachain/x/guard/types"
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

		meta, found := tokenKeeper.GetNft(ctx, collectionIndex, index)

		if !found {
			return false, sdkerrors.Wrapf(types.ErrTokenNftNotFound, "token nft not found, address %s", address.String())
		}

		owner := nftKeeper.GetOwner(ctx, string(collectionIndex), string(index))

		if !address.Equals(owner) {
			return false, sdkerrors.Wrapf(types.ErrIncorrectNftOwner, "incorrect nft owner, address %s", address.String())
		}

		if len(meta.Attributes) == 0 || meta.Attributes[0].Type != "AccPerm" {
			return false, sdkerrors.Wrapf(types.ErrTokenNftAttributesIncorrectOrNotFound, "incorrect nft attributes or not found, address %s", address.String())
		}

		accPermCat := meta.Attributes[0].Value

		accPerm, found := k.GetAccPerm(ctx, accPermCat)
		if !found {
			return false, sdkerrors.Wrapf(types.ErrAccPermNotFound, "acc perm not found, address %s", address.String())
		}

		if len(accPerm.WhlCurr) == 0 {
			return false, sdkerrors.Wrapf(types.ErrAccPermCatIncorrectOrNotFound, "incorrect acc perm cat or not found, address %s", address.String())
		}

		if accPerm.WhlCurr[0] == "*" {
			continue
		}

		for _, coin := range amount {
			if !slices.Contains(accPerm.WhlCurr, coin.Denom) {
				return false, sdkerrors.Wrapf(types.ErrAccPermCatIncorrectOrNotFound, "incorrect acc perm cat or not found, address %s", address.String())
			}
		}
	}

	return true, nil
}
