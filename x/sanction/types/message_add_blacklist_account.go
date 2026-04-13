package types

func NewMsgAddBlacklistAccounts(creator string, blacklistAccounts []string) *MsgAddBlacklistAccounts {
	return &MsgAddBlacklistAccounts{
		Authority:         creator,
		BlacklistAccounts: blacklistAccounts,
	}
}

func (m MsgAddBlacklistAccounts) Validate() error {
	return validateBlacklistMessage(m.Authority, m.BlacklistAccounts)
}
