package keeper

import (
	"math"

	"github.com/AumegaChain/aumega/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) CalculateRewards(ctx sdk.Context, receiver string, pairId uint64, provider types.Provider, params *types.ClaimParams) types.Provider {
	logger := k.Logger(ctx)
	conf := k.GetParams(ctx)
	minDepositTime := conf.MinDepositTime
	startClaimedSnapshotId := uint64(0)
	endClaimedSnapshotId := uint64(0)

	receiverAcc, err := sdk.AccAddressFromBech32(receiver)
	if err != nil {
		return types.Provider{}
	}

	if params != nil && params.EndClaimedSnapshotId != nil {
		endClaimedSnapshotId = *params.EndClaimedSnapshotId
	} else {
		endClaimedSnapshotId = k.GetEndClaimedSnapshotId(ctx, pairId)
	}

	pairIdx, found := provider.PairIdToIdx[pairId]
	if !found {
		return provider
	}

	providerPair := provider.Pairs[pairIdx]

	if providerPair == nil {
		return provider
	}

	if params != nil && params.StartClaimedSnapshotId != nil {
		startClaimedSnapshotId = *params.StartClaimedSnapshotId
	} else {
		startClaimedSnapshotId = providerPair.LastClaimedSnapshotId

		// If LastClaimedSnapshotId is zero then the provider has never claimed rewards
		if startClaimedSnapshotId == uint64(math.MaxUint64) {
			startClaimedSnapshotId = 0
		} else {
			startClaimedSnapshotId++
		}

		snapshotStartId, found := k.GetSnapshotStartId(ctx, pairId)

		if !found {
			snapshotStartId = types.SnapshotStartId{
				PairId:     pairId,
				SnapshotId: 0,
			}
		}

		// If we purge the snapshots then we need to start from the first available snapshot if it is greater than the last claimed snapshot
		if snapshotStartId.SnapshotId > startClaimedSnapshotId {
			startClaimedSnapshotId = snapshotStartId.SnapshotId
		}
	}

	// If the provider has never claimed rewards then we need to start from the first snapshot
	snapshots := k.GetSnapshotsInRange(ctx, pairId, startClaimedSnapshotId, endClaimedSnapshotId)

	for _, snapshot := range snapshots {
		if !snapshot.Distributed {
			continue
		}

		if minDepositTime > 0 && snapshot.DistributedAt.Unix() < providerPair.LastDepositTime.Unix()+int64(minDepositTime) {
			continue
		}

		for _, pool := range snapshot.Pools {
			liquidityPool, found := k.liquidityKeeper.GetPool(ctx, pool.PoolId)

			if !found {
				logger.Error("No pool found for pair", "pair_id", pairId, "pool_id", pool.PoolId)
				continue
			}

			if liquidityPool.Disabled {
				continue
			}

			balance := k.bankKeeper.GetBalance(ctx, receiverAcc, liquidityPool.PoolCoinDenom)

			if balance.IsPositive() {
				for _, rewardPerToken := range pool.RewardsPerToken {
					reward := sdk.NewDecCoinFromDec(rewardPerToken.Denom, rewardPerToken.Amount.Mul(sdk.NewDecFromInt(balance.Amount)))
					providerPair.OwedRewards = providerPair.OwedRewards.Add(reward)

					if params != nil && !params.IsQuery {
						snapshot.Remaining = snapshot.Remaining.Sub(sdk.NewDecCoins(reward))
						k.SetSnapshot(ctx, snapshot)
					}
				}
			}
		}
	}

	provider.Pairs[pairIdx].LastClaimedSnapshotId = endClaimedSnapshotId

	return provider
}
