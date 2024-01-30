package keeper

import (
	"math"

	liquiditytypes "github.com/MANTRA-Finance/mantrachain/x/liquidity/types"
	"github.com/MANTRA-Finance/mantrachain/x/rewards/types"
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

func (h Hooks) OnProvideLiquidity(ctx sdk.Context, receiver sdk.Address, pairId uint64, poolId uint64, poolCoin sdk.Coin) {
	lastSnapshot, found := h.k.GetLastSnapshot(ctx, pairId)
	lastClaimedSnapshotId := uint64(math.MaxUint64)

	if !found {
		// Create the first snapshot for the pair
		lastSnapshot = types.Snapshot{
			Id:          uint64(math.MaxUint64),
			PairId:      pairId,
			Pools:       []*types.SnapshotPool{},
			PoolIdToIdx: map[uint64]uint64{},
			Distributed: false,
		}
	} else if lastSnapshot.Distributed {
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

		lastClaimedSnapshotId = lastSnapshot.Id
		lastSnapshot = newSnapshot
	} else if lastSnapshot.Id > 0 {
		lastClaimedSnapshotId = lastSnapshot.Id - 1
	}

	// Update the snapshot pools
	_, found = lastSnapshot.PoolIdToIdx[poolId]
	if !found {
		// Create a new snapshot pool
		poolIdx := uint64(len(lastSnapshot.Pools))
		lastSnapshot.Pools = append(lastSnapshot.Pools, &types.SnapshotPool{
			PoolId:          poolId,
			RewardsPerToken: sdk.NewDecCoins(),
		})
		lastSnapshot.PoolIdToIdx[poolId] = poolIdx
	}

	if lastSnapshot.Id == uint64(math.MaxUint64) {
		// Create a new snapshot
		h.k.AppendSnapshot(ctx, lastSnapshot)
	} else {
		// Update the last snapshot
		h.k.SetSnapshot(ctx, lastSnapshot)
	}

	provider, found := h.k.GetProvider(ctx, receiver.String())

	if !found {
		// Create a new provider
		provider = types.Provider{
			Index:       receiver.String(),
			Pairs:       []*types.ProviderPair{},
			PairIdToIdx: map[uint64]uint64{},
		}
	}

	pairIdx, found := provider.PairIdToIdx[pairId]

	if !found {
		// Create a new provider pair
		pairIdx = uint64(len(provider.Pairs))
		provider.Pairs = append(provider.Pairs, &types.ProviderPair{
			PairId:                pairId,
			LastClaimedSnapshotId: lastClaimedSnapshotId,
			OwedRewards:           sdk.NewDecCoins(),
		})
		provider.PairIdToIdx[pairId] = pairIdx
	} else {
		provider = h.k.CalculateRewards(ctx, receiver.String(), pairId, provider, &types.ClaimParams{IsQuery: false})
	}

	blockTime := ctx.BlockTime()
	// Update the provider pair
	provider.Pairs[pairIdx].LastDepositTime = &blockTime

	// Update the provider
	h.k.SetProvider(ctx, provider)
}

func (h Hooks) OnWithdrawLiquidity(ctx sdk.Context, receiver sdk.Address, pairId uint64, poolId uint64, poolCoin sdk.Coin) {
	logger := h.k.Logger(ctx)
	lastSnapshot, found := h.k.GetLastSnapshot(ctx, pairId)

	if !found {
		logger.Error("no snapshot found for pair", "pair_id", pairId)
		return
	} else if lastSnapshot.Distributed {
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

	if lastSnapshot.Id == uint64(math.MaxUint64) {
		// Create a new snapshot
		h.k.AppendSnapshot(ctx, lastSnapshot)
	} else {
		// Update the last snapshot
		h.k.SetSnapshot(ctx, lastSnapshot)
	}

	provider, found := h.k.GetProvider(ctx, receiver.String())

	if !found {
		logger.Error("no provider found for address", "receiver", receiver.String())
		return
	}

	// Update the provider pair
	provider = h.k.CalculateRewards(ctx, receiver.String(), pairId, provider, &types.ClaimParams{IsQuery: false})

	// Update the provider
	h.k.SetProvider(ctx, provider)
}
