package types

import (
	"errors"

	"cosmossdk.io/math"
)

// ValidateBasic implements Authorization.ValidateBasic.
func (msg MsgUpdateParams) ValidateBasic() error {
	proportion := math.LegacyZeroDec()
	if msg.McaTax != "" {
		var err error
		proportion, err = math.LegacyNewDecFromStr(msg.McaTax)
		if err != nil {
			return err
		}
	}
	if proportion.IsNegative() {
		return errors.New("proportion cannot be negative")
	}
	return nil
}
