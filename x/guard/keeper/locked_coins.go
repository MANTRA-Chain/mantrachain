package keeper

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/LimeChain/mantrachain/x/guard/types"
)

func (k Keeper) ValidateAreCoinsLocked(ctx sdk.Context, address sdk.Address, coins sdk.Coins) error {
	if len(coins) == 0 {
		return nil
	}

	conf := k.GetParams(ctx)
	admin := sdk.MustAccAddressFromBech32(conf.AdminAccount)

	// Check if it is a module address or it is an admin wallet address
	// TODO: change `admin.Equals(...` to update to `k.hasRole(...` when implemented
	if k.modAccAddrs[address.String()] ||
		admin.Equals(address) {
		return nil
	}

	for _, coin := range coins {
		denom := coin.GetDenom()

		exists := k.HasLocked(ctx, []byte(denom), types.LockedCoin)

		if exists {
			return errors.Wrapf(types.ErrCoinLocked, "coin %s locked", denom)
		}
	}

	return nil
}
