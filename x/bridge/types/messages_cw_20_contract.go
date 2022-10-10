package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateCw20Contract = "create_cw_20_contract"
	TypeMsgUpdateCw20Contract = "update_cw_20_contract"
	TypeMsgDeleteCw20Contract = "delete_cw_20_contract"
)

var _ sdk.Msg = &MsgCreateCw20Contract{}

func NewMsgCreateCw20Contract(creator string, codeId uint64, ver string, path string) *MsgCreateCw20Contract {
	return &MsgCreateCw20Contract{
		Creator: creator,
		CodeId:  codeId,
		Ver:     ver,
		Path:    path,
	}
}

func (msg *MsgCreateCw20Contract) Route() string {
	return RouterKey
}

func (msg *MsgCreateCw20Contract) Type() string {
	return TypeMsgCreateCw20Contract
}

func (msg *MsgCreateCw20Contract) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateCw20Contract) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateCw20Contract) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateCw20Contract{}

func NewMsgUpdateCw20Contract(creator string, codeId uint64, ver string, path string) *MsgUpdateCw20Contract {
	return &MsgUpdateCw20Contract{
		Creator: creator,
		CodeId:  codeId,
		Ver:     ver,
		Path:    path,
	}
}

func (msg *MsgUpdateCw20Contract) Route() string {
	return RouterKey
}

func (msg *MsgUpdateCw20Contract) Type() string {
	return TypeMsgUpdateCw20Contract
}

func (msg *MsgUpdateCw20Contract) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateCw20Contract) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateCw20Contract) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeleteCw20Contract{}

func NewMsgDeleteCw20Contract(creator string) *MsgDeleteCw20Contract {
	return &MsgDeleteCw20Contract{
		Creator: creator,
	}
}
func (msg *MsgDeleteCw20Contract) Route() string {
	return RouterKey
}

func (msg *MsgDeleteCw20Contract) Type() string {
	return TypeMsgDeleteCw20Contract
}

func (msg *MsgDeleteCw20Contract) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteCw20Contract) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteCw20Contract) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
