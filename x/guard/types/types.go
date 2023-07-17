package types

import (
	"fmt"
)

type RequiredPrivilegesKind string

// Important: Update ExportGenesis when adding new constants,
// otherwise the data related to them will be lost on chain upgrade.
const (
	RequiredPrivilegesCoin RequiredPrivilegesKind = "coin"
)

func ParseRequiredPrivilegesKind(s string) (c RequiredPrivilegesKind, err error) {
	requiredPrivilegesKind := map[RequiredPrivilegesKind]struct{}{
		RequiredPrivilegesCoin: {},
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
