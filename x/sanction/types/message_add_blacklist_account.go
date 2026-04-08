package types

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewMsgAddBlacklistAccounts(creator string, blacklistAccounts []string) *MsgAddBlacklistAccounts {
	return &MsgAddBlacklistAccounts{
		Authority:         creator,
		BlacklistAccounts: blacklistAccounts,
	}
}

func (m MsgAddBlacklistAccounts) Validate() error {
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
