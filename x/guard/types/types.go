package types

import (
	"fmt"
)

type RequiredPrivilegesKind string

const (
	Coin RequiredPrivilegesKind = "coin"
)

func (c RequiredPrivilegesKind) String() string {
	return string(c)
}

func (c RequiredPrivilegesKind) Bytes() []byte {
	return []byte(c)
}

func ParseRequiredPrivilegesKind(s string) (c RequiredPrivilegesKind, err error) {
	requiredPrivilegesKind := map[RequiredPrivilegesKind]struct{}{
		Coin: {},
	}

	cap := RequiredPrivilegesKind(s)
	_, ok := requiredPrivilegesKind[cap]
	if !ok {
		return c, fmt.Errorf(`cannot parse:[%s] as required privileges kind`, s)
	}
	return cap, nil
}
