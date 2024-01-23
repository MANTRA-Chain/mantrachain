package types

import (
	"encoding/binary"
)

const (
	// ModuleName defines the module name
	ModuleName = "rewards"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_rewards"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	SnapshotCountKey              = "Snapshot/count/"
	SnapshotStartIdKey            = "Snapshot/startId/"
	SnapshotsLastDistributedAtKey = "Snapshots/lastDistributedAt/"
	DistributionPairsIdsKey       = "Distribution/pairsIds/"
	PurgePairsIdsKey              = "Purge/pairsIds/"
)

var (
	snapshotStoreKey        = "snapshot-store"
	snapshotCountStoreKey   = "snapshot-count-store"
	snapshotStartIdStoreKey = "snapshot-start-id-store"
	Delimiter               = []byte{0x00}
	Placeholder             = []byte{0x01}
)

func SnapshotStoreKey(pairId uint64) []byte {
	var key []byte
	if pairId == 0 {
		key = make([]byte, len(snapshotStoreKey)+len(Delimiter))
		copy(key, snapshotStoreKey)
		copy(key[len(snapshotStoreKey):], Delimiter)
	} else {
		pairIdBytes := make([]byte, binary.MaxVarintLen64)
		binary.PutUvarint(pairIdBytes, pairId)

		key = make([]byte, len(snapshotStoreKey)+len(Delimiter)+len(pairIdBytes)+len(Delimiter))
		copy(key, snapshotStoreKey)
		copy(key[len(snapshotStoreKey):], Delimiter)
		copy(key[len(snapshotStoreKey)+len(Delimiter):], pairIdBytes)
		copy(key[len(snapshotStoreKey)+len(Delimiter)+len(pairIdBytes):], Delimiter)
	}

	return key
}

func SnapshotCountStoreKey(pairId uint64) []byte {
	var key []byte
	if pairId == 0 {
		key = make([]byte, len(snapshotCountStoreKey)+len(Delimiter))
		copy(key, snapshotCountStoreKey)
		copy(key[len(snapshotCountStoreKey):], Delimiter)
	} else {
		pairIdBytes := make([]byte, binary.MaxVarintLen64)
		binary.PutUvarint(pairIdBytes, pairId)

		key = make([]byte, len(snapshotCountStoreKey)+len(Delimiter)+len(pairIdBytes)+len(Delimiter))
		copy(key, snapshotCountStoreKey)
		copy(key[len(snapshotCountStoreKey):], Delimiter)
		copy(key[len(snapshotCountStoreKey)+len(Delimiter):], pairIdBytes)
		copy(key[len(snapshotCountStoreKey)+len(Delimiter)+len(pairIdBytes):], Delimiter)
	}

	return key
}

func SnapshotStartIdStoreKey(pairId uint64) []byte {
	var key []byte
	if pairId == 0 {
		key = make([]byte, len(snapshotStartIdStoreKey)+len(Delimiter))
		copy(key, snapshotStartIdStoreKey)
		copy(key[len(snapshotStartIdStoreKey):], Delimiter)
	} else {
		pairIdBytes := make([]byte, binary.MaxVarintLen64)
		binary.PutUvarint(pairIdBytes, pairId)

		key = make([]byte, len(snapshotStartIdStoreKey)+len(Delimiter)+len(pairIdBytes)+len(Delimiter))
		copy(key, snapshotStartIdStoreKey)
		copy(key[len(snapshotStartIdStoreKey):], Delimiter)
		copy(key[len(snapshotStartIdStoreKey)+len(Delimiter):], pairIdBytes)
		copy(key[len(snapshotStartIdStoreKey)+len(Delimiter)+len(pairIdBytes):], Delimiter)
	}

	return key
}