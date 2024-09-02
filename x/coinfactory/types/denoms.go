package types

import (
	fmt "fmt"
	"strings"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleDenomPrefix = "factory"
	// See the CoinFactory readme for a derivation of these.
	// TL;DR, MaxSubdenomLength + MaxHrpLength = 60 comes from SDK max denom length = 128
	// and the structure of coinfactory denoms.
	MaxSubdenomLength = 44
	MaxHrpLength      = 16
	// MaxCreatorLength = 59 + MaxHrpLength
	MaxCreatorLength = 59 + MaxHrpLength
)

// GetTokenDenom constructs a denom string for tokens created by coinfactory
// based on an input creator address and a subdenom
// The denom constructed is factory/{creator}/{subdenom}
func GetTokenDenom(creator, subdenom string) (string, error) {
	if len(subdenom) > MaxSubdenomLength {
		return "", ErrSubdenomTooLong
	}
	if len(creator) > MaxCreatorLength {
		return "", ErrCreatorTooLong
	}
	if strings.Contains(creator, "/") {
		return "", ErrInvalidCreator
	}
	denom := strings.Join([]string{ModuleDenomPrefix, creator, subdenom}, "/")
	return denom, sdk.ValidateDenom(denom)
}

// DeconstructDenom takes a token denom string and verifies that it is a valid
// denom of the coinfactory module, and is of the form `factory/{creator}/{subdenom}`
// If valid, it returns the creator address and subdenom
func DeconstructDenom(denom string) (creator string, subdenom string, err error) {
	err = sdk.ValidateDenom(denom)
	if err != nil {
		return "", "", err
	}

	strParts := strings.Split(denom, "/")
	if len(strParts) < 3 {
		return "", "", errors.Wrapf(ErrInvalidDenom, "not enough parts of denom %s", denom)
	}

	if strParts[0] != ModuleDenomPrefix {
		return "", "", errors.Wrapf(ErrInvalidDenom, "denom prefix is incorrect. Is: %s.  Should be: %s", strParts[0], ModuleDenomPrefix)
	}

	creator = strParts[1]
	creatorAddr, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return "", "", errors.Wrapf(ErrInvalidDenom, "Invalid creator address (%s)", err)
	}

	// Handle the case where a denom has a slash in its subdenom. For example,
	// when we did the split, we'd turn factory/accaddr/atomderivative/sikka into ["factory", "accaddr", "atomderivative", "sikka"]
	// So we have to join [2:] with a "/" as the delimiter to get back the correct subdenom which should be "atomderivative/sikka"
	subdenom = strings.Join(strParts[2:], "/")

	return creatorAddr.String(), subdenom, nil
}

// NewCoinFactoryDenomMintCoinsRestriction creates and returns a function that only allows minting of
// valid coinfactory denoms
func NewCoinFactoryDenomMintCoinsRestriction() func(ctx sdk.Context, coinsToMint sdk.Coins) error {
	return func(ctx sdk.Context, coinsToMint sdk.Coins) error {
		for _, coin := range coinsToMint {
			_, _, err := DeconstructDenom(coin.Denom)
			if err != nil {
				return fmt.Errorf("does not have permission to mint %s", coin.Denom)
			}
		}
		return nil
	}
}
