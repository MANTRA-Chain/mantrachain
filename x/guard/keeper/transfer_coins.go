package keeper

import (
	"strings"

	coinfactorytypes "github.com/MANTRA-Finance/mantrachain/x/coinfactory/types"
	tokentypes "github.com/MANTRA-Finance/mantrachain/x/token/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
)

func (k Keeper) CheckCanTransferCoins(ctx sdk.Context, address sdk.AccAddress, coins sdk.Coins) error {
	var indexes [][]byte

	for _, coin := range coins {
		denom := coin.GetDenom()
		denomBytes := []byte(denom)

		// verify that denom is an x/coinfactory denom
		_, _, err := coinfactorytypes.DeconstructDenom(denom)
		if err == nil {
			coinAdmin, found := k.ck.GetAdmin(ctx, denom)

			if found && coinAdmin.Equals(address) {
				continue
			}

			indexes = append(indexes, denomBytes)
		}
	}

	if len(indexes) > 0 {
		conf := k.GetParams(ctx)

		collectionCreator := conf.AccountPrivilegesTokenCollectionCreator
		collectionId := conf.AccountPrivilegesTokenCollectionId

		if strings.TrimSpace(collectionId) == "" {
			return errors.Wrap(types.ErrInvalidTokenCollectionId, "nft collection id should not be empty")
		}

		creator, err := sdk.AccAddressFromBech32(collectionCreator)

		if err != nil {
			return errors.Wrap(types.ErrInvalidTokenCollectionCreator, "collection creator should not be empty")
		}

		collectionIndex := tokentypes.GetNftCollectionIndex(creator, collectionId)
		index := tokentypes.GetNftIndex(collectionIndex, address.String())
		owner := k.nk.GetOwner(ctx, string(collectionIndex), string(index))

		if owner.Empty() || !address.Equals(owner) {
			return errors.Wrapf(types.ErrIncorrectNftOwner, "incorrect nft owner, address %s", address)
		}

		requiredPrivilegesList := k.GetRequiredPrivilegesMany(ctx, indexes, types.RequiredPrivilegesCoin)

		if len(requiredPrivilegesList) == 0 || len(requiredPrivilegesList) != len(indexes) {
			return errors.Wrap(types.ErrCoinRequiredPrivilegesNotFound, "coin required privileges not found")
		}

		hasPrivileges, err := k.CheckAccountFulfillsRequiredPrivileges(ctx, address, requiredPrivilegesList)

		if err != nil {
			return err
		}

		if !hasPrivileges {
			k.Logger(ctx).Error("insufficient privileges", "address", address, "coins", coins)
			return errors.Wrapf(types.ErrInsufficientPrivileges, "insufficient privileges, address %s", address)
		}
	}

	return nil
}

func (k Keeper) ValidateCoinsTransfers(ctx sdk.Context, inputs []banktypes.Input, outputs []banktypes.Output) error {
	if !k.HasGuardTransferCoins(ctx) {
		return nil
	}

	if len(inputs) == 0 && len(outputs) == 0 {
		return errors.Wrapf(sdkerrors.ErrLogic, "inputs and outputs length not equal")
	}

	conf := k.GetParams(ctx)
	admin := sdk.MustAccAddressFromBech32(conf.AdminAccount)

	for i, in := range inputs {
		out := outputs[i]

		inAddress, err := sdk.AccAddressFromBech32(in.Address)
		if err != nil {
			return err
		}

		outAddress, err := sdk.AccAddressFromBech32(out.Address)
		if err != nil {
			return err
		}

		if k.whlstTransfersSendersAccAddrs[in.Address] ||
			admin.Equals(inAddress) {
			return nil
		}

		err = k.CheckCanTransferCoins(ctx, inAddress, in.Coins)

		if err != nil {
			return err
		}

		err = k.CheckCanTransferCoins(ctx, outAddress, out.Coins)

		if err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) WhlstTransferSendersAccAddresses(ctx sdk.Context, addresses []string, isWhitelisted bool) {
	for _, address := range addresses {
		val, ok := k.whlstTransfersSendersAccAddrs[address]

		if ok && !isWhitelisted {
			delete(k.whlstTransfersSendersAccAddrs, address)
		} else if !val && isWhitelisted {
			k.whlstTransfersSendersAccAddrs[address] = isWhitelisted
		}
	}
}
