package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/MANTRA-Finance/mantrachain/x/coinfactory/types"
)

func (k Keeper) mintTo(ctx sdk.Context, amount sdk.Coin, mintTo string) error {
	// verify that denom is an x/coinfactory denom
	_, _, err := types.DeconstructDenom(amount.Denom)
	if err != nil {
		return err
	}

	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(amount))
	if err != nil {
		return err
	}

	addr, err := sdk.AccAddressFromBech32(mintTo)
	if err != nil {
		return err
	}
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName,
		addr,
		sdk.NewCoins(amount))
}

func (k Keeper) burnFrom(ctx sdk.Context, amount sdk.Coin, burnFrom string) error {
	// verify that denom is an x/coinfactory denom
	_, _, err := types.DeconstructDenom(amount.Denom)
	if err != nil {
		return err
	}

	addr, err := sdk.AccAddressFromBech32(burnFrom)
	if err != nil {
		return err
	}

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx,
		addr,
		types.ModuleName,
		sdk.NewCoins(amount))
	if err != nil {
		return err
	}

	return k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(amount))
}

func (k Keeper) forceTransfer(ctx sdk.Context, amount sdk.Coin, fromAddr string, toAddr string) error {
	from, err := sdk.AccAddressFromBech32(fromAddr)
	if err != nil {
		return err
	}

	to, err := sdk.AccAddressFromBech32(toAddr)
	if err != nil {
		return err
	}

	whitelisted := k.guardKeeper.AddTransferAccAddressesWhitelist(ctx, []sdk.AccAddress{from})
	err = k.bankKeeper.SendCoins(ctx, from, to, sdk.NewCoins(amount))
	k.guardKeeper.RemoveTransferAccAddressesWhitelist(ctx, whitelisted)
	if err != nil {
		return err
	}

	return nil
}
