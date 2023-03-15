package types

import "encoding/binary"

var _ binary.ByteOrder

const (
    // RequiredPrivilegesKeyPrefix is the prefix to retrieve all RequiredPrivileges
	RequiredPrivilegesKeyPrefix = "RequiredPrivileges/value/"
)

// RequiredPrivilegesKey returns the store key to retrieve a RequiredPrivileges from the index fields
func RequiredPrivilegesKey(
index string,
) []byte {
	var key []byte
    
    indexBytes := []byte(index)
    key = append(key, indexBytes...)
    key = append(key, []byte("/")...)
    
	return key
}