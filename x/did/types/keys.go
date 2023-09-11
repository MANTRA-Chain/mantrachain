package types

var (
	// DidDocumentKey prefix for each key to a DidDocument
	DidDocumentKey = []byte{0x61}
	// DidMetadataKey prefix for each key of a DidMetadata
	DidMetadataKey = []byte{0x62}

	Delimiter   = []byte{0x00}
	Placeholder = []byte{0x01}
)

const (
	// ModuleName defines the module name
	ModuleName = "did"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// DidChainPrefix defines the did prefix for this chain
	DidChainPrefix = "did:cosmos:net:"

	// DidKeyPrefix defines the did key prefix
	DidKeyPrefix = "did:cosmos:key:"

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_capability_did"
)
