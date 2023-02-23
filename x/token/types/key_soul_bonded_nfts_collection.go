package types

import "encoding/binary"

var _ binary.ByteOrder

const (
    // SoulBondedNftsCollectionKeyPrefix is the prefix to retrieve all SoulBondedNftsCollection
	SoulBondedNftsCollectionKeyPrefix = "SoulBondedNftsCollection/value/"
)

// SoulBondedNftsCollectionKey returns the store key to retrieve a SoulBondedNftsCollection from the index fields
func SoulBondedNftsCollectionKey(
index string,
) []byte {
	var key []byte
    
    indexBytes := []byte(index)
    key = append(key, indexBytes...)
    key = append(key, []byte("/")...)
    
	return key
}