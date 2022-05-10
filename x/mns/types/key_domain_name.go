package types

import "encoding/binary"

var _ binary.ByteOrder

const (
    // DomainNameKeyPrefix is the prefix to retrieve all DomainName
	DomainNameKeyPrefix = "DomainName/value/"
)

// DomainNameKey returns the store key to retrieve a DomainName from the index fields
func DomainNameKey(
index string,
) []byte {
	var key []byte
    
    indexBytes := []byte(index)
    key = append(key, indexBytes...)
    key = append(key, []byte("/")...)
    
	return key
}