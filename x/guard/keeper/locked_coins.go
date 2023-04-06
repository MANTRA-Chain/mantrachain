package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	coinfactorytypes "github.com/MANTRA-Finance/mantrachain/x/coinfactory/types"
	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
)

func (k Keeper) ValidateCoinsLocked(ctx sdk.Context, address sdk.Address, coins sdk.Coins) error {
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
