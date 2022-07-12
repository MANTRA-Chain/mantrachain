package types

import (
	"github.com/LimeChain/mantrachain/internal/conv"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	AttributeKeyNftCollectionId      = "nft_collection_id"
	AttributeKeyNftCollectionCreator = "nft_collection_creator"
	AttributeKeyNftsIds              = "nfts_ids"
	AttributeKeyNftId                = "nft_id"
	AttributeKeySigner               = "signer"
	AttributeKeyOwner                = "owner"
	AttributeKeyReceiver             = "receiver"
	AttributeKeyApproved             = "approved"
)

var (
	nftIndex            = "nft-id"
	nftCollectionIndex  = "nft-collection-id"
	nftApprovedAllIndex = "nft-approved-all-index"

	nftStoreKey            = "nft-store"
	nftCollectionStoreKey  = "nft-collection-store"
	nftApprovedStoreKey    = "nft-approved-store"
	nftApprovedAllStoreKey = "nft-approved-all-store"

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
	var key []byte
	if creator == nil {
		key = make([]byte, len(nftCollectionStoreKey)+len(delimiter))
		copy(key, nftCollectionStoreKey)
		copy(key[len(nftCollectionStoreKey):], delimiter)
	} else {
		creator = address.MustLengthPrefix(creator)

		key = make([]byte, len(nftCollectionStoreKey)+len(delimiter)+len(creator)+len(delimiter))
		copy(key, nftCollectionStoreKey)
		copy(key[len(nftCollectionStoreKey):], delimiter)
		copy(key[len(nftCollectionStoreKey)+len(delimiter):], creator)
		copy(key[len(nftCollectionStoreKey)+len(delimiter)+len(creator):], delimiter)
	}

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

func NftApprovedStoreKey(collectionIndex []byte) []byte {
	key := make([]byte, len(nftApprovedStoreKey)+len(delimiter)+len(collectionIndex)+len(delimiter))
	copy(key, nftApprovedStoreKey)
	copy(key[len(nftApprovedStoreKey):], delimiter)
	copy(key[len(nftApprovedStoreKey)+len(delimiter):], collectionIndex)
	copy(key[len(nftApprovedStoreKey)+len(delimiter)+len(collectionIndex):], delimiter)
	return key
}

func NftApprovedAllStoreKey() []byte {
	key := make([]byte, len(nftApprovedAllStoreKey)+len(delimiter))
	copy(key, nftApprovedAllStoreKey)
	copy(key[len(nftApprovedAllStoreKey):], delimiter)
	return key
}

func GetNftApprovedAllIndex(owner sdk.AccAddress) []byte {
	owner = address.MustLengthPrefix(owner)

	key := make([]byte, len(nftApprovedAllIndex)+len(delimiter)+len(owner)+len(delimiter))
	copy(key, nftApprovedAllIndex)
	copy(key[len(nftApprovedAllIndex):], delimiter)
	copy(key[len(nftApprovedAllIndex)+len(delimiter):], owner)
	copy(key[len(nftApprovedAllIndex)+len(delimiter)+len(owner):], delimiter)
	return key
}
