package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// AccPermKeyPrefix is the prefix to retrieve all AccPerm
	AccPermKeyPrefix = "AccPerm/value/"
)

// AccPermKey returns the store key to retrieve a AccPerm from the index fields
func AccPermKey(
	id string,
) []byte {
	var key []byte

	idBytes := []byte(id)
	key = append(key, idBytes...)
	key = append(key, []byte("/")...)

	return key
}
