package keeper

import (
	"strconv"

	"github.com/LimeChain/mantrachain/x/vault/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) InitEpoch(ctx sdk.Context, chain string, validator string, bh int64) {
	newEpoch := types.Epoch{
		BlockStart:     bh,
		BlockEnd:       types.UndefinedBlockHeight,
		Rewards:        nil,
		Staked:         sdk.NewDec(0),
		PrevEpochBlock: types.UndefinedBlockHeight,
		NextEpochBlock: types.UndefinedBlockHeight,
		StartAt:        ctx.BlockHeader().Time.Unix(),
	}

	k.SetEpoch(ctx, chain, validator, bh, newEpoch)
	k.SetLastEpochBlock(ctx, chain, validator, types.LastEpochBlock{
		BlockHeight: bh,
	})
}

func (k Keeper) SetEpochEndNative(
	ctx sdk.Context,
	chain string,
	validator string,
	bh int64,
	lastEpochBlockHeight int64,
	minEpochWithdraw string,
) error {
	var parsed *sdk.Coin = nil
	lastEpoch, found := k.GetEpoch(ctx, chain, validator, lastEpochBlockHeight)

	if !found {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "last epoch not found %s", lastEpochBlockHeight)
	}

	if minEpochWithdraw != "" {
		minEpochWithdrawCoin, err := sdk.ParseCoinNormalized(minEpochWithdraw)

		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "min epoch withdraw is invalid %s", minEpochWithdraw)
		}

		parsed = &minEpochWithdrawCoin
	}

	de := NewDistributionExecutor(ctx, k.ac, k.sk, k.dk)
	rewards, err := de.GetDelegationRewards(validator)

	if err != nil {
		return err
	}

	// TODO: sum all rewards before comparison
	for _, reward := range rewards {
		if (parsed == nil && reward.Amount.GTE(sdk.NewDec(1))) ||
			(reward.Denom == parsed.Denom && reward.Amount.GTE(sdk.NewDecFromInt(parsed.Amount))) {
			withdrawn, err := de.WithdrawDelegationRewards(validator)

			if err != nil {
				return err
			}

			lastEpoch.Rewards = withdrawn

			break
		}
	}

	newEpoch := types.Epoch{
		BlockStart:     bh,
		BlockEnd:       types.UndefinedBlockHeight,
		PrevEpochBlock: lastEpochBlockHeight,
		NextEpochBlock: types.UndefinedBlockHeight,
		StartAt:        ctx.BlockHeader().Time.Unix(),
	}

	se := NewStakingExecutor(ctx, k.ac, k.bk, k.sk)
	staked, err := se.GetDelegatorDelegation(validator)

	if err != nil {
		return err
	}

	newEpoch.Staked = staked
	lastEpoch.BlockEnd = bh
	lastEpoch.NextEpochBlock = bh
	lastEpoch.EndAt = ctx.BlockHeader().Time.Unix()

	k.SetEpoch(ctx, chain, validator, lastEpochBlockHeight, lastEpoch)
	k.SetEpoch(ctx, chain, validator, bh, newEpoch)
	k.SetLastEpochBlock(ctx, chain, validator, types.LastEpochBlock{
		BlockHeight: bh,
	})

	return nil
}

func (k Keeper) HasEpoch(
	ctx sdk.Context, chain string, validator string, epochId int64,
) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.EpochStoreKey(chain))
	index := types.GetEpochIndex(validator, []byte(strconv.FormatInt(epochId, 10)))
	return store.Has(index)
}

func (k Keeper) SetEpoch(
	ctx sdk.Context, chain string, validator string, epochId int64, epoch types.Epoch,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.EpochStoreKey(chain))
	index := types.GetEpochIndex(validator, []byte(strconv.FormatInt(epochId, 10)))
	b := k.cdc.MustMarshal(&epoch)
	store.Set(index, b)
}

func (k Keeper) GetEpoch(
	ctx sdk.Context, chain string, validator string, epochId int64,
) (val types.Epoch, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.EpochStoreKey(chain))

	if !k.HasEpoch(ctx, chain, validator, epochId) {
		return types.Epoch{}, false
	}

	index := types.GetEpochIndex(validator, []byte(strconv.FormatInt(epochId, 10)))

	b := store.Get(index)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) SetLastEpochBlock(ctx sdk.Context, chain string, validator string, lastEpochBlock types.LastEpochBlock) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.EpochStoreKey(chain))
	b := k.cdc.MustMarshal(&lastEpochBlock)
	store.Set(types.GetLastEpochBlockIndex(validator), b)
}

func (k Keeper) GetLastEpochBlock(ctx sdk.Context, chain string, validator string) (val types.LastEpochBlock, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.EpochStoreKey(chain))

	if !k.HasLastEpochBlock(ctx, chain, validator) {
		return val, false
	}

	b := store.Get(types.GetLastEpochBlockIndex(validator))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) HasLastEpochBlock(ctx sdk.Context, chain string, validator string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.EpochStoreKey(chain))
	return store.Has(types.GetLastEpochBlockIndex(validator))
}

func (k Keeper) GetAllEpoch(ctx sdk.Context, chain string, validator string) (list []types.Epoch) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.EpochStoreKey(chain))
	iterator := sdk.KVStorePrefixIterator(store, types.GetEpochIndex(validator, nil))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Epoch
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return

}

func (k Keeper) GetNextRewardsEpochsFromPrevEpochId(
	ctx sdk.Context, chain string, validator string, epochId int64,
) []*types.Epoch {
	var epoch types.Epoch
	var epochs []*types.Epoch = []*types.Epoch{}
	var nextEpochBlock int64

	epoch, found := k.GetEpoch(ctx, chain, validator, epochId)

	if !found || epoch.NextEpochBlock == types.UndefinedBlockHeight {
		return epochs
	}

	nextEpochBlock = epoch.NextEpochBlock

	for {
		if nextEpochBlock == types.UndefinedBlockHeight {
			break
		}

		epoch, found := k.GetEpoch(ctx, chain, validator, nextEpochBlock)

		if !found || epoch.BlockEnd == types.UndefinedBlockHeight {
			break
		}

		if len(epoch.Rewards) > 0 && !epoch.Staked.IsZero() {
			epochs = append(epochs, &epoch)
		}

		nextEpochBlock = epoch.NextEpochBlock
	}

	return epochs
}
