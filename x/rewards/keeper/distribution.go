package keeper

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"

	utils "github.com/AumegaChain/aumega/types"
	"github.com/AumegaChain/aumega/x/rewards/types"
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
	r := rand.New(rand.NewSource(0))

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

	if len(pairsIds) != 0 {
		cnt := uint64(len(pairsIds))
		if cnt > params.PairsCycleMaxCount {
			cnt = params.PairsCycleMaxCount
		}

		for range make([]struct{}, cnt) {
			pairIdIndex := utils.RandomUint(r, 0, uint64(len(pairsIds)-1))
			pairId := pairsIds[pairIdIndex]
			pairsIds = append(pairsIds[:pairIdIndex], pairsIds[pairIdIndex+1:]...)
			err := k.DistributeRewardsForPair(ctx, pairId)
			if err != nil {
				logger.Error("Error distributing rewards for pair", "pair_id", pairId, "error", err.Error())
			}
		}

		// The cycle of distribution is over
		if len(pairsIds) == 0 {
			k.SetSnapshotsLastDistributedAt(ctx, uint64(ctx.BlockTime().Unix()))
		}

		k.SetDistributionPairsIds(ctx, pairsIds)
	}
}

func (k Keeper) DistributeRewardsForPair(ctx sdk.Context, pairId uint64) error {
	logger := k.Logger(ctx)
	params := k.GetParams(ctx)
	admin := k.gk.GetAdmin(ctx)

	distributionFeeRate := params.DistributionFeeRate

	pair, found := k.liquidityKeeper.GetPair(ctx, pairId)

	if !found {
		logger.Error("No pair found", "pair_id", pairId)
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
				PoolId:          pool.PoolId,
				RewardsPerToken: sdk.NewDecCoins(),
			})
			newSnapshot.PoolIdToIdx[pool.PoolId] = uint64(len(newSnapshot.Pools) - 1)
		}

		lastSnapshot = newSnapshot
	}

	pairCummulativeTotalSupply := sdk.ZeroDec()

	for _, pool := range lastSnapshot.Pools {
		liquidityPool, found := k.liquidityKeeper.GetPool(ctx, pool.PoolId)

		if !found {
			logger.Error("No pool found for pair", "pair_id", pairId, "pool_id", pool.PoolId)
			continue
		}

		if !liquidityPool.Disabled {
			poolCoinSupply := k.liquidityKeeper.GetPoolCoinSupply(ctx, liquidityPool)

			if poolCoinSupply.IsPositive() {
				pairCummulativeTotalSupply = pairCummulativeTotalSupply.Add(sdk.NewDecFromInt(poolCoinSupply))
			}
		}
	}

	if pairCummulativeTotalSupply.IsZero() {
		return nil
	}

	swapFeeCollectorAddress := pair.GetSwapFeeCollectorAddress()

	balances := k.bankKeeper.GetAllBalances(ctx, swapFeeCollectorAddress)

	if len(balances) == 0 {
		return nil
	}

	logger.Info("Distributing rewards for pair", "pair_id", pairId)

	for _, balance := range balances {
		availableBalance := sdk.NewDecFromInt(balance.Amount)

		if distributionFeeRate.IsPositive() {
			rewardsBalance := availableBalance.Mul(sdk.OneDec().Sub(distributionFeeRate))
			distributionFeeRateCoin := sdk.NewCoin(balance.Denom, availableBalance.Sub(rewardsBalance).TruncateInt())

			if !distributionFeeRateCoin.IsZero() {
				whitelisted := k.gk.WhitelistTransferAccAddresses([]string{swapFeeCollectorAddress.String()}, true)
				if err := k.bankKeeper.SendCoins(ctx, swapFeeCollectorAddress, admin, sdk.NewCoins(distributionFeeRateCoin)); err != nil {
					k.gk.WhitelistTransferAccAddresses(whitelisted, false)
					return err
				}
				k.gk.WhitelistTransferAccAddresses(whitelisted, false)
			}

			availableBalance = rewardsBalance
			// Update the balance
			balance.Amount = balance.Amount.Sub(distributionFeeRateCoin.Amount)
		}

		for _, pool := range lastSnapshot.Pools {
			liquidityPool, found := k.liquidityKeeper.GetPool(ctx, pool.PoolId)

			if !found {
				logger.Error("No pool found for pair", "pair_id", pairId, "pool_id", pool.PoolId)
				continue
			}

			if liquidityPool.Disabled {
				continue
			}

			poolCoinSupply := k.liquidityKeeper.GetPoolCoinSupply(ctx, liquidityPool)

			if poolCoinSupply.IsPositive() {
				requestedPoolShare := sdk.NewDecFromInt(poolCoinSupply).Quo(pairCummulativeTotalSupply)
				rewardAmount := requestedPoolShare.Mul(availableBalance)
				pool.RewardsPerToken = pool.RewardsPerToken.Add(sdk.NewDecCoinFromDec(balance.Denom, rewardAmount.Quo(sdk.NewDecFromInt(poolCoinSupply))))
			}
		}

		if !balance.IsZero() {
			whitelisted := k.gk.WhitelistTransferAccAddresses([]string{swapFeeCollectorAddress.String()}, true)
			err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, swapFeeCollectorAddress, types.ModuleName, sdk.NewCoins(balance))
			if err != nil {
				k.gk.WhitelistTransferAccAddresses(whitelisted, false)
				return err
			}
			k.gk.WhitelistTransferAccAddresses(whitelisted, false)
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
