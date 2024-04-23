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

func (k Keeper) GetPurgePairsIdsBytes(ctx sdk.Context) []byte {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.PurgePairsIdsKey)
	bz := store.Get(byteKey)

	return bz
}

func (k Keeper) GetPurgePairsIds(ctx sdk.Context) []uint64 {
	bz := k.GetPurgePairsIdsBytes(ctx)

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

func (k Keeper) SetPurgePairsIdsBytes(ctx sdk.Context, pairsIds []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.PurgePairsIdsKey)

	if len(pairsIds) == 0 {
		store.Delete(byteKey)
		return
	}

	store.Set(byteKey, pairsIds)
}

func (k Keeper) SetPurgePairsIds(ctx sdk.Context, pairsIds []uint64) {
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
	k.SetPurgePairsIdsBytes(ctx, bz)
}

func (k Keeper) PurgeSnapshots(ctx sdk.Context) {
	logger := k.Logger(ctx)
	params := k.GetParams(ctx)

	pairsIds := k.GetPurgePairsIds(ctx)

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
		err := k.PurgeSnapshotsForPair(ctx, pairId)
		if err != nil {
			logger.Error("fail to purge snapshots for pair", "pair_id", pairId, "error", err.Error())
		}
	}

	k.SetPurgePairsIds(ctx, rest)
}

func (k Keeper) PurgeSnapshotsForPair(ctx sdk.Context, pairId uint64) error {
	logger := k.Logger(ctx)
	conf := k.GetParams(ctx)
	maxSnapshotsCount := conf.MaxSnapshotsCount
	lastSnapshot, found := k.GetLastSnapshot(ctx, pairId)

	if !found {
		return nil
	}

	snapshotStartId, found := k.GetSnapshotStartId(ctx, pairId)

	if !found {
		snapshotStartId = types.SnapshotStartId{
			PairId:     pairId,
			SnapshotId: 0,
		}
	}

	if maxSnapshotsCount == math.MaxUint64 || lastSnapshot.Id < maxSnapshotsCount || lastSnapshot.Id-snapshotStartId.SnapshotId < maxSnapshotsCount {
		return nil
	}

	snapshotEndId := lastSnapshot.Id - maxSnapshotsCount

	if conf.MaxPurgedRangeLength > 0 && snapshotEndId-snapshotStartId.SnapshotId > conf.MaxPurgedRangeLength {
		snapshotEndId = snapshotStartId.SnapshotId + conf.MaxPurgedRangeLength
	}

	admin := k.gk.GetAdmin(ctx)
	remaining := sdk.NewCoins()

	logger.Info("purge snapshots for pair", "pair_id", pairId, "snapshots_count", snapshotEndId-snapshotStartId.SnapshotId+1)

	for i := snapshotStartId.SnapshotId; i <= snapshotEndId; i++ {
		snapshot, found := k.GetSnapshot(ctx, pairId, i)

		if !found {
			logger.Error("no snapshot found for pair", "pair_id", pairId, "snapshot_id", i)
			continue
		}

		if snapshot.Distributed {
			for _, decCoin := range snapshot.Remaining {
				if decCoin.Amount.IsPositive() {
					remaining = remaining.Add(sdk.NewCoin(decCoin.Denom, decCoin.Amount.TruncateInt()))
				}
			}
		}

		k.RemoveSnapshot(ctx, pairId, i)
	}

	if !remaining.IsZero() {
		// TODO: maybe we should keep the remaining coins in the rewards module account and redistribute them later?
		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, admin, remaining)
		if err != nil {
			return err
		}
	}

	snapshotStartId.SnapshotId = snapshotEndId + 1

	k.SetSnapshotStartId(ctx, snapshotStartId)

	return nil
}
