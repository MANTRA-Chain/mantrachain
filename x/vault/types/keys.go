package types

import (
	"github.com/LimeChain/mantrachain/internal/conv"
)

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
	nftStakeStoreKey    = "nft-stake-store"
	epochStoreKey       = "epoch-store"
	epochYieldStoreKey  = "epoch-yield-store"
	lastEpochBlockIndex = "last-epoch-block-id"
	stakedIndex         = "staked-id"
	epochIndex          = "epoch-id"

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

func EpochStoreKey(chain string, validator string) []byte {
	chainBz := conv.UnsafeStrToBytes(chain)
	validatorBz := conv.UnsafeStrToBytes(validator)
	key := make([]byte, len(epochStoreKey)+len(delimiter)+len(chainBz)+len(delimiter)+len(validatorBz)+len(delimiter))
	copy(key, epochStoreKey)
	copy(key[len(epochStoreKey):], delimiter)
	copy(key[len(epochStoreKey)+len(delimiter):], chainBz)
	copy(key[len(epochStoreKey)+len(delimiter)+len(chainBz):], delimiter)
	copy(key[len(epochStoreKey)+len(delimiter)+len(chainBz)+len(delimiter):], validatorBz)
	copy(key[len(epochStoreKey)+len(delimiter)+len(chainBz)+len(delimiter)+len(validatorBz):], delimiter)
	return key
}

func GetLastEpochBlockIndex(denom string) []byte {
	denomBz := conv.UnsafeStrToBytes(denom)
	key := make([]byte, len(lastEpochBlockIndex)+len(delimiter)+len(denomBz)+len(delimiter))
	copy(key, lastEpochBlockIndex)
	copy(key[len(lastEpochBlockIndex):], delimiter)
	copy(key[len(lastEpochBlockIndex)+len(delimiter):], denomBz)
	copy(key[len(lastEpochBlockIndex)+len(delimiter)+len(denomBz):], delimiter)
	return key
}

func GetStakedIndex(denom string) []byte {
	denomBz := conv.UnsafeStrToBytes(denom)
	key := make([]byte, len(stakedIndex)+len(delimiter)+len(denomBz)+len(delimiter))
	copy(key, stakedIndex)
	copy(key[len(stakedIndex):], delimiter)
	copy(key[len(stakedIndex)+len(delimiter):], denomBz)
	copy(key[len(stakedIndex)+len(delimiter)+len(denomBz):], delimiter)
	return key
}

func GetEpochIndex(denom string, index []byte) []byte {
	var key []byte
	denomBz := conv.UnsafeStrToBytes(denom)
	if index == nil {
		key = make([]byte, len(epochIndex)+len(delimiter)+len(denomBz)+len(delimiter))
		copy(key, epochIndex)
		copy(key[len(epochIndex):], delimiter)
		copy(key[len(epochIndex)+len(delimiter):], denomBz)
		copy(key[len(epochIndex)+len(delimiter)+len(denomBz):], delimiter)
	} else {
		denomBz := conv.UnsafeStrToBytes(denom)
		key = make([]byte, len(epochIndex)+len(delimiter)+len(denomBz)+len(delimiter)+len(index)+len(delimiter))
		copy(key, epochIndex)
		copy(key[len(epochIndex):], delimiter)
		copy(key[len(epochIndex)+len(delimiter):], denomBz)
		copy(key[len(epochIndex)+len(delimiter)+len(denomBz):], delimiter)
		copy(key[len(epochIndex)+len(delimiter)+len(denomBz)+len(delimiter):], index)
		copy(key[len(epochIndex)+len(delimiter)+len(denomBz)+len(delimiter)+len(index):], delimiter)
	}

	return key
}

func KeyPrefix(p string) []byte {
	return []byte(p)
}
