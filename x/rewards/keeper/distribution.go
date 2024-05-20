package keeper

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/MANTRA-Finance/mantrachain/x/rewards/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetDistributionPairsIdsBytes(ctx sdk.Context) []byte {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.DistributionPairsIdsKey)
	bz := store.Get(byteKey)

	return bz
}

func (k Keeper) GetDistributionPairsIds(ctx sdk.Context) []uint64 {
	bz := k.GetDistributionPairsIdsBytes(ctx)

	if bz == nil {
		return []uint64{}
	}

	buf := bytes.NewReader(bz)

	// Slice to hold the uint64 values
	var pairsIds []uint64

	// Read until the end of the byte slice
	for {
		var num uint64
		err := binary.Read(buf, binary.BigEndian, &num)
		if err != nil {
			// If EOF, break the loop; otherwise, report an error
			if err == bytes.ErrTooLarge {
				break
			}
			fmt.Println("binary.Read failed:", err)
			return pairsIds
		}
		pairsIds = append(pairsIds, num)
	}

	return pairsIds
}

func (k Keeper) SetDistributionPairsIdsBytes(ctx sdk.Context, pairsIds []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.DistributionPairsIdsKey)

	if len(pairsIds) == 0 {
		store.Delete(byteKey)
		return
	}

	store.Set(byteKey, pairsIds)
}

func (k Keeper) SetDistributionPairsIds(ctx sdk.Context, pairsIds []uint64) {
	buf := new(bytes.Buffer)

	// Iterate over the uint64 slice and write each value to the buffer
	for _, pairId := range pairsIds {
		err := binary.Write(buf, binary.BigEndian, pairId)
		if err != nil {
			fmt.Println("binary.Write failed:", err)
			return
		}
	}

	// The bytes slice
	bz := buf.Bytes()
	k.SetDistributionPairsIdsBytes(ctx, bz)
}

func (k Keeper) DistributeRewards(ctx sdk.Context) {
	logger := k.Logger(ctx)
	params := k.GetParams(ctx)

	snapshotsLastDistributedAt := k.GetSnapshotsLastDistributedAt(ctx)
	blockTime := ctx.BlockTime().Unix()
	distributionPeriod := params.DistributionPeriod

	if distributionPeriod > 0 && blockTime < int64(snapshotsLastDistributedAt)+int64(distributionPeriod) {
		return
	}

	pairsIds := k.GetDistributionPairsIds(ctx)

	if len(pairsIds) == 0 {
		pairsIds = k.liquidityKeeper.GetAllPairsIds(ctx)
	}

	if len(pairsIds) == 0 {
		return
	}

	current := pairsIds
	rest := []uint64{}
	if params.PairsCycleMaxCount < uint64(len(pairsIds))-1 {
		current = pairsIds[:params.PairsCycleMaxCount]
		rest = pairsIds[params.PairsCycleMaxCount:]
	}

	for _, pairId := range current {
		err := k.DistributeRewardsForPair(ctx, pairId)
		if err != nil {
			logger.Error("error distributing rewards for pair", "pair_id", pairId, "error", err.Error())
		}
	}

	// The cycle of distribution is over
	if len(rest) == 0 {
		k.SetSnapshotsLastDistributedAt(ctx, uint64(ctx.BlockTime().Unix()))
	}

	k.SetDistributionPairsIds(ctx, rest)
}

func (k Keeper) DistributeRewardsForPair(ctx sdk.Context, pairId uint64) error {
	logger := k.Logger(ctx)
	params := k.GetParams(ctx)
	admin := k.gk.GetAdmin(ctx)

	distributionFeeRate := params.DistributionFeeRate

	pair, found := k.liquidityKeeper.GetPair(ctx, pairId)

	if !found {
		logger.Error("no pair found", "pair_id", pairId)
		return nil
	}

	lastSnapshot, found := k.GetLastSnapshot(ctx, pairId)

	if !found {
		return nil
	}

	if lastSnapshot.Distributed {
		// Create a new snapshot for the pair
		newSnapshot := types.Snapshot{
			Id:          uint64(math.MaxUint64),
			PairId:      pairId,
			Pools:       []*types.SnapshotPool{},
			PoolIdToIdx: map[uint64]uint64{},
			Distributed: false,
		}

		for _, pool := range lastSnapshot.Pools {
			newSnapshot.Pools = append(newSnapshot.Pools, &types.SnapshotPool{
				PoolId:                pool.PoolId,
				CumulativeTotalSupply: pool.CumulativeTotalSupply,
				RewardsPerToken:       sdk.NewDecCoins(),
			})
			newSnapshot.PoolIdToIdx[pool.PoolId] = uint64(len(newSnapshot.Pools) - 1)
		}

		lastSnapshot = newSnapshot
	}

	pairCummulativeTotalSupply := sdk.ZeroDec()

	for _, pool := range lastSnapshot.Pools {
		pairCummulativeTotalSupply = pairCummulativeTotalSupply.Add(pool.CumulativeTotalSupply)
	}

	if pairCummulativeTotalSupply.IsZero() {
		return nil
	}

	swapFeeCollectorAddress := pair.GetSwapFeeCollectorAddress()

	balances := k.bankKeeper.GetAllBalances(ctx, swapFeeCollectorAddress)

	if len(balances) == 0 {
		return nil
	}

	logger.Info("distributing rewards for pair", "pair_id", pairId)

	for _, balance := range balances {
		availableBalance := sdk.NewDecFromInt(balance.Amount)

		if distributionFeeRate.IsPositive() {
			rewardsBalance := availableBalance.Mul(sdk.OneDec().Sub(distributionFeeRate))
			distributionFeeRateCoin := sdk.NewCoin(balance.Denom, availableBalance.Sub(rewardsBalance).TruncateInt())

			if !distributionFeeRateCoin.IsZero() {
				whitelisted := k.gk.WhitelistTransferAccAddresses([]string{swapFeeCollectorAddress.String()}, true)
				err := k.bankKeeper.SendCoins(ctx, swapFeeCollectorAddress, admin, sdk.NewCoins(distributionFeeRateCoin))
				k.gk.WhitelistTransferAccAddresses(whitelisted, false)
				if err != nil {
					return err
				}
			}

			availableBalance = rewardsBalance
			// Update the balance
			balance.Amount = balance.Amount.Sub(distributionFeeRateCoin.Amount)
		}

		for _, pool := range lastSnapshot.Pools {
			if pool.CumulativeTotalSupply.IsZero() {
				continue
			}

			requestedPoolShare := pool.CumulativeTotalSupply.Quo(pairCummulativeTotalSupply)
			rewardAmount := requestedPoolShare.Mul(availableBalance)
			pool.RewardsPerToken = pool.RewardsPerToken.Add(sdk.NewDecCoinFromDec(balance.Denom, rewardAmount.Quo(pool.CumulativeTotalSupply)))
		}

		if !balance.IsZero() {
			whitelisted := k.gk.WhitelistTransferAccAddresses([]string{swapFeeCollectorAddress.String()}, true)
			err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, swapFeeCollectorAddress, types.ModuleName, sdk.NewCoins(balance))
			k.gk.WhitelistTransferAccAddresses(whitelisted, false)
			if err != nil {
				return err
			}
		}

		lastSnapshot.Remaining = lastSnapshot.Remaining.Add(sdk.NewDecCoinFromDec(balance.Denom, sdk.NewDecFromInt(balance.Amount)))
	}

	var distributedAt = ctx.BlockTime()

	lastSnapshot.Distributed = true
	lastSnapshot.DistributedAt = &distributedAt

	// Update the last snapshot
	if lastSnapshot.Id == uint64(math.MaxUint64) {
		// Create a new snapshot
		k.AppendSnapshot(ctx, lastSnapshot)
	} else {
		// Update the last snapshot
		k.SetSnapshot(ctx, lastSnapshot)
	}

	return nil
}
