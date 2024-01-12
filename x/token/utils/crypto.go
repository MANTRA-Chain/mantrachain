package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	types "github.com/AumegaChain/aumega/x/token/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

func GetIndexHex(index []byte) string {
	b := sha256.Sum256([]byte(index))
	return hex.EncodeToString(b[:])
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
