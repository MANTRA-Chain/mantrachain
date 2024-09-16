package types

import (
	"fmt"

	"cosmossdk.io/math"
)

// DefaultProportion represents the Proportion default value.
// TODO: Determine the default value.
var DefaultProportion = math.LegacyMustNewDecFromStr("0.5")

// NewParams creates a new Params instance.
func NewParams(
	proportion math.LegacyDec,
) Params {
	return Params{
		Proportion: proportion,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		DefaultProportion,
	)
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if p.Proportion.IsNegative() {
		return fmt.Errorf("proportion cannot be negative: %s", p.Proportion)
	}

	return nil
}
