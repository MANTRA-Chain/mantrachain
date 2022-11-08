package types

import "encoding/binary"

var _ binary.ByteOrder

const (
    // ChainValidatorBridgeKeyPrefix is the prefix to retrieve all ChainValidatorBridge
	ChainValidatorBridgeKeyPrefix = "ChainValidatorBridge/value/"
)

// ChainValidatorBridgeKey returns the store key to retrieve a ChainValidatorBridge from the index fields
func ChainValidatorBridgeKey(
index string,
) []byte {
	var key []byte
    
    indexBytes := []byte(index)
    key = append(key, indexBytes...)
    key = append(key, []byte("/")...)
    
	return key
}