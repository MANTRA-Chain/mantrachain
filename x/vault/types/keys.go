package types

const (
	// ModuleName defines the module name
	ModuleName = "vault"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_vault"
)

var (
	nftStakeStoreKey = "nft-stake-store"

	delimiter = []byte{0x00}
)

func NftStakeStoreKey(marketplaceIndex []byte, collectionIndex []byte) []byte {
	key := make([]byte, len(nftStakeStoreKey)+len(delimiter)+len(marketplaceIndex)+len(delimiter)+len(collectionIndex)+len(delimiter))
	copy(key, nftStakeStoreKey)
	copy(key[len(nftStakeStoreKey):], delimiter)
	copy(key[len(nftStakeStoreKey)+len(delimiter):], marketplaceIndex)
	copy(key[len(nftStakeStoreKey)+len(delimiter)+len(marketplaceIndex):], delimiter)
	copy(key[len(nftStakeStoreKey)+len(delimiter)+len(marketplaceIndex)+len(delimiter):], collectionIndex)
	copy(key[len(nftStakeStoreKey)+len(delimiter)+len(marketplaceIndex)+len(delimiter)+len(collectionIndex):], delimiter)
	return key
}

func KeyPrefix(p string) []byte {
	return []byte(p)
}
