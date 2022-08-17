package keeper

import (
	"strconv"

	"github.com/LimeChain/mantrachain/x/vault/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) InitEpoch(ctx sdk.Context, chain string, validator string, denom string, bh int64) {
	newEpoch := types.Epoch{
		BlockStart:     bh,
		BlockEnd:       types.UndefinedBlockHeight,
		Rewards:        sdk.NewCoin(denom, sdk.NewInt(0)),
		Staked:         sdk.NewDec(0),
		PrevEpochBlock: types.UndefinedBlockHeight,
		NextEpochBlock: types.UndefinedBlockHeight,
		StartAt:        ctx.BlockHeader().Time.Unix(),
	}

	k.SetEpoch(ctx, chain, validator, denom, bh, newEpoch)
	k.SetLastEpochBlock(ctx, chain, validator, denom, types.LastEpochBlock{
		BlockHeight: bh,
	})
}

func (k Keeper) SetEpochEnd(
	ctx sdk.Context,
	chain string,
	validator string,
	denom string,
	bh int64,
	lastEpochBlock int64,
	minEpochWithdrawAmount int64,
) error {
	lastEpoch, found := k.GetEpoch(ctx, chain, validator, denom, lastEpochBlock)

	if !found {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "last epoch not found %s", lastEpochBlock)
	}

	de := NewDistributionExecutor(ctx, k.ac, k.sk, k.dk)
	amount, err := de.GetDelegationRewards(validator, denom)

	if err != nil {
		return err
	}

	newEpoch := types.Epoch{
		BlockStart:     bh,
		BlockEnd:       types.UndefinedBlockHeight,
		Rewards:        sdk.NewCoin(denom, sdk.NewInt(0)),
		PrevEpochBlock: lastEpochBlock,
		NextEpochBlock: types.UndefinedBlockHeight,
		StartAt:        ctx.BlockHeader().Time.Unix(),
	}

	se := NewStakingExecutor(ctx, k.ac, k.bk, k.sk)
	staked, err := se.GetDelegatorDelegation(validator)

	if err != nil {
		return err
	}

	newEpoch.Staked = staked

	if amount.GTE(sdk.NewDec(minEpochWithdrawAmount)) {
		rewards, err := de.WithdrawDelegationRewards(validator, denom)

		if err != nil {
			return err
		}

		lastEpoch.Rewards = rewards
	}

	lastEpoch.BlockEnd = bh
	lastEpoch.NextEpochBlock = bh
	lastEpoch.EndAt = ctx.BlockHeader().Time.Unix()

	k.SetEpoch(ctx, chain, validator, denom, lastEpochBlock, lastEpoch)

	newEpoch.Rewards = sdk.NewCoin(denom, sdk.NewInt(0))

	k.SetEpoch(ctx, chain, validator, denom, bh, newEpoch)
	k.SetLastEpochBlock(ctx, chain, validator, denom, types.LastEpochBlock{
		BlockHeight: bh,
	})

	return nil
}

func (k Keeper) HasEpoch(
	ctx sdk.Context, chain string, validator string, denom string, epochId int64,
) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.EpochStoreKey(chain, validator))
	index := types.GetEpochIndex(denom, []byte(strconv.FormatInt(epochId, 10)))
	return store.Has(index)
}

func (k Keeper) SetEpoch(
	ctx sdk.Context, chain string, validator string, denom string, epochId int64, epoch types.Epoch,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.EpochStoreKey(chain, validator))
	index := types.GetEpochIndex(denom, []byte(strconv.FormatInt(epochId, 10)))
	b := k.cdc.MustMarshal(&epoch)
	store.Set(index, b)
}

func (k Keeper) GetEpoch(
	ctx sdk.Context, chain string, validator string, denom string, epochId int64,
) (val types.Epoch, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.EpochStoreKey(chain, validator))

	if !k.HasEpoch(ctx, chain, validator, denom, epochId) {
		return types.Epoch{}, false
	}

	index := types.GetEpochIndex(denom, []byte(strconv.FormatInt(epochId, 10)))

	b := store.Get(index)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) SetLastEpochBlock(ctx sdk.Context, chain string, validator string, denom string, lastEpochBlock types.LastEpochBlock) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.EpochStoreKey(chain, validator))
	b := k.cdc.MustMarshal(&lastEpochBlock)
	store.Set(types.GetLastEpochBlockIndex(denom), b)
}

func (k Keeper) GetLastEpochBlock(ctx sdk.Context, chain string, validator string, denom string) (val types.LastEpochBlock, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.EpochStoreKey(chain, validator))

	if !k.HasLastEpochBlock(ctx, chain, validator, denom) {
		return val, false
	}

	b := store.Get(types.GetLastEpochBlockIndex(denom))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) HasLastEpochBlock(ctx sdk.Context, chain string, validator string, denom string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.EpochStoreKey(chain, validator))
	return store.Has(types.GetLastEpochBlockIndex(denom))
}

func (k Keeper) GetAllEpoch(ctx sdk.Context, chain string, validator string, denom string) (list []types.Epoch) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.EpochStoreKey(chain, validator))
	iterator := sdk.KVStorePrefixIterator(store, types.GetEpochIndex(denom, nil))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Epoch
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return

}

func (k Keeper) GetNextRewardsEpochsFromPrevEpochId(
	ctx sdk.Context, chain string, validator string, denom string, epochId int64,
) []*types.Epoch {
	var epoch types.Epoch
	var epochs []*types.Epoch = []*types.Epoch{}
	var nextEpochBlock int64

	epoch, found := k.GetEpoch(ctx, chain, validator, denom, epochId)

	if !found || epoch.NextEpochBlock == types.UndefinedBlockHeight {
		return epochs
	}

	nextEpochBlock = epoch.NextEpochBlock

	for {
		if nextEpochBlock == types.UndefinedBlockHeight {
			break
		}

		epoch, found := k.GetEpoch(ctx, chain, validator, denom, nextEpochBlock)

		if !found || epoch.BlockEnd == types.UndefinedBlockHeight {
			break
		}

		if !epoch.Rewards.IsZero() && !epoch.Staked.IsZero() {
			epochs = append(epochs, &epoch)
		}

		nextEpochBlock = epoch.NextEpochBlock
	}

	return epochs
}
