package types

import (
	"fmt"
	"strconv"
)

type RequiredPrivilegesKind string

// Important: Update ExportGenesis when adding new constants,
// otherwise the data related to them will be lost on chain upgrade.
const (
	RequiredPrivilegesCoin  RequiredPrivilegesKind = "coin"
	RequiredPrivilegesAuthz RequiredPrivilegesKind = "authz"
)

func ParseRequiredPrivilegesKind(s string) (c RequiredPrivilegesKind, err error) {
	requiredPrivilegesKind := map[RequiredPrivilegesKind]struct{}{
		RequiredPrivilegesCoin:  {},
		RequiredPrivilegesAuthz: {},
	}

	cap := RequiredPrivilegesKind(s)
	_, ok := requiredPrivilegesKind[cap]
	if !ok {
		return c, fmt.Errorf(`cannot parse:[%s] as required privileges kind`, s)
	}
	return cap, nil
}

func (c RequiredPrivilegesKind) String() string {
	return string(c)
}

func (c RequiredPrivilegesKind) Bytes() []byte {
	return []byte(c)
}

func BinaryStringToBytes(s string) ([]byte, error) {
	var bytes []byte
	for i := 0; i < len(s); i += 8 {
		end := i + 8
		// Ensure not going past the end of the string
		if end > len(s) {
			end = len(s)
		}
		// Parse the binary string to uint64, then cast to byte
		b, err := strconv.ParseUint(s[i:end], 2, 8)
		if err != nil {
			return nil, err
		}
		bytes = append(bytes, byte(b))
	}
	return bytes, nil
}

func BytesToBinaryString(bytes []byte) string {
	var binaryStr string
	for _, b := range bytes {
		// Convert each byte to a binary string
		binaryStr += fmt.Sprintf("%08b", b)
	}
	return binaryStr
}
