package types

import (
	"github.com/LimeChain/mantrachain/internal/conv"
)

var (
	domainIndex        = "domain-id"
	domainNameIndex    = "domain-name-id"
	domainStoreKey     = "domain-store"
	domainNameStoreKey = "domain-name-store"

	delimiter = []byte{0x00}
)

const (
	// ModuleName defines the module name
	ModuleName = "mns"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_mns"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	AttributeKeyDomain     = "domain"
	AttributeKeyDomainName = "domain_name"
	AttributeKeyDid        = "did"
	AttributeKeyDomainType = "domain_type"
	AttributeKeyOwner      = "owner"
	AttributeKeyCreator    = "creator"
)

func GetDomainIndex(domain string) []byte {
	domainBz := conv.UnsafeStrToBytes(domain)

	key := make([]byte, len(domainIndex)+len(delimiter)+len(domainBz)+len(delimiter))
	copy(key, domainIndex)
	copy(key[len(domainIndex):], delimiter)
	copy(key[len(domainIndex)+len(delimiter):], domainBz)
	copy(key[len(domainIndex)+len(delimiter)+len(domainBz):], delimiter)
	return key
}

func DomainStoreKey() []byte {
	key := make([]byte, len(domainStoreKey)+len(delimiter))
	copy(key, domainStoreKey)
	copy(key[len(domainStoreKey):], delimiter)
	return key
}

func GetDomainNameIndex(domain, domainName string) []byte {
	domainBz := conv.UnsafeStrToBytes(domain)
	domainNameBz := conv.UnsafeStrToBytes(domainName)

	key := make([]byte, len(domainNameIndex)+len(delimiter)+len(domainBz)+len(delimiter)+len(domainNameBz)+len(delimiter))
	copy(key, domainNameIndex)
	copy(key[len(domainNameIndex):], delimiter)
	copy(key[len(domainNameIndex)+len(delimiter):], domainBz)
	copy(key[len(domainNameIndex)+len(delimiter)+len(domainBz):], delimiter)
	copy(key[len(domainNameIndex)+len(delimiter)+len(domainBz)+len(delimiter):], domainNameBz)
	copy(key[len(domainNameIndex)+len(delimiter)+len(domainBz)+len(delimiter)+len(domainNameBz):], delimiter)
	return key
}

func DomainNameStoreKey(domain string) []byte {
	domainBz := conv.UnsafeStrToBytes(domain)

	key := make([]byte, len(domainNameStoreKey)+len(delimiter)+len(domainBz)+len(delimiter))
	copy(key, domainNameStoreKey)
	copy(key[len(domainNameStoreKey):], delimiter)
	copy(key[len(domainNameStoreKey)+len(delimiter):], domainBz)
	copy(key[len(domainNameStoreKey)+len(delimiter)+len(domainBz):], delimiter)
	return key
}
