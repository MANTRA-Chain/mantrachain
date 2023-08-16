package keeper

import (
	"strings"

	coinfactorytypes "github.com/MANTRA-Finance/mantrachain/x/coinfactory/types"
	tokentypes "github.com/MANTRA-Finance/mantrachain/x/token/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

			if !found {
				return sdkerrors.Wrapf(types.ErrCoinAdminNotFound, "missing coin admin, denom %s", denom)
			}

			// The coin admin should be able to transfer without checking the privileges
			if coinAdmin.Equals(address) {
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
			return sdkerrors.Wrap(types.ErrInvalidTokenCollectionId, "nft collection id should not be empty")
		}

		creator, err := sdk.AccAddressFromBech32(collectionCreator)

		if err != nil {
			return sdkerrors.Wrap(types.ErrInvalidTokenCollectionCreator, "collection creator should not be empty")
		}

		collectionIndex := tokentypes.GetNftCollectionIndex(creator, collectionId)
		index := tokentypes.GetNftIndex(collectionIndex, address.String())
		owner := k.nk.GetOwner(ctx, string(collectionIndex), string(index))

		if owner.Empty() || !address.Equals(owner) {
			return sdkerrors.Wrapf(types.ErrMissingSoulBondNft, "missing soul bond nft, address %s", address)
		}

		requiredPrivilegesList := k.GetRequiredPrivilegesMany(ctx, indexes, types.RequiredPrivilegesCoin)

		for i, privileges := range requiredPrivilegesList {
			if privileges == nil {
				return sdkerrors.Wrapf(types.ErrCoinRequiredPrivilegesNotFound, "coin required privileges not found, denom %s", string(indexes[i]))
			}

			hasPrivileges, err := k.CheckAccountFulfillsRequiredPrivileges(ctx, address, privileges)

			if err != nil || !hasPrivileges {
				k.Logger(ctx).Error("insufficient privileges", "address", address, "denom", string(indexes[i]))
				return sdkerrors.Wrapf(types.ErrInsufficientPrivileges, "insufficient privileges, address %s, denom %s", address, string(indexes[i]))
			}
		}
	}

	return nil
}

func (k Keeper) ValidateCoinsTransfers(ctx sdk.Context, inputs []banktypes.Input, outputs []banktypes.Output) error {
	if !k.HasGuardTransferCoins(ctx) {
		return nil
	}

	if len(inputs) == 0 && len(outputs) == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrLogic, "inputs and outputs length not equal")
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

		// The admin can send coins to any address no matter if the recipient has soul bond nft and/or
		// the account privileges and no matter of the coin required privileges
		if k.whitelistTransfersAccAddrs[in.Address] ||
			admin.Equals(inAddress) {
			return nil
		}

		err = k.CheckCanTransferCoins(ctx, inAddress, in.Coins)

		if err != nil {
			return err
		}

		if k.whitelistTransfersAccAddrs[out.Address] ||
			admin.Equals(outAddress) {
			return nil
		}

		err = k.CheckCanTransferCoins(ctx, outAddress, out.Coins)

		if err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) WhitelistTransferAccAddresses(addresses []string, isWhitelisted bool) []string {
	updated := make([]string, 0)

	if len(addresses) == 0 {
		return updated
	}

	for _, address := range addresses {
		val, ok := k.whitelistTransfersAccAddrs[address]

		if ok && !isWhitelisted {
			delete(k.whitelistTransfersAccAddrs, address)
			updated = append(updated, address)
		} else if !val && isWhitelisted {
			k.whitelistTransfersAccAddrs[address] = isWhitelisted
			updated = append(updated, address)
		}
	}

	return updated
}
