package types

import (
	"errors"

	"cosmossdk.io/math"
)

func NewMsgUpsertFeeDenom(creator string, denom string, multiplier string) *MsgUpsertFeeDenom {
	return &MsgUpsertFeeDenom{
		Authority:  creator,
		Denom:      denom,
		Multiplier: math.LegacyMustNewDecFromStr(multiplier),
	}
}

func (m MsgUpsertFeeDenom) Validate() error {
	if m.Authority == "" {
		return errors.New("authority cannot be empty")
	}
	if m.Denom == "" {
		return errors.New("denom cannot be empty")
	}
	if m.Multiplier.IsNil() {
		return errors.New("multiplier cannot be nil")
	}
	if m.Multiplier.LTE(math.LegacyZeroDec()) {
		return errors.New("multiplier cannot be less than equal 0")
	}
	return nil
}
