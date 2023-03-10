package keeper

import (
	"strings"

	"cosmossdk.io/errors"
	ante "github.com/LimeChain/mantrachain/x/guard/ante"
	tokentypes "github.com/LimeChain/mantrachain/x/token/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/LimeChain/mantrachain/x/guard/types"
)

func (k Keeper) CheckCanTransferCoins(ctx sdk.Context, nftKeeper ante.NFTKeeper, addresses []sdk.AccAddress, amount sdk.Coins) (bool, error) {
	conf := k.GetParams(ctx)

	if !k.HasGuardTransferCoins(ctx) {
		return true, nil
	}

	collectionCreator := conf.AccountPrivilegesTokenCollectionCreator
	collectionId := conf.AccountPrivilegesTokenCollectionId

	if strings.TrimSpace(collectionId) == "" {
		return false, errors.Wrap(types.ErrInvalidTokenCollectionId, "nft collection id should not be empty")
	}

	creator, err := sdk.AccAddressFromBech32(collectionCreator)

	if err != nil {
		return false, errors.Wrap(types.ErrInvalidTokenCollectionCreator, "collection creator should not be empty")
	}

	collectionIndex := tokentypes.GetNftCollectionIndex(creator, collectionId)

	for _, address := range addresses {
		index := tokentypes.GetNftIndex(collectionIndex, address.String())

		owner := nftKeeper.GetOwner(ctx, string(collectionIndex), string(index))

		if owner.Empty() || !address.Equals(owner) {
			return false, errors.Wrapf(types.ErrIncorrectNftOwner, "incorrect nft owner, address %s", address.String())
		}

		// TODO: not correct, use the corrct denom privileges
		hasPrivileges, err := k.HasPrivileges(ctx, address, []byte{})

		if err != nil {
			return false, err
		}

		if !hasPrivileges {
			return false, errors.Wrapf(types.ErrInsufficientPrivileges, "insufficient privileges, address %s", address.String())
		}
	}

	return true, nil
}
