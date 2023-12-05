package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ LiquidityHooks = MultiLiquidityHooks{}

type MultiLiquidityHooks []LiquidityHooks

func NewMultiLiquidityHooks(hooks ...LiquidityHooks) MultiLiquidityHooks {
	return hooks
}

func (h MultiLiquidityHooks) AfterPoolCoinMinted(ctx sdk.Context, receiver sdk.Address, pairId uint64, poolId uint64, poolCoin sdk.Coin) {
	for i := range h {
		h[i].AfterPoolCoinMinted(ctx, receiver, pairId, poolId, poolCoin)
	}
}

func (h MultiLiquidityHooks) AfterPoolCoinBurned(ctx sdk.Context, receiver sdk.Address, pairId uint64, poolId uint64, poolCoin sdk.Coin) {
	for i := range h {
		h[i].AfterPoolCoinBurned(ctx, receiver, pairId, poolId, poolCoin)
	}
}
