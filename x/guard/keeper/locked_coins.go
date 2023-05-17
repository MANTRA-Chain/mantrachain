package keeper

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	coinfactorytypes "github.com/MANTRA-Finance/mantrachain/x/coinfactory/types"
	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
)

func (k Keeper) ValidateCoinLocked(ctx sdk.Context, coin sdk.Coin) error {
	denom := coin.GetDenom()

	_, _, err := coinfactorytypes.DeconstructDenom(denom)
	if err == nil {
		if k.HasLocked(ctx, []byte(denom), types.LockedCoin) {
			return errors.Wrapf(types.ErrCoinLocked, "coin %s locked", denom)
		}
	}

	return nil
}

func (k Keeper) ValidateCoinLockedByDenom(ctx sdk.Context, denom string) error {
	if strings.TrimSpace(denom) == "" {
		return errors.Wrapf(types.ErrInvalidDenom, "invalid denom")
	}

	_, _, err := coinfactorytypes.DeconstructDenom(denom)
	if err == nil {
		if k.HasLocked(ctx, []byte(denom), types.LockedCoin) {
			return errors.Wrapf(types.ErrCoinLocked, "coin %s locked", denom)
		}
	}

	return nil
}

func (k Keeper) ValidateCoinsLocked(ctx sdk.Context, coins sdk.Coins) error {
	if len(coins) == 0 {
		return nil
	}

	for _, coin := range coins {
		denom := coin.GetDenom()

		// verify that denom is an x/coinfactory denom
		_, _, err := coinfactorytypes.DeconstructDenom(denom)
		if err == nil {
			if k.HasLocked(ctx, []byte(denom), types.LockedCoin) {
				return errors.Wrapf(types.ErrCoinLocked, "coin %s locked", denom)
			}
		}
	}

	return nil
}

func (k Keeper) ValidateCoinsLockedByDenoms(ctx sdk.Context, denoms []string) error {
	if len(denoms) == 0 {
		return nil
	}

	for _, denom := range denoms {
		if strings.TrimSpace(denom) == "" {
			return errors.Wrapf(types.ErrInvalidDenom, "invalid denom")
		}

		// verify that denom is an x/coinfactory denom
		_, _, err := coinfactorytypes.DeconstructDenom(denom)
		if err == nil {
			if k.HasLocked(ctx, []byte(denom), types.LockedCoin) {
				return errors.Wrapf(types.ErrCoinLocked, "coin %s locked", denom)
			}
		}
	}

	return nil
}
