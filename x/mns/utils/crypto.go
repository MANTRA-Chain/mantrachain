package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	types "github.com/LimeChain/mantrachain/x/mns/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

func GetDomainId(domain string) string {
	return fmt.Sprintf("%s/domain/%s",
		types.ModuleName,
		domain,
	)
}

func GetDomainIndex(domain string) string {
	index := sha256.Sum256([]byte(GetDomainId(domain)))
	return hex.EncodeToString(index[:])
}

func GetDomainNameId(domain string, domainName string) string {
	return fmt.Sprintf("%s/domain/%s/domain-name/%s",
		types.ModuleName,
		domain,
		domainName,
	)
}

func GetDomainNameIndex(domain string, domainName string) string {
	index := sha256.Sum256([]byte(GetDomainNameId(domain, domainName)))
	return hex.EncodeToString(index[:])
}

func GetPubKeyHex(pubKey cryptotypes.PubKey) string {
	return fmt.Sprint("F", hex.EncodeToString(pubKey.Bytes()))
}

// derivePubKeyType derive the public key type from a public key
func DerivePubKeyType(pubKey cryptotypes.PubKey) (pubKeyType string, err error) {
	switch pubKey.(type) {
	case *ed25519.PubKey:
		pubKeyType = "ed25519"
	case *secp256k1.PubKey:
		pubKeyType = "secp256k1"
	default:
		err = types.ErrKeyFormatNotSupported
	}
	return
}
