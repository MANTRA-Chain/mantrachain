package types

import (
	"strings"

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
	bridgeCreator string,
	bridgeId string,

) *MsgCreateChainValidatorBridge {
	return &MsgCreateChainValidatorBridge{
		Creator:       creator,
		Chain:         chain,
		Validator:     validator,
		BridgeCreator: bridgeCreator,
		BridgeId:      bridgeId,
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
	_, err = sdk.AccAddressFromBech32(msg.BridgeCreator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid bridge creator address (%s)", err)
	}
	if strings.TrimSpace(msg.Chain) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "chain should not be empty")
	}
	if strings.TrimSpace(msg.Validator) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "validator should not be empty")
	}
	if strings.TrimSpace(msg.BridgeId) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "bridge id should not be empty")
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateChainValidatorBridge{}

func NewMsgUpdateChainValidatorBridge(
	creator string,
	chain string,
	validator string,
	bridgeCreator string,
	bridgeId string,

) *MsgUpdateChainValidatorBridge {
	return &MsgUpdateChainValidatorBridge{
		Creator:       creator,
		Chain:         chain,
		Validator:     validator,
		BridgeCreator: bridgeCreator,
		BridgeId:      bridgeId,
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
	_, err = sdk.AccAddressFromBech32(msg.BridgeCreator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid bridge creator address (%s)", err)
	}
	if strings.TrimSpace(msg.Chain) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "chain should not be empty")
	}
	if strings.TrimSpace(msg.Validator) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "validator should not be empty")
	}
	if strings.TrimSpace(msg.BridgeId) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "bridge id should not be empty")
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
	if strings.TrimSpace(msg.Chain) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "chain should not be empty")
	}
	if strings.TrimSpace(msg.Validator) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "validator should not be empty")
	}
	return nil
}
