package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgUpdateGuardTransfer = "update_guard_transfer"
)

var _ sdk.Msg = &MsgUpdateGuardTransfer{}

func NewMsgUpdateGuardTransfer(creator string, enabled bool) *MsgUpdateGuardTransfer {
	return &MsgUpdateGuardTransfer{
		Creator: creator,
		Enabled: enabled,
	}
}

func (msg *MsgUpdateGuardTransfer) Route() string {
	return RouterKey
}

func (msg *MsgUpdateGuardTransfer) Type() string {
	return TypeMsgUpdateGuardTransfer
}

func (msg *MsgUpdateGuardTransfer) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateGuardTransfer) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateGuardTransfer) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
