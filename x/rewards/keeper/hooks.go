package keeper

import (
	liquiditytypes "github.com/MANTRA-Finance/aumega/x/liquidity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Wrapper struct
type Hooks struct {
	k Keeper
}

var _ liquiditytypes.LiquidityHooks = Hooks{}

func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

func (h Hooks) AfterPoolCoinMinted(ctx sdk.Context, receiver sdk.Address, pairId uint64, poolId uint64, poolCoin sdk.Coin) {

}

func (h Hooks) AfterPoolCoinBurned(ctx sdk.Context, receiver sdk.Address, pairId uint64, poolId uint64, poolCoin sdk.Coin) {

}
