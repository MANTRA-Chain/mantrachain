package keeper

import (
	"context"
	"time"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/MANTRA-Finance/mantrachain/x/lpfarm/types"
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

func (k Keeper) GetPrivatePlanCreationFee(ctx context.Context) (fee sdk.Coins) {
	return k.GetParams(ctx).PrivatePlanCreationFee
}

func (k Keeper) SetPrivatePlanCreationFee(ctx context.Context, fee sdk.Coins) {
	params := k.GetParams(ctx)
	params.PrivatePlanCreationFee = fee
	if err := k.SetParams(ctx, params); err != nil {
		panic(err)
	}
}

func (k Keeper) GetFeeCollector(ctx context.Context) (feeCollector string) {
	return k.GetParams(ctx).FeeCollector
}

func (k Keeper) SetFeeCollector(ctx context.Context, feeCollector string) {
	params := k.GetParams(ctx)
	params.FeeCollector = feeCollector
	if err := k.SetParams(ctx, params); err != nil {
		panic(err)
	}
}

func (k Keeper) GetMaxNumPrivatePlans(ctx context.Context) (num uint32) {
	return k.GetParams(ctx).MaxNumPrivatePlans
}

func (k Keeper) SetMaxNumPrivatePlans(ctx context.Context, num uint32) {
	params := k.GetParams(ctx)
	params.MaxNumPrivatePlans = num
	if err := k.SetParams(ctx, params); err != nil {
		panic(err)
	}
}

func (k Keeper) GetMaxBlockDuration(ctx context.Context) (d time.Duration) {
	return k.GetParams(ctx).MaxBlockDuration
}

func (k Keeper) SetMaxBlockDuration(ctx context.Context, d time.Duration) {
	params := k.GetParams(ctx)
	params.MaxBlockDuration = d
	if err := k.SetParams(ctx, params); err != nil {
		panic(err)
	}
}
