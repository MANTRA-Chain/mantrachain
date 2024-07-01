package keeper

import (
	"math"

	sdkmath "cosmossdk.io/math"
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

func (h Hooks) OnProvideLiquidity(ctx sdk.Context, receiver sdk.Address, pairId uint64, poolId uint64, poolCoin sdk.Coin) error {
	var err error
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
				CoinSupply:      pool.CoinSupply,
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
	poolIdx, found := lastSnapshot.PoolIdToIdx[poolId]
	if !found {
		// Create a new snapshot pool
		poolIdx = uint64(len(lastSnapshot.Pools))
		lastSnapshot.Pools = append(lastSnapshot.Pools, &types.SnapshotPool{
			PoolId:          poolId,
			CoinSupply:      sdkmath.LegacyNewDec(0),
			RewardsPerToken: sdk.NewDecCoins(),
		})
		lastSnapshot.PoolIdToIdx[poolId] = poolIdx
	}

	// Update the last snapshot
	lastSnapshot.Pools[poolIdx].CoinSupply = lastSnapshot.Pools[poolIdx].CoinSupply.Add(sdkmath.LegacyNewDecFromInt(poolCoin.Amount))

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
			Balances:              sdk.NewCoins(),
			PoolIdToIdx:           map[uint64]uint64{},
			OwedRewards:           sdk.NewDecCoins(),
		})
		provider.PairIdToIdx[pairId] = pairIdx
	} else {
		provider, err = h.k.CalculateRewards(ctx, pairId, provider, &types.ClaimParams{
			IsQuery: false,
		})
		if err != nil {
			h.k.logger.Error("failed to calculate rewards", "error", err)
			return err
		}
	}

	poolIdx, found = provider.Pairs[pairIdx].PoolIdToIdx[poolId]

	if !found {
		// Create a new provider pair pool
		poolIdx = uint64(len(provider.Pairs[pairIdx].Balances))
		provider.Pairs[pairIdx].Balances = append(provider.Pairs[pairIdx].Balances, sdk.NewCoin(poolCoin.Denom, sdkmath.ZeroInt()))
		provider.Pairs[pairIdx].PoolIdToIdx[poolId] = poolIdx
	}

	// Update the provider pair
	provider.Pairs[pairIdx].Balances[poolIdx] = provider.Pairs[pairIdx].Balances[poolIdx].Add(sdk.NewCoin(poolCoin.Denom, poolCoin.Amount))

	blockTime := ctx.BlockTime()
	provider.Pairs[pairIdx].LastDepositTime = &blockTime

	// Update the provider
	h.k.SetProvider(ctx, provider)

	return nil
}

func (h Hooks) OnWithdrawLiquidity(ctx sdk.Context, receiver sdk.Address, pairId uint64, poolId uint64, poolCoin sdk.Coin) error {
	lastSnapshot, found := h.k.GetLastSnapshot(ctx, pairId)

	if !found {
		h.k.logger.Error("no snapshot found for pair", "pair_id", pairId)
		return types.ErrSnapshotNotFound
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
				CoinSupply:      pool.CoinSupply,
				RewardsPerToken: sdk.NewDecCoins(),
			})
			newSnapshot.PoolIdToIdx[pool.PoolId] = uint64(len(newSnapshot.Pools) - 1)
		}

		lastSnapshot = newSnapshot
	}

	poolIdx, found := lastSnapshot.PoolIdToIdx[poolId]

	if !found {
		h.k.logger.Error("No snapshot pool found for pair", "pool_id", poolId, "pair_id", pairId)
		return types.ErrSnapshotPoolNotFound
	}

	// Update the last snapshot
	lastSnapshot.Pools[poolIdx].CoinSupply = lastSnapshot.Pools[poolIdx].CoinSupply.Sub(sdkmath.LegacyNewDecFromInt(poolCoin.Amount))

	if lastSnapshot.Id == uint64(math.MaxUint64) {
		// Create a new snapshot
		h.k.AppendSnapshot(ctx, lastSnapshot)
	} else {
		// Update the last snapshot
		h.k.SetSnapshot(ctx, lastSnapshot)
	}

	provider, found := h.k.GetProvider(ctx, receiver.String())

	if !found {
		h.k.logger.Error("no provider found for address", "receiver", receiver.String())
		return types.ErrProviderNotFound
	}

	pairIdx, found := provider.PairIdToIdx[pairId]

	if !found {
		h.k.logger.Error("No provider pair found for pair", "pair_id", pairId)
		return types.ErrProviderPairNotFound
	}

	// Update the provider pair
	provider, err := h.k.CalculateRewards(ctx, pairId, provider, &types.ClaimParams{
		IsQuery: false,
	})
	if err != nil {
		h.k.logger.Error("failed to calculate rewards", "error", err)
		return err
	}

	if !found {
		h.k.logger.Error("No provider pool found for pair", "pool_id", poolId, "pair_id", pairId)
		return types.ErrProviderPoolNotFound
	}

	poolIdx, found = provider.Pairs[pairIdx].PoolIdToIdx[poolId]

	if !found {
		h.k.logger.Error("No provider pool found for pair", "pool_id", poolId, "pair_id", pairId)
		return types.ErrProviderPoolNotFound
	}

	if provider.Pairs[pairIdx].Balances[poolIdx].IsLT(sdk.NewCoin(poolCoin.Denom, poolCoin.Amount)) {
		h.k.logger.Error("balance mismatch", "balance", provider.Pairs[pairIdx].Balances[poolIdx], "pool_coin", poolCoin)
		return types.ErrBalanceMismatch
	}

	// Update the provider pair pool
	provider.Pairs[pairIdx].Balances[poolIdx] = provider.Pairs[pairIdx].Balances[poolIdx].Sub(sdk.NewCoin(poolCoin.Denom, poolCoin.Amount))

	// Update the provider
	h.k.SetProvider(ctx, provider)

	return nil
}
