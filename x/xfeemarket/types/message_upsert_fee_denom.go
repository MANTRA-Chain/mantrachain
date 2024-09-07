package types

func NewMsgUpsertFeeDenom(creator string, denom string, multiplier string) *MsgUpsertFeeDenom {
  return &MsgUpsertFeeDenom{
		Creator: creator,
    Denom: denom,
    Multiplier: multiplier,
	}
}