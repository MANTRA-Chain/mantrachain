package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ LiquidityHooks = MultiLiquidityHooks{}

type MultiLiquidityHooks []LiquidityHooks

func NewMultiLiquidityHooks(hooks ...LiquidityHooks) MultiLiquidityHooks {
	return hooks
}

func (h MultiLiquidityHooks) OnProvideLiquidity(ctx sdk.Context, receiver sdk.Address, pairId uint64, poolId uint64, poolCoin sdk.Coin) error {
	var err error
	for i := range h {
		err = h[i].OnProvideLiquidity(ctx, receiver, pairId, poolId, poolCoin)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h MultiLiquidityHooks) OnWithdrawLiquidity(ctx sdk.Context, receiver sdk.Address, pairId uint64, poolId uint64, poolCoin sdk.Coin) error {
	var err error
	for i := range h {
		err = h[i].OnWithdrawLiquidity(ctx, receiver, pairId, poolId, poolCoin)
		if err != nil {
			return err
		}
	}
	return nil
}
