package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateChainValidatorBridge = "create_chain_validator_bridge"
	TypeMsgUpdateChainValidatorBridge = "update_chain_validator_bridge"
	TypeMsgDeleteChainValidatorBridge = "delete_chain_validator_bridge"
)

var _ sdk.Msg = &MsgCreateChainValidatorBridge{}

func NewMsgCreateChainValidatorBridge(
	creator string,
	chain string,
	validator string,
	bridgeId string,

) *MsgCreateChainValidatorBridge {
	return &MsgCreateChainValidatorBridge{
		Creator:   creator,
		Chain:     chain,
		Validator: validator,
		BridgeId:  bridgeId,
	}
}

func (msg *MsgCreateChainValidatorBridge) Route() string {
	return RouterKey
}

func (msg *MsgCreateChainValidatorBridge) Type() string {
	return TypeMsgCreateChainValidatorBridge
}

func (msg *MsgCreateChainValidatorBridge) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateChainValidatorBridge) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateChainValidatorBridge) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateChainValidatorBridge{}

func NewMsgUpdateChainValidatorBridge(
	creator string,
	chain string,
	validator string,
	bridgeId string,

) *MsgUpdateChainValidatorBridge {
	return &MsgUpdateChainValidatorBridge{
		Creator:   creator,
		Chain:     chain,
		Validator: validator,
		BridgeId:  bridgeId,
	}
}

func (msg *MsgUpdateChainValidatorBridge) Route() string {
	return RouterKey
}

func (msg *MsgUpdateChainValidatorBridge) Type() string {
	return TypeMsgUpdateChainValidatorBridge
}

func (msg *MsgUpdateChainValidatorBridge) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateChainValidatorBridge) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateChainValidatorBridge) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeleteChainValidatorBridge{}

func NewMsgDeleteChainValidatorBridge(
	creator string,
	chain string,
	validator string,

) *MsgDeleteChainValidatorBridge {
	return &MsgDeleteChainValidatorBridge{
		Creator:   creator,
		Chain:     chain,
		Validator: validator,
	}
}
func (msg *MsgDeleteChainValidatorBridge) Route() string {
	return RouterKey
}

func (msg *MsgDeleteChainValidatorBridge) Type() string {
	return TypeMsgDeleteChainValidatorBridge
}

func (msg *MsgDeleteChainValidatorBridge) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteChainValidatorBridge) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteChainValidatorBridge) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
