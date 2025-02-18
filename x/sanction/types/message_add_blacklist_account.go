package types

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewMsgAddBlacklistAccount(creator string, blacklistAccount string) *MsgAddBlacklistAccount {
	return &MsgAddBlacklistAccount{
		Authority:        creator,
		BlacklistAccount: blacklistAccount,
	}
}

func (m MsgAddBlacklistAccount) Validate() error {
	if m.Authority == "" {
		return errors.New("authority cannot be empty")
	}
	if m.BlacklistAccount == "" {
		return errors.New("blacklistAccount cannot be empty")
	}
	_, err := sdk.AccAddressFromBech32(m.BlacklistAccount)
	if err != nil {
		return fmt.Errorf("invalid account %s", m.BlacklistAccount)
	}
	return nil
}
