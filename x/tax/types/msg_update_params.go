package types

import (
	"errors"

	"cosmossdk.io/math"
)

// ValidateBasic implements Authorization.ValidateBasic.
func (msg MsgUpdateParams) ValidateBasic() error {
	mcaTax := math.LegacyZeroDec()
	if msg.McaTax != "" {
		var err error
		mcaTax, err = math.LegacyNewDecFromStr(msg.McaTax)
		if err != nil {
			return err
		}
	}
	if mcaTax.IsNegative() {
		return errors.New("mcaTax cannot be negative")
	}

	return nil
}
