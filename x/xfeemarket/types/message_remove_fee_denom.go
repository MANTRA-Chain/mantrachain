package types

import "errors"

func NewMsgRemoveFeeDenom(creator string, denom string) *MsgRemoveFeeDenom {
	return &MsgRemoveFeeDenom{
		Authority: creator,
		Denom:     denom,
	}
}

func (m MsgRemoveFeeDenom) Validate() error {
	if m.Authority == "" {
		return errors.New("authority cannot be empty")
	}
	if m.Denom == "" {
		return errors.New("denom cannot be empty")
	}
	return nil
}
