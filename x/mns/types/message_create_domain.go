package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgCreateDomain = "create_domain"

var _ sdk.Msg = &MsgCreateDomain{}

func NewMsgCreateDomain(creator string, domain string) *MsgCreateDomain {
  return &MsgCreateDomain{
		Creator: creator,
    Domain: domain,
	}
}

func (msg *MsgCreateDomain) Route() string {
  return RouterKey
}

func (msg *MsgCreateDomain) Type() string {
  return TypeMsgCreateDomain
}

func (msg *MsgCreateDomain) GetSigners() []sdk.AccAddress {
  creator, err := sdk.AccAddressFromBech32(msg.Creator)
  if err != nil {
    panic(err)
  }
  return []sdk.AccAddress{creator}
}

func (msg *MsgCreateDomain) GetSignBytes() []byte {
  bz := ModuleCdc.MustMarshalJSON(msg)
  return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateDomain) ValidateBasic() error {
  _, err := sdk.AccAddressFromBech32(msg.Creator)
  	if err != nil {
  		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
  	}
  return nil
}

