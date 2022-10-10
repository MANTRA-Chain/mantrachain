package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRegisterBridge = "register_bridge"

var _ sdk.Msg = &MsgRegisterBridge{}

func NewMsgRegisterBridge(creator string, bridge *MsgBridgeMetadata) *MsgRegisterBridge {
	return &MsgRegisterBridge{
		Creator: creator,
		Bridge:  bridge,
	}
}

func (msg *MsgRegisterBridge) Route() string {
	return RouterKey
}

func (msg *MsgRegisterBridge) Type() string {
	return TypeMsgRegisterBridge
}

func (msg *MsgRegisterBridge) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRegisterBridge) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRegisterBridge) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
