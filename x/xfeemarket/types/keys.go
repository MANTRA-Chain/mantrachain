package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "xfeemarket"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_xfeemarket"
)

var ParamsKey = collections.NewPrefix("p_xfeemarket")

func KeyPrefix(p string) []byte {
	return []byte(p)
}
