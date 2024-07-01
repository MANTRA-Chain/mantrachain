package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
)

const (
	TypeMsgCreateFeeToken = "create_fee_token"
	TypeMsgUpdateFeeToken = "update_fee_token"
	TypeMsgDeleteFeeToken = "delete_fee_token"
)

var _ legacytx.LegacyMsg = &MsgCreateFeeToken{}
var _ sdk.Msg = &MsgCreateFeeToken{}

func NewMsgCreateFeeToken(
	creator string,
	denom string,
	pairId uint64,

) *MsgCreateFeeToken {
	return &MsgCreateFeeToken{
		Creator: creator,
		Denom:   denom,
		PairId:  pairId,
	}
}

func (msg *MsgCreateFeeToken) Route() string {
	return RouterKey
}

func (msg *MsgCreateFeeToken) Type() string {
	return TypeMsgCreateFeeToken
}

func (msg *MsgCreateFeeToken) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateFeeToken) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ legacytx.LegacyMsg = &MsgUpdateFeeToken{}
var _ sdk.Msg = &MsgUpdateFeeToken{}

func NewMsgUpdateFeeToken(
	creator string,
	denom string,
	pairId uint64,

) *MsgUpdateFeeToken {
	return &MsgUpdateFeeToken{
		Creator: creator,
		Denom:   denom,
		PairId:  pairId,
	}
}

func (msg *MsgUpdateFeeToken) Route() string {
	return RouterKey
}

func (msg *MsgUpdateFeeToken) Type() string {
	return TypeMsgUpdateFeeToken
}

func (msg *MsgUpdateFeeToken) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateFeeToken) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ legacytx.LegacyMsg = &MsgDeleteFeeToken{}
var _ sdk.Msg = &MsgDeleteFeeToken{}

func NewMsgDeleteFeeToken(
	creator string,
	denom string,

) *MsgDeleteFeeToken {
	return &MsgDeleteFeeToken{
		Creator: creator,
		Denom:   denom,
	}
}
func (msg *MsgDeleteFeeToken) Route() string {
	return RouterKey
}

func (msg *MsgDeleteFeeToken) Type() string {
	return TypeMsgDeleteFeeToken
}

func (msg *MsgDeleteFeeToken) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteFeeToken) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
