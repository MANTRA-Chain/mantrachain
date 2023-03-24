package keeper

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	coinfactorytypes "github.com/LimeChain/mantrachain/x/coinfactory/types"
	"github.com/LimeChain/mantrachain/x/guard/types"
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
