package params

import (
	"cosmossdk.io/math"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	HumanCoinUnit = "om"
	BaseCoinUnit  = "uom"
	OmExponent    = 6

	DefaultBondDenom = BaseCoinUnit
)

var (
	Bech32Prefix = "mantra"
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key.
	Bech32PrefixAccPub = Bech32Prefix + "pub"
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address.
	Bech32PrefixValAddr = Bech32Prefix + "valoper"
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key.
	Bech32PrefixValPub = Bech32Prefix + "valoperpub"
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address.
	Bech32PrefixConsAddr = Bech32Prefix + "valcons"
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key.
	Bech32PrefixConsPub = Bech32Prefix + "valconspub"
)

func init() {
	sdk.SetCoinDenomRegex(MantraCoinDenomRegex)
	SetAddressPrefixes()
	RegisterDenoms()
}

// MantraCoinDenomRegex returns the mantra regex string
// this is used to override the default sdk coin denom regex
func MantraCoinDenomRegex() string {
	return `[a-zA-Z][a-zA-Z0-9/:._-]{1,127}`
}

// RegisterDenoms registers token denoms.
func RegisterDenoms() {
	err := sdk.RegisterDenom(HumanCoinUnit, math.LegacyOneDec())
	if err != nil {
		panic(err)
	}
	err = sdk.RegisterDenom(BaseCoinUnit, math.LegacyNewDecWithPrec(1, OmExponent))
	if err != nil {
		panic(err)
	}
}

// SetAddressPrefixes builds the Config with Bech32 addressPrefix and publKeyPrefix for accounts, validators, and consensus nodes and verifies that addreeses have correct format.
func SetAddressPrefixes() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(Bech32Prefix, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
	config.SetAddressVerifier(wasmtypes.VerifyAddressLen())
	config.Seal()
}
