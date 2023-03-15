package types

import (
	"fmt"
)

type RequiredPrivilegesKind string
type LockedKind string

const (
	RequiredPrivilegesCoin RequiredPrivilegesKind = "coin"
	LockedCoin             LockedKind             = "coin"
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

func ParseLockedKind(s string) (c LockedKind, err error) {
	lockedKind := map[LockedKind]struct{}{
		LockedCoin: {},
	}

	cap := LockedKind(s)
	_, ok := lockedKind[cap]
	if !ok {
		return c, fmt.Errorf(`cannot parse:[%s] as locked kind`, s)
	}
	return cap, nil
}

func (c LockedKind) String() string {
	return string(c)
}

func (c LockedKind) Bytes() []byte {
	return []byte(c)
}
