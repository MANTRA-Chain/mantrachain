package types

const (
	// ModuleName defines the module name
	ModuleName = "did"

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// DidChainPrefix defines the did prefix for this chain
	DidChainPrefix = "did:cosmos:net:"

	// DidKeyPrefix defines the did key prefix
	DidKeyPrefix = "did:cosmos:key:"
)

var (
	ParamsKey = []byte("p_did")

	// DidDocumentKey prefix for each key to a DidDocument
	DidDocumentKey = []byte{0x61}
	// DidMetadataKey prefix for each key of a DidMetadata
	DidMetadataKey = []byte{0x62}

	Delimiter   = []byte{0x00}
	Placeholder = []byte{0x01}
)
