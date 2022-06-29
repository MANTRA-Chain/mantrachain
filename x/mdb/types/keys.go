package types

import (
	"github.com/LimeChain/mantrachain/internal/conv"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

var (
	nftCollectionIndex    = "nft-collection-id"
	nftCollectionStoreKey = "nft-collection-store"
	nftIndex              = "nft-id"
	nftStoreKey           = "nft-store"

	delimiter = []byte{0x00}
)

const (
	// ModuleName defines the module name
	ModuleName = "mdb"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_mdb"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func GetNftCollectionIndex(creator sdk.AccAddress, id string) []byte {
	creator = address.MustLengthPrefix(creator)
	idBz := conv.UnsafeStrToBytes(id)

	key := make([]byte, len(nftCollectionIndex)+len(delimiter)+len(creator)+len(delimiter)+len(idBz)+len(delimiter))
	copy(key, nftCollectionIndex)
	copy(key[len(nftCollectionIndex):], delimiter)
	copy(key[len(nftCollectionIndex)+len(delimiter):], creator)
	copy(key[len(nftCollectionIndex)+len(delimiter)+len(creator):], delimiter)
	copy(key[len(nftCollectionIndex)+len(delimiter)+len(creator)+len(delimiter):], idBz)
	copy(key[len(nftCollectionIndex)+len(delimiter)+len(creator)+len(delimiter)+len(idBz):], delimiter)
	return key
}

func NftCollectionStoreKey(creator sdk.AccAddress) []byte {
	creator = address.MustLengthPrefix(creator)

	key := make([]byte, len(nftCollectionStoreKey)+len(delimiter)+len(creator)+len(delimiter))
	copy(key, nftCollectionStoreKey)
	copy(key[len(nftCollectionStoreKey):], delimiter)
	copy(key[len(nftCollectionStoreKey)+len(delimiter):], creator)
	copy(key[len(nftCollectionStoreKey)+len(delimiter)+len(creator):], delimiter)
	return key
}

func GetNftIndex(collectionIndex []byte, id string) []byte {
	idBz := conv.UnsafeStrToBytes(id)
	key := make([]byte, len(nftIndex)+len(delimiter)+len(collectionIndex)+len(delimiter)+len(idBz)+len(delimiter))
	copy(key, nftIndex)
	copy(key[len(nftIndex):], delimiter)
	copy(key[len(nftIndex)+len(delimiter):], collectionIndex)
	copy(key[len(nftIndex)+len(delimiter)+len(collectionIndex):], delimiter)
	copy(key[len(nftIndex)+len(delimiter)+len(collectionIndex)+len(delimiter):], id)
	copy(key[len(nftIndex)+len(delimiter)+len(collectionIndex)+len(delimiter)+len(id):], delimiter)
	return key
}

func NftStoreKey(collectionIndex []byte) []byte {
	key := make([]byte, len(nftStoreKey)+len(delimiter)+len(collectionIndex)+len(delimiter))
	copy(key, nftStoreKey)
	copy(key[len(nftStoreKey):], delimiter)
	copy(key[len(nftStoreKey)+len(delimiter):], collectionIndex)
	copy(key[len(nftStoreKey)+len(delimiter)+len(collectionIndex):], delimiter)
	return key
}
