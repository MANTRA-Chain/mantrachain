package types

import (
	"fmt"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// DefaultProportion represents the Proportion default value.
var (
	DefaultProportion    = "0.5"
	DefaultMcaAddress, _ = bech32.ConvertAndEncode("mantra", authtypes.NewModuleAddress(govtypes.ModuleName))
)

// NewParams creates a new Params instance.
func NewParams(
	proportionStr string,
	mcaAddress string,
) Params {
	proportion := math.LegacyMustNewDecFromStr(proportionStr)
	return Params{
		Proportion: proportion,
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
	if p.Proportion.IsNegative() {
		return fmt.Errorf("proportion cannot be negative: %s", p.Proportion)
	}
	_, _, err := bech32.DecodeAndConvert(p.McaAddress)
	if err != nil {
		return fmt.Errorf("invalid mca address: %s", p.McaAddress)
	}

	return nil
}
