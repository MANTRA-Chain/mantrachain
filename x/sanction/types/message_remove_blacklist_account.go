package types

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewMsgRemoveBlacklistAccount(creator string, blacklistAccount string) *MsgRemoveBlacklistAccount {
	return &MsgRemoveBlacklistAccount{
		Authority:        creator,
		BlacklistAccount: blacklistAccount,
	}
}

func (m MsgRemoveBlacklistAccount) Validate() error {
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
