package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "tax"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_tax"
)

var ParamsKey = collections.NewPrefix("p_tax")

func KeyPrefix(p string) []byte {
	return []byte(p)
}
