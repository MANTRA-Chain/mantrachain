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

var (
	ParamsKey             = collections.NewPrefix("p_xfeemarket")
	PrefixDenomMultiplier = []byte{0x01}

	EventTypeTipRefund = "tip_refund"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func KeyDenomMultiplier(denom string) []byte {
	return append(PrefixDenomMultiplier, []byte(denom)...)
}
