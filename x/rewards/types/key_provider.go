package types

import "encoding/binary"

var _ binary.ByteOrder

const (
    // ProviderKeyPrefix is the prefix to retrieve all Provider
	ProviderKeyPrefix = "Provider/value/"
)

// ProviderKey returns the store key to retrieve a Provider from the index fields
func ProviderKey(
index string,
) []byte {
	var key []byte
    
    indexBytes := []byte(index)
    key = append(key, indexBytes...)
    key = append(key, []byte("/")...)
    
	return key
}