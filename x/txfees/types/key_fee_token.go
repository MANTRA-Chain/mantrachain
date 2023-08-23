package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// FeeTokenKeyPrefix is the prefix to retrieve all FeeToken
	FeeTokenKeyPrefix = "FeeToken/value/"
)

// FeeTokenKey returns the store key to retrieve a FeeToken from the index fields
func FeeTokenKey(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
