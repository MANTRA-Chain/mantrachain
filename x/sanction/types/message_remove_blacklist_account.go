package types

func NewMsgRemoveBlacklistAccounts(creator string, blacklistAccounts []string) *MsgRemoveBlacklistAccounts {
	return &MsgRemoveBlacklistAccounts{
		Authority:         creator,
		BlacklistAccounts: blacklistAccounts,
	}
}

func (m MsgRemoveBlacklistAccounts) Validate() error {
	return validateBlacklistMessage(m.Authority, m.BlacklistAccounts, 0)
}
