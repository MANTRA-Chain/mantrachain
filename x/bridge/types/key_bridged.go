package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// BridgedKeyPrefix is the prefix to retrieve all Bridged
	BridgedKeyPrefix = "Bridged/value/"
)

// BridgedKey returns the store key to retrieve a Bridged from the index fields
func BridgedKey(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
