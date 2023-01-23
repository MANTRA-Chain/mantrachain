package types

import "encoding/binary"

var _ binary.ByteOrder

const (
    // AccPermKeyPrefix is the prefix to retrieve all AccPerm
	AccPermKeyPrefix = "AccPerm/value/"
)

// AccPermKey returns the store key to retrieve a AccPerm from the index fields
func AccPermKey(
cat string,
) []byte {
	var key []byte
    
    catBytes := []byte(cat)
    key = append(key, catBytes...)
    key = append(key, []byte("/")...)
    
	return key
}