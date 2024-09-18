package types

import (
	"fmt"
	"strings"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

// DefaultProportion represents the Proportion default value.
var (
	DefaultProportion = "0.5"
	DefaultMcaAddress = "mantra15m77x4pe6w9vtpuqm22qxu0ds7vn4ehzwx8pls"
)

// NewParams creates a new Params instance.
func NewParams(
	mcaTaxStr string,
	mcaAddress string,
) Params {
	mcaTax := math.LegacyMustNewDecFromStr(mcaTaxStr)
	return Params{
		McaTax:     mcaTax,
		McaAddress: mcaAddress,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		DefaultProportion,
		DefaultMcaAddress,
	)
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if p.McaTax.IsNegative() {
		return fmt.Errorf("mca tax cannot be negative: %s", p.McaTax)
	}
	_, _, err := bech32.DecodeAndConvert(p.McaAddress)
	if err != nil {
		return fmt.Errorf("invalid mca address: %s", p.McaAddress)
	}

	return nil
}

func ValidateMcaTax(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	mcaTax, err := math.LegacyNewDecFromStr(v)
	if err != nil {
		return fmt.Errorf("invalid mca tax: %s", err)
	}

	if mcaTax.IsNegative() || mcaTax.GT(math.LegacyOneDec()) {
		return fmt.Errorf("mca tax must be between 0 and 1")
	}

	return nil
}

func ValidateMcaAddress(address string) error {
	if address == "" {
		return fmt.Errorf("mca address cannot be empty")
	}
	_, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return fmt.Errorf("invalid mca address: %w", err)
	}
	if !strings.HasPrefix(address, "mantra") {
		return fmt.Errorf("mca address must have 'mantra' prefix")
	}
	return nil
}
