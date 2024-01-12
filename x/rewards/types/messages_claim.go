package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgClaim = "claim"
)

var _ sdk.Msg = &MsgClaim{}

func NewMsgClaim(
	creator string,
	pairId uint64,

) *MsgClaim {
	return &MsgClaim{
		Creator: creator,
		PairId:  pairId,
	}
}

func (msg *MsgClaim) Route() string {
	return RouterKey
}

func (msg *MsgClaim) Type() string {
	return TypeMsgClaim
}

func (msg *MsgClaim) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgClaim) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClaim) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
