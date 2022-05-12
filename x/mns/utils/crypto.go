package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	types "github.com/LimeChain/mantrachain/x/mns/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

func GetDomainIndex(domain string) string {
	index := sha256.Sum256([]byte(fmt.Sprintf("%s/domain/%s", types.ModuleName, domain)))
	return strings.ToUpper(fmt.Sprint("F", hex.EncodeToString(index[:])))
}

func GetDomainNameIndex(domain string, domainName string) string {
	index := sha256.Sum256([]byte(fmt.Sprintf("%s/domain/%s/domain-name/%s", types.ModuleName, domain, domainName)))
	return strings.ToUpper(fmt.Sprint("F", hex.EncodeToString(index[:])))
}

func GetPubKeyHex(pubKey cryptotypes.PubKey) string {
	return strings.ToUpper(fmt.Sprint("F", hex.EncodeToString(pubKey.Bytes())))
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
