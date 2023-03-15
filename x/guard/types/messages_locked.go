package types

import (
	"strings"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgUpdateLocked = "update_locked"
)

var _ sdk.Msg = &MsgUpdateLocked{}

func NewMsgUpdateLocked(
	creator string,
	index []byte,
	locked bool,
	kind string,
) *MsgUpdateLocked {
	return &MsgUpdateLocked{
		Creator: creator,
		Index:   index,
		Locked:  locked,
		Kind:    kind,
	}
}

func (msg *MsgUpdateLocked) Route() string {
	return RouterKey
}

func (msg *MsgUpdateLocked) Type() string {
	return TypeMsgUpdateLocked
}

func (msg *MsgUpdateLocked) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateLocked) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateLocked) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if len(msg.Index) == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "index should not be empty")
	}
	if strings.TrimSpace(msg.Kind) == "" {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "kind should not be empty")
	}
	return nil
}
