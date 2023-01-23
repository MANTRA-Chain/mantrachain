package types

import (
	"github.com/LimeChain/mantrachain/internal/conv"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	AttributeKeyMarketplaceId      = "marketplace_id"
	AttributeKeyMarketplaceCreator = "marketplace_creator"
	AttributeKeyCollectionId       = "collection_id"
	AttributeKeyCollectionCreator  = "collection_creator"
	AttributeKeyNftId              = "nft_id"
	AttributeKeySigner             = "signer"
	AttributeKeyOwner              = "owner"
	AttributeKeyReceiver           = "receiver"
	AttributeKeyDelegated          = "delegated"
	AttributeKeyStakingChain       = "staking_chain"
	AttributeKeyStakingValidator   = "staking_validator"
)

var (
	marketplaceIndex              = "marketplace-id"
	marketplaceStoreKey           = "marketplace-store"
	marketplaceCollectionStoreKey = "marketplace-collection-store"
	marketplaceNftStoreKey        = "marketplace-nft-store"

	delimiter = []byte{0x00}
)

const (
	// ModuleName defines the module name
	ModuleName = "marketplace"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_marketplace"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func GetMarketplaceIndex(creator sdk.AccAddress, id string) []byte {
	creator = address.MustLengthPrefix(creator)
	idBz := conv.UnsafeStrToBytes(id)

	key := make([]byte, len(marketplaceIndex)+len(delimiter)+len(creator)+len(delimiter)+len(idBz)+len(delimiter))
	copy(key, marketplaceIndex)
	copy(key[len(marketplaceIndex):], delimiter)
	copy(key[len(marketplaceIndex)+len(delimiter):], creator)
	copy(key[len(marketplaceIndex)+len(delimiter)+len(creator):], delimiter)
	copy(key[len(marketplaceIndex)+len(delimiter)+len(creator)+len(delimiter):], idBz)
	copy(key[len(marketplaceIndex)+len(delimiter)+len(creator)+len(delimiter)+len(idBz):], delimiter)
	return key
}

func MarketplaceStoreKey(creator sdk.AccAddress) []byte {
	var key []byte
	if creator == nil {
		key = make([]byte, len(marketplaceStoreKey)+len(delimiter))
		copy(key, marketplaceStoreKey)
		copy(key[len(marketplaceStoreKey):], delimiter)
	} else {
		creator = address.MustLengthPrefix(creator)

		key = make([]byte, len(marketplaceStoreKey)+len(delimiter)+len(creator)+len(delimiter))
		copy(key, marketplaceStoreKey)
		copy(key[len(marketplaceStoreKey):], delimiter)
		copy(key[len(marketplaceStoreKey)+len(delimiter):], creator)
		copy(key[len(marketplaceStoreKey)+len(delimiter)+len(creator):], delimiter)
	}

	return key
}

func MarketplaceCollectionStoreKey(marketplaceIndex []byte) []byte {
	var key []byte
	if marketplaceIndex == nil {
		key = make([]byte, len(marketplaceCollectionStoreKey)+len(delimiter))
		copy(key, marketplaceCollectionStoreKey)
		copy(key[len(marketplaceCollectionStoreKey):], delimiter)
	} else {
		key = make([]byte, len(marketplaceCollectionStoreKey)+len(delimiter)+len(marketplaceIndex)+len(delimiter))
		copy(key, marketplaceCollectionStoreKey)
		copy(key[len(marketplaceCollectionStoreKey):], delimiter)
		copy(key[len(marketplaceCollectionStoreKey)+len(delimiter):], marketplaceIndex)
		copy(key[len(marketplaceCollectionStoreKey)+len(delimiter)+len(marketplaceIndex):], delimiter)
	}

	return key
}

func MarketplaceNftStoreKey(marketplaceIndex []byte, collectionIndex []byte) []byte {
	var key []byte
	if marketplaceIndex == nil && collectionIndex == nil {
		key = make([]byte, len(marketplaceNftStoreKey)+len(delimiter))
		copy(key, marketplaceNftStoreKey)
		copy(key[len(marketplaceNftStoreKey):], delimiter)
	} else {
		key = make([]byte, len(marketplaceNftStoreKey)+len(delimiter)+len(marketplaceIndex)+len(delimiter)+len(collectionIndex)+len(delimiter))
		copy(key, marketplaceNftStoreKey)
		copy(key[len(marketplaceNftStoreKey):], delimiter)
		copy(key[len(marketplaceNftStoreKey)+len(delimiter):], marketplaceIndex)
		copy(key[len(marketplaceNftStoreKey)+len(delimiter)+len(marketplaceIndex):], delimiter)
		copy(key[len(marketplaceNftStoreKey)+len(delimiter)+len(marketplaceIndex)+len(delimiter):], collectionIndex)
		copy(key[len(marketplaceNftStoreKey)+len(delimiter)+len(marketplaceIndex)+len(delimiter)+len(collectionIndex):], delimiter)
	}

	return key
}
