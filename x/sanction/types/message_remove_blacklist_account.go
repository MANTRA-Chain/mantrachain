package types

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewMsgRemoveBlacklistAccounts(creator string, blacklistAccounts []string) *MsgRemoveBlacklistAccounts {
	return &MsgRemoveBlacklistAccounts{
		Authority:         creator,
		BlacklistAccounts: blacklistAccounts,
	}
}

func (m MsgRemoveBlacklistAccounts) Validate() error {
	if m.Authority == "" {
		return errors.New("authority cannot be empty")
	}
	if len(m.BlacklistAccounts) == 0 || m.BlacklistAccounts == nil {
		return errors.New("blacklistAccounts cannot be empty")
	}
	for _, account := range m.BlacklistAccounts {
		_, err := sdk.AccAddressFromBech32(account)
		if err != nil {
			return fmt.Errorf("invalid account %s", account)
		}
	}
	return nil
}
