package types

import (
	"encoding/binary"

	"github.com/MANTRA-Finance/mantrachain/internal/conv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName defines the module name
	ModuleName = "airdrop"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_airdrop"
)

const (
	AttributeKeyCampaignId      = "campaign_id"
	AttributeKeyCampaignAddress = "campaign_address"
	AttributeKeyCreator         = "creator"
)

const (
	LastCampaignIdKey = "Campaign/lastId/"
)

var (
	campaignIndex = "campaign-id"
	claimedIndex  = "claimed-id"

	campaignStoreKey = "campaign-store"
	claimedStoreKey  = "claimed-store"

	Delimiter   = []byte{0x00}
	Placeholder = []byte{0x01}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func GetCampaignIndex(id string) []byte {
	idBz := conv.UnsafeStrToBytes(id)
	key := make([]byte, len(campaignIndex)+len(Delimiter)+len(idBz)+len(Delimiter))
	copy(key, campaignIndex)
	copy(key[len(campaignIndex):], Delimiter)
	copy(key[len(campaignIndex)+len(Delimiter):], idBz)
	copy(key[len(campaignIndex)+len(Delimiter)+len(idBz):], Delimiter)
	return key
}

func GetClaimedIndex(creator sdk.AccAddress, id uint64) []byte {
	creator = address.MustLengthPrefix(creator)
	idBz := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(idBz, id)
	idBz = idBz[:n]

	key := make([]byte, len(claimedIndex)+len(Delimiter)+len(creator)+len(Delimiter)+len(idBz)+len(Delimiter))
	copy(key, claimedIndex)
	copy(key[len(claimedIndex):], Delimiter)
	copy(key[len(claimedIndex)+len(Delimiter):], creator)
	copy(key[len(claimedIndex)+len(Delimiter)+len(creator):], Delimiter)
	copy(key[len(claimedIndex)+len(Delimiter)+len(creator)+len(Delimiter):], idBz)
	copy(key[len(claimedIndex)+len(Delimiter)+len(creator)+len(Delimiter)+len(idBz):], Delimiter)
	return key
}

func CampaignStoreKey() []byte {
	key := make([]byte, len(campaignStoreKey)+len(Delimiter))
	copy(key, campaignStoreKey)
	copy(key[len(campaignStoreKey):], Delimiter)

	return key
}

func ClaimedStoreKey() []byte {
	key := make([]byte, len(claimedStoreKey)+len(Delimiter))
	copy(key, claimedStoreKey)
	copy(key[len(claimedStoreKey):], Delimiter)

	return key
}
