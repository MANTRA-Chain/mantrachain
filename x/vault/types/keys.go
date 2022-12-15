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
	nftStakeStoreKey             = "nft-stake-store"
	epochStoreKey                = "epoch-store"
	lastEpochBlockStoreKey       = "last-epoch-block-store"
	lastEpochBlockIndex          = "last-epoch-block-id"
	stakedIndex                  = "staked-id"
	epochIndex                   = "epoch-id"
	chainValidatorBridgeStoreKey = "chain-validator-bridge-store"
	chainValidatorBridgeIndex    = "chain-validator-bridge-id"

	delimiter = []byte{0x00}
)

const (
	TypeMsgEpochEnd = "epoch_end"

	AttributeKeyDenom    = "denom"
	AttributeBlockHeight = "block_height"

	AttributeKeyMarketplaceId       = "marketplace_id"
	AttributeKeyMarketplaceCreator  = "marketplace_creator"
	AttributeKeyCollectionId        = "collection_id"
	AttributeKeyCollectionCreator   = "collection_creator"
	AttributeKeyNftId               = "nft_id"
	AttributeKeyChain               = "chain"
	AttributeKeyValidator           = "validator"
	AttributeKeyNftStakeStakedIndex = "nft_stake_staked_index"
	AttributeNftStakeStakedIndex    = "nft_stake_staked_index"
	AttributeKeyDelegated           = "delegated"
)

const (
	TypeNftStakeStakedCreated = "nft_stake_staked_created"
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

func ChainValidatorBridgeStoreKey(chain string) []byte {
	chainBz := conv.UnsafeStrToBytes(chain)
	key := make([]byte, len(chainValidatorBridgeStoreKey)+len(delimiter)+len(chainBz)+len(delimiter))
	copy(key, chainValidatorBridgeStoreKey)
	copy(key[len(chainValidatorBridgeStoreKey):], delimiter)
	copy(key[len(chainValidatorBridgeStoreKey)+len(delimiter):], chainBz)
	copy(key[len(chainValidatorBridgeStoreKey)+len(delimiter)+len(chainBz):], delimiter)
	return key
}

func GetChainValidatorBridgeIndex(validator string) []byte {
	validatorBz := conv.UnsafeStrToBytes(validator)
	key := make([]byte, len(chainValidatorBridgeIndex)+len(delimiter)+len(validatorBz)+len(delimiter))
	copy(key, chainValidatorBridgeIndex)
	copy(key[len(chainValidatorBridgeIndex):], delimiter)
	copy(key[len(chainValidatorBridgeIndex)+len(delimiter):], validatorBz)
	copy(key[len(chainValidatorBridgeIndex)+len(delimiter)+len(validatorBz):], delimiter)
	return key
}

func EpochStoreKey(chain string) []byte {
	chainBz := conv.UnsafeStrToBytes(chain)
	key := make([]byte, len(epochStoreKey)+len(delimiter)+len(chainBz)+len(delimiter))
	copy(key, epochStoreKey)
	copy(key[len(epochStoreKey):], delimiter)
	copy(key[len(epochStoreKey)+len(delimiter):], chainBz)
	copy(key[len(epochStoreKey)+len(delimiter)+len(chainBz):], delimiter)
	return key
}

func LastEpochBlockStoreKey(chain string) []byte {
	chainBz := conv.UnsafeStrToBytes(chain)
	key := make([]byte, len(lastEpochBlockStoreKey)+len(delimiter)+len(chainBz)+len(delimiter))
	copy(key, lastEpochBlockStoreKey)
	copy(key[len(lastEpochBlockStoreKey):], delimiter)
	copy(key[len(lastEpochBlockStoreKey)+len(delimiter):], chainBz)
	copy(key[len(lastEpochBlockStoreKey)+len(delimiter)+len(chainBz):], delimiter)
	return key
}

func GetLastEpochBlockIndex(validator string) []byte {
	validatorBz := conv.UnsafeStrToBytes(validator)
	key := make([]byte, len(lastEpochBlockIndex)+len(delimiter)+len(validatorBz)+len(delimiter))
	copy(key, lastEpochBlockIndex)
	copy(key[len(lastEpochBlockIndex):], delimiter)
	copy(key[len(lastEpochBlockIndex)+len(delimiter):], validatorBz)
	copy(key[len(lastEpochBlockIndex)+len(delimiter)+len(validatorBz):], delimiter)
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

func GetEpochIndex(validator string, index []byte) []byte {
	var key []byte
	validatorBz := conv.UnsafeStrToBytes(validator)
	if index == nil {
		key = make([]byte, len(epochIndex)+len(delimiter)+len(validatorBz)+len(delimiter))
		copy(key, epochIndex)
		copy(key[len(epochIndex):], delimiter)
		copy(key[len(epochIndex)+len(delimiter):], validatorBz)
		copy(key[len(epochIndex)+len(delimiter)+len(validatorBz):], delimiter)
	} else {
		validatorBz := conv.UnsafeStrToBytes(validator)
		key = make([]byte, len(epochIndex)+len(delimiter)+len(validatorBz)+len(delimiter)+len(index)+len(delimiter))
		copy(key, epochIndex)
		copy(key[len(epochIndex):], delimiter)
		copy(key[len(epochIndex)+len(delimiter):], validatorBz)
		copy(key[len(epochIndex)+len(delimiter)+len(validatorBz):], delimiter)
		copy(key[len(epochIndex)+len(delimiter)+len(validatorBz)+len(delimiter):], index)
		copy(key[len(epochIndex)+len(delimiter)+len(validatorBz)+len(delimiter)+len(index):], delimiter)
	}

	return key
}

func KeyPrefix(p string) []byte {
	return []byte(p)
}
