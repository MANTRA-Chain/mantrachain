package types

import (
	"github.com/LimeChain/mantrachain/internal/conv"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	AttributeKeyBridgeId            = "bridge_id"
	AttributeKeyBridgeAccount       = "bridge_account"
	AttributeKeyBridgeCreator       = "bridge_creator"
	AttributeKeySigner              = "signer"
	AttributeKeyOwner               = "owner"
	AttributeKeyReceivers           = "receivers"
	AttributeKeyStaked              = "staked"
	AttributeKeyCw20ContractAddress = "cw-20-contract-address"
	AttributeKeyCw20ContractCreator = "cw-20-contract-creator"
	AttributeKeyCw20ContractCodeId  = "cw-20-contract-code-id"
	AttributeKeyCw20ContractVersion = "cw-20-contract-version"
)

var (
	bridgeIndex    = "bridge-id"
	bridgeStoreKey = "bridge-store"
	txHashStoreKey = "tx-hash-store"

	delimiter = []byte{0x00}
)

const (
	// ModuleName defines the module name
	ModuleName = "bridge"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_bridge"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func GetTxHashIndex(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)

	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}

func GetBridgeIndex(creator sdk.AccAddress, id string) []byte {
	creator = address.MustLengthPrefix(creator)
	idBz := conv.UnsafeStrToBytes(id)

	key := make([]byte, len(bridgeIndex)+len(delimiter)+len(creator)+len(delimiter)+len(idBz)+len(delimiter))
	copy(key, bridgeIndex)
	copy(key[len(bridgeIndex):], delimiter)
	copy(key[len(bridgeIndex)+len(delimiter):], creator)
	copy(key[len(bridgeIndex)+len(delimiter)+len(creator):], delimiter)
	copy(key[len(bridgeIndex)+len(delimiter)+len(creator)+len(delimiter):], idBz)
	copy(key[len(bridgeIndex)+len(delimiter)+len(creator)+len(delimiter)+len(idBz):], delimiter)
	return key
}

func BridgeStoreKey(creator sdk.AccAddress) []byte {
	var key []byte
	if creator == nil {
		key = make([]byte, len(bridgeStoreKey)+len(delimiter))
		copy(key, bridgeStoreKey)
		copy(key[len(bridgeStoreKey):], delimiter)
	} else {
		creator = address.MustLengthPrefix(creator)

		key = make([]byte, len(bridgeStoreKey)+len(delimiter)+len(creator)+len(delimiter))
		copy(key, bridgeStoreKey)
		copy(key[len(bridgeStoreKey):], delimiter)
		copy(key[len(bridgeStoreKey)+len(delimiter):], creator)
		copy(key[len(bridgeStoreKey)+len(delimiter)+len(creator):], delimiter)
	}

	return key
}

func TxHashStoreKey(bridgeIndex []byte) []byte {
	var key []byte

	if bridgeIndex == nil {
		key = make([]byte, len(txHashStoreKey)+len(delimiter))
		copy(key, txHashStoreKey)
		copy(key[len(txHashStoreKey):], delimiter)
		return key
	} else {
		key = make([]byte, len(txHashStoreKey)+len(delimiter)+len(bridgeIndex)+len(delimiter))
		copy(key, txHashStoreKey)
		copy(key[len(txHashStoreKey):], delimiter)
		copy(key[len(txHashStoreKey)+len(delimiter):], bridgeIndex)
		copy(key[len(txHashStoreKey)+len(delimiter)+len(bridgeIndex):], delimiter)
	}

	return key
}

const (
	Cw20ContractKey = "Cw20Contract-value-"
)
