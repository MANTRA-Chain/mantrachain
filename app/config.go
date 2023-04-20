package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DisplayDenom = "aum"
	BaseDenom    = "uaum"
	AumExponent  = 6
)

func InitSDKConfig() {
	SetAddressPrefixes()
	RegisterDenoms()
}

// RegisterDenoms registers token denoms.
func RegisterDenoms() {
	err := sdk.RegisterDenom(DisplayDenom, sdk.OneDec())
	if err != nil {
		panic(err)
	}
	err = sdk.RegisterDenom(BaseDenom, sdk.NewDecWithPrec(1, AumExponent))
	if err != nil {
		panic(err)
	}
}

// SetAddressPrefixes builds the Config with Bech32 addressPrefix and publKeyPrefix for accounts, validators, and consensus nodes and verifies that addreeses have correct format.
func SetAddressPrefixes() {
	// Set prefixes
	accountPubKeyPrefix := AccountAddressPrefix + "pub"
	validatorAddressPrefix := AccountAddressPrefix + "valoper"
	validatorPubKeyPrefix := AccountAddressPrefix + "valoperpub"
	consNodeAddressPrefix := AccountAddressPrefix + "valcons"
	consNodePubKeyPrefix := AccountAddressPrefix + "valconspub"

	// Set and seal config
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AccountAddressPrefix, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
	config.Seal()
}
