package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "sanction"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_sanction"
)

var (
	ParamsKey               = collections.NewPrefix("p_sanction")
	PrefixBlackListAccounts = []byte{0x01}
	PrefixSanctioners       = []byte{0x02}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
