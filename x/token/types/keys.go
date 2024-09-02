package types

import (
	"github.com/MANTRA-Finance/mantrachain/internal/conv"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName defines the module name
	ModuleName = "token"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	AttributeKeyNftCollectionId      = "nft_collection_id"
	AttributeKeyNftCollectionCreator = "nft_collection_creator"
	AttributeKeyNftsIds              = "nfts_ids"
	AttributeKeyNftsCount            = "nfts_count"
	AttributeKeyNftId                = "nft_id"
	AttributeKeySigner               = "signer"
	AttributeKeyOwner                = "owner"
	AttributeKeyReceiver             = "receiver"
	AttributeKeyApproved             = "approved"
)

var (
	ParamsKey = []byte("p_token")

	nftIndex            = "nft-id"
	nftCollectionIndex  = "nft-collection-id"
	nftApprovedAllIndex = "nft-approved-all-index"

	nftStoreKey            = "nft-store"
	nftCollectionStoreKey  = "nft-collection-store"
	nftApprovedStoreKey    = "nft-approved-store"
	nftApprovedAllStoreKey = "nft-approved-all-store"

	NftCollectionOwnerKeyPrefix       = "NftCollectionOwner/value/"
	OpenedNftsCollectionKeyPrefix     = "OpenedNftsCollection/value/"
	RestrictedNftsCollectionKeyPrefix = "RestrictedNftsCollection/value/"
	SoulBondedNftsCollectionKeyPrefix = "SoulBondedNftsCollection/value/"

	Delimiter   = []byte{0x00}
	Placeholder = []byte{0x01}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func GetNftCollectionIndex(creator sdk.AccAddress, id string) []byte {
	creator = address.MustLengthPrefix(creator)
	idBz := conv.GetByteKey(id)

	key := make([]byte, len(nftCollectionIndex)+len(Delimiter)+len(creator)+len(Delimiter)+len(idBz)+len(Delimiter))
	copy(key, nftCollectionIndex)
	copy(key[len(nftCollectionIndex):], Delimiter)
	copy(key[len(nftCollectionIndex)+len(Delimiter):], creator)
	copy(key[len(nftCollectionIndex)+len(Delimiter)+len(creator):], Delimiter)
	copy(key[len(nftCollectionIndex)+len(Delimiter)+len(creator)+len(Delimiter):], idBz)
	copy(key[len(nftCollectionIndex)+len(Delimiter)+len(creator)+len(Delimiter)+len(idBz):], Delimiter)
	return key
}

func NftCollectionStoreKey(creator sdk.AccAddress) []byte {
	var key []byte
	if creator == nil {
		key = make([]byte, len(nftCollectionStoreKey)+len(Delimiter))
		copy(key, nftCollectionStoreKey)
		copy(key[len(nftCollectionStoreKey):], Delimiter)
	} else {
		creator = address.MustLengthPrefix(creator)

		key = make([]byte, len(nftCollectionStoreKey)+len(Delimiter)+len(creator)+len(Delimiter))
		copy(key, nftCollectionStoreKey)
		copy(key[len(nftCollectionStoreKey):], Delimiter)
		copy(key[len(nftCollectionStoreKey)+len(Delimiter):], creator)
		copy(key[len(nftCollectionStoreKey)+len(Delimiter)+len(creator):], Delimiter)
	}

	return key
}

func GetNftIndex(collectionIndex []byte, id string) []byte {
	idBz := conv.GetByteKey(id)
	key := make([]byte, len(nftIndex)+len(Delimiter)+len(collectionIndex)+len(Delimiter)+len(idBz)+len(Delimiter))
	copy(key, nftIndex)
	copy(key[len(nftIndex):], Delimiter)
	copy(key[len(nftIndex)+len(Delimiter):], collectionIndex)
	copy(key[len(nftIndex)+len(Delimiter)+len(collectionIndex):], Delimiter)
	copy(key[len(nftIndex)+len(Delimiter)+len(collectionIndex)+len(Delimiter):], id)
	copy(key[len(nftIndex)+len(Delimiter)+len(collectionIndex)+len(Delimiter)+len(id):], Delimiter)
	return key
}

func NftStoreKey(collectionIndex []byte) []byte {
	var key []byte
	if collectionIndex == nil {
		key = make([]byte, len(nftStoreKey)+len(Delimiter))
		copy(key, nftStoreKey)
		copy(key[len(nftStoreKey):], Delimiter)
	} else {
		key = make([]byte, len(nftStoreKey)+len(Delimiter)+len(collectionIndex)+len(Delimiter))
		copy(key, nftStoreKey)
		copy(key[len(nftStoreKey):], Delimiter)
		copy(key[len(nftStoreKey)+len(Delimiter):], collectionIndex)
		copy(key[len(nftStoreKey)+len(Delimiter)+len(collectionIndex):], Delimiter)
	}

	return key
}

func NftApprovedStoreKey(collectionIndex []byte) []byte {
	key := make([]byte, len(nftApprovedStoreKey)+len(Delimiter)+len(collectionIndex)+len(Delimiter))
	copy(key, nftApprovedStoreKey)
	copy(key[len(nftApprovedStoreKey):], Delimiter)
	copy(key[len(nftApprovedStoreKey)+len(Delimiter):], collectionIndex)
	copy(key[len(nftApprovedStoreKey)+len(Delimiter)+len(collectionIndex):], Delimiter)
	return key
}

func NftApprovedAllStoreKey() []byte {
	key := make([]byte, len(nftApprovedAllStoreKey)+len(Delimiter))
	copy(key, nftApprovedAllStoreKey)
	copy(key[len(nftApprovedAllStoreKey):], Delimiter)
	return key
}

func GetNftApprovedAllIndex(owner sdk.AccAddress) []byte {
	owner = address.MustLengthPrefix(owner)

	key := make([]byte, len(nftApprovedAllIndex)+len(Delimiter)+len(owner)+len(Delimiter))
	copy(key, nftApprovedAllIndex)
	copy(key[len(nftApprovedAllIndex):], Delimiter)
	copy(key[len(nftApprovedAllIndex)+len(Delimiter):], owner)
	copy(key[len(nftApprovedAllIndex)+len(Delimiter)+len(owner):], Delimiter)
	return key
}
