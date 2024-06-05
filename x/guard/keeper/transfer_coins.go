package keeper

import (
	"math/big"
	"strings"

	coinfactorytypes "github.com/MANTRA-Finance/mantrachain/x/coinfactory/types"
	tokentypes "github.com/MANTRA-Finance/mantrachain/x/token/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
)

func (k Keeper) CheckCanTransferCoins(ctx sdk.Context, address sdk.AccAddress, coins sdk.Coins) error {
	if !k.HasGuardTransferCoins(ctx) {
		return nil
	}

	conf := k.GetParams(ctx)

	defaultPrivileges := big.NewInt(0).SetBytes(conf.DefaultPrivileges)
	inverseDefaultPrilileges := big.NewInt(0).Not(defaultPrivileges)

	nftCollectionCreator := sdk.MustAccAddressFromBech32(conf.AccountPrivilegesTokenCollectionCreator)
	nftCollectionIndex := tokentypes.GetNftCollectionIndex(nftCollectionCreator, conf.AccountPrivilegesTokenCollectionId)

	for _, coin := range coins {
		denom := coin.GetDenom()
		denomBytes := []byte(denom)

		if denom == conf.BaseDenom || strings.HasPrefix(denom, "pool") {
			continue
		}

		if strings.HasPrefix(denom, "factory/") {
			// verify that denom is an x/coinfactory denom
			_, _, err := coinfactorytypes.DeconstructDenom(denom)
			if err != nil {
				return err
			}

			coinAdmin, found := k.ck.GetAdmin(ctx, denom)

			if !found {
				return sdkerrors.Wrapf(types.ErrCoinAdminNotFound, "missing coin admin, denom %s", denom)
			}

			// The coin admin should be able to transfer without checking the privileges
			if coinAdmin.Equals(address) {
				continue
			}
		}

		err := k.CheckCanTransferCoin(ctx, nftCollectionCreator, nftCollectionIndex, inverseDefaultPrilileges, address, denomBytes)

		if err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) CheckCanTransferCoin(ctx sdk.Context, nftCollectionCreator sdk.AccAddress, nftCollectionIndex []byte, inverseDefaultPrilileges *big.Int, address sdk.AccAddress, denom []byte) error {
	nftIndex := tokentypes.GetNftIndex(nftCollectionIndex, address.String())
	nftOwner := k.nk.GetOwner(ctx, string(nftCollectionIndex), string(nftIndex))

	if nftOwner.Empty() || !address.Equals(nftOwner) {
		return sdkerrors.Wrapf(types.ErrMissingSoulBondNft, "missing soul bond nft, address %s", address)
	}

	privileges, found := k.GetRequiredPrivileges(ctx, denom, types.RequiredPrivilegesCoin)

	if !found || privileges == nil {
		return sdkerrors.Wrapf(types.ErrCoinRequiredPrivilegesNotFound, "coin required privileges not found, denom %s", string(denom))
	}

	requiredPrivileges := types.PrivilegesFromBytes(privileges)
	requiredPrivilegesWithoutDefault := big.NewInt(0).And(inverseDefaultPrilileges, requiredPrivileges.BigInt())

	if requiredPrivilegesWithoutDefault.Cmp(big.NewInt(0)) == 0 {
		return sdkerrors.Wrapf(types.ErrCoinRequiredPrivilegesNotSet, "coin required privileges not set, denom %s", string(denom))
	}

	hasPrivileges, err := k.CheckAccountFulfillsRequiredPrivileges(ctx, address, privileges)

	if err != nil {
		return err
	}

	if !hasPrivileges {
		k.Logger(ctx).Error("insufficient privileges", "address", address, "denom", string(denom))
		return sdkerrors.Wrapf(types.ErrInsufficientPrivileges, "insufficient privileges, address %s, denom %s", address, string(denom))
	}

	return nil
}

func (k Keeper) ValidateCoinsTransfers(ctx sdk.Context, inputs []banktypes.Input, outputs []banktypes.Output) error {
	if !k.HasGuardTransferCoins(ctx) {
		return nil
	}

	if len(inputs) == 0 {
		return nil
	}

	conf := k.GetParams(ctx)
	admin := sdk.MustAccAddressFromBech32(conf.AdminAccount)

	for i, out := range outputs {
		var in banktypes.Input
		if len(inputs) == 1 {
			in = inputs[0]
		} else {
			in = inputs[i]
		}

		inAddress, err := sdk.AccAddressFromBech32(in.Address)
		if err != nil {
			return err
		}

		// The admin can send coins to any address no matter if the recipient has soul bond nft and/or
		// the account privileges and no matter of the coin required privileges
		if k.whitelistTransfersAccAddrs[in.Address] || admin.Equals(inAddress) {
			if len(inputs) == 1 {
				return nil
			} else {
				continue
			}
		}

		err = k.CheckCanTransferCoins(ctx, inAddress, in.Coins)

		if err != nil {
			return err
		}

		outAddress, err := sdk.AccAddressFromBech32(out.Address)
		if err != nil {
			return err
		}

		if k.whitelistTransfersAccAddrs[out.Address] || admin.Equals(outAddress) {
			continue
		}

		err = k.CheckCanTransferCoins(ctx, outAddress, out.Coins)

		if err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) AddTransferAccAddressesWhitelist(addresses []string) []string {
	updated := make([]string, 0)

	if len(addresses) == 0 {
		return updated
	}

	for _, address := range addresses {
		val, ok := k.whitelistTransfersAccAddrs[address]

		if !ok || !val {
			k.whitelistTransfersAccAddrs[address] = true
			updated = append(updated, address)
		}
	}

	return updated
}

func (k Keeper) RemoveTransferAccAddressesWhitelist(addresses []string) {
	for _, address := range addresses {
		_, ok := k.whitelistTransfersAccAddrs[address]

		if ok {
			delete(k.whitelistTransfersAccAddrs, address)
		}
	}
}
