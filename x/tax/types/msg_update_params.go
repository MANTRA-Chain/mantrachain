package types

import (
	"errors"

	"cosmossdk.io/math"
)

// ValidateBasic implements Authorization.ValidateBasic.
func (msg MsgUpdateParams) ValidateBasic() error {
	proportion := math.LegacyZeroDec()
	if msg.Proportion != "" {
		proportion = math.LegacyMustNewDecFromStr(msg.Proportion)
	}
	if proportion.IsNegative() {
		return errors.New("proportion cannot be negative")
	}
	return nil
}
