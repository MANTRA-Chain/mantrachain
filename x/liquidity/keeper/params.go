package keeper

import (
	"context"
	"time"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/MANTRA-Finance/mantrachain/x/liquidity/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx context.Context) (params types.Params) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams set the params
func (k Keeper) SetParams(ctx context.Context, params types.Params) error {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz, err := k.cdc.Marshal(&params)
	if err != nil {
		return err
	}
	store.Set(types.ParamsKey, bz)

	return nil
}

// GetBatchSize returns the current batch size parameter.
func (k Keeper) GetBatchSize(ctx sdk.Context) (batchSize uint32) {
	return k.GetParams(ctx).BatchSize
}

// GetTickPrecision returns the current tick precision parameter.
func (k Keeper) GetTickPrecision(ctx sdk.Context) (tickPrec uint32) {
	return k.GetParams(ctx).TickPrecision
}

// GetSwapFeeRate returns the current swap fee rate parameter.
func (k Keeper) GetSwapFeeRate(ctx sdk.Context) (swapFeeRate sdkmath.LegacyDec) {
	return k.GetParams(ctx).SwapFeeRate
}

// GetSwapFeeRate returns the current swap fee rate parameter.
func (k Keeper) GetPairCreatorSwapFeeRatio(ctx sdk.Context) (pairCreatorSwapFeeRatio sdkmath.LegacyDec) {
	return k.GetParams(ctx).PairCreatorSwapFeeRatio
}

// GetFeeCollector returns the current fee collector address parameter.
func (k Keeper) GetFeeCollector(ctx sdk.Context) sdk.AccAddress {
	feeCollectorAddr := k.GetParams(ctx).FeeCollectorAddress
	addr, err := sdk.AccAddressFromBech32(feeCollectorAddr)
	if err != nil {
		panic(err)
	}
	return addr
}

// GetDustCollector returns the current dust collector address parameter.
func (k Keeper) GetDustCollector(ctx sdk.Context) sdk.AccAddress {
	dustCollectorAddr := k.GetParams(ctx).DustCollectorAddress
	addr, err := sdk.AccAddressFromBech32(dustCollectorAddr)
	if err != nil {
		panic(err)
	}
	return addr
}

// GetMinInitialPoolCoinSupply returns the current minimum pool coin supply
// parameter.
func (k Keeper) GetMinInitialPoolCoinSupply(ctx sdk.Context) (i sdkmath.Int) {
	return k.GetParams(ctx).MinInitialPoolCoinSupply
}

// GetPairCreationFee returns the current pair creation fee parameter.
func (k Keeper) GetPairCreationFee(ctx sdk.Context) (fee sdk.Coins) {
	return k.GetParams(ctx).PairCreationFee
}

// GetPoolCreationFee returns the current pool creation fee parameter.
func (k Keeper) GetPoolCreationFee(ctx sdk.Context) (fee sdk.Coins) {
	return k.GetParams(ctx).PoolCreationFee
}

// GetMinInitialDepositAmount returns the current minimum initial deposit
// amount parameter.
func (k Keeper) GetMinInitialDepositAmount(ctx sdk.Context) (amt sdkmath.Int) {
	return k.GetParams(ctx).MinInitialDepositAmount
}

// GetMaxPriceLimitRatio returns the current maximum price limit ratio
// parameter.
func (k Keeper) GetMaxPriceLimitRatio(ctx sdk.Context) (ratio sdkmath.LegacyDec) {
	return k.GetParams(ctx).MaxPriceLimitRatio
}

// GetMaxNumMarketMakingOrderTicks returns the current maximum number of
// market making order ticks.
func (k Keeper) GetMaxNumMarketMakingOrderTicks(ctx sdk.Context) (i uint32) {
	return k.GetParams(ctx).MaxNumMarketMakingOrderTicks
}

// GetMaxOrderLifespan returns the current maximum order lifespan
// parameter.
func (k Keeper) GetMaxOrderLifespan(ctx sdk.Context) (maxLifespan time.Duration) {
	return k.GetParams(ctx).MaxOrderLifespan
}

// GetWithdrawFeeRate returns the current withdraw fee rate parameter.
func (k Keeper) GetWithdrawFeeRate(ctx sdk.Context) (feeRate sdkmath.LegacyDec) {
	return k.GetParams(ctx).WithdrawFeeRate
}

// GetDepositExtraGas returns the current deposit extra gas parameter.
func (k Keeper) GetDepositExtraGas(ctx sdk.Context) (gas storetypes.Gas) {
	return k.GetParams(ctx).DepositExtraGas
}

// GetWithdrawExtraGas returns the current withdraw extra gas parameter.
func (k Keeper) GetWithdrawExtraGas(ctx sdk.Context) (gas storetypes.Gas) {
	return k.GetParams(ctx).WithdrawExtraGas
}

// GetOrderExtraGas returns the current order extra gas parameter.
func (k Keeper) GetOrderExtraGas(ctx sdk.Context) (gas storetypes.Gas) {
	return k.GetParams(ctx).OrderExtraGas
}

// SetMaxNumMarketMakingOrderTicks sets max num market making order ticks
func (k Keeper) SetMaxNumMarketMakingOrderTicks(ctx sdk.Context, input uint32) {
	params := k.GetParams(ctx)
	params.MaxNumMarketMakingOrderTicks = input
	if err := k.SetParams(ctx, params); err != nil {
		panic(err)
	}
}

// GetMaxNumActivePoolsPerPair returns the current maximum number of active
// pools per pair.
func (k Keeper) GetMaxNumActivePoolsPerPair(ctx sdk.Context) (i uint32) {
	return k.GetParams(ctx).MaxNumActivePoolsPerPair
}

// SetMaxNumActivePoolsPerPair sets the maximum number of active pools per pair.
func (k Keeper) SetMaxNumActivePoolsPerPair(ctx sdk.Context, i uint32) {
	params := k.GetParams(ctx)
	params.MaxNumActivePoolsPerPair = i
	if err := k.SetParams(ctx, params); err != nil {
		panic(err)
	}
}
