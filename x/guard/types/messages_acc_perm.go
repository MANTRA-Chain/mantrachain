package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateAccPerm = "create_acc_perm"
	TypeMsgUpdateAccPerm = "update_acc_perm"
	TypeMsgDeleteAccPerm = "delete_acc_perm"
)

var _ sdk.Msg = &MsgCreateAccPerm{}

func NewMsgCreateAccPerm(
    creator string,
    cat string,
    whlCurr []string,
    
) *MsgCreateAccPerm {
  return &MsgCreateAccPerm{
		Creator : creator,
		Cat: cat,
		WhlCurr: whlCurr,
        
	}
}

func (msg *MsgCreateAccPerm) Route() string {
  return RouterKey
}

func (msg *MsgCreateAccPerm) Type() string {
  return TypeMsgCreateAccPerm
}

func (msg *MsgCreateAccPerm) GetSigners() []sdk.AccAddress {
  creator, err := sdk.AccAddressFromBech32(msg.Creator)
  if err != nil {
    panic(err)
  }
  return []sdk.AccAddress{creator}
}

func (msg *MsgCreateAccPerm) GetSignBytes() []byte {
  bz := ModuleCdc.MustMarshalJSON(msg)
  return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateAccPerm) ValidateBasic() error {
  _, err := sdk.AccAddressFromBech32(msg.Creator)
  	if err != nil {
  		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
  	}
  return nil
}

var _ sdk.Msg = &MsgUpdateAccPerm{}

func NewMsgUpdateAccPerm(
    creator string,
    cat string,
    whlCurr []string,
    
) *MsgUpdateAccPerm {
  return &MsgUpdateAccPerm{
		Creator: creator,
        Cat: cat,
        WhlCurr: whlCurr,
        
	}
}

func (msg *MsgUpdateAccPerm) Route() string {
  return RouterKey
}

func (msg *MsgUpdateAccPerm) Type() string {
  return TypeMsgUpdateAccPerm
}

func (msg *MsgUpdateAccPerm) GetSigners() []sdk.AccAddress {
  creator, err := sdk.AccAddressFromBech32(msg.Creator)
  if err != nil {
    panic(err)
  }
  return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateAccPerm) GetSignBytes() []byte {
  bz := ModuleCdc.MustMarshalJSON(msg)
  return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateAccPerm) ValidateBasic() error {
  _, err := sdk.AccAddressFromBech32(msg.Creator)
  if err != nil {
    return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
  }
   return nil
}

var _ sdk.Msg = &MsgDeleteAccPerm{}

func NewMsgDeleteAccPerm(
    creator string,
    cat string,
    
) *MsgDeleteAccPerm {
  return &MsgDeleteAccPerm{
		Creator: creator,
		Cat: cat,
        
	}
}
func (msg *MsgDeleteAccPerm) Route() string {
  return RouterKey
}

func (msg *MsgDeleteAccPerm) Type() string {
  return TypeMsgDeleteAccPerm
}

func (msg *MsgDeleteAccPerm) GetSigners() []sdk.AccAddress {
  creator, err := sdk.AccAddressFromBech32(msg.Creator)
  if err != nil {
    panic(err)
  }
  return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteAccPerm) GetSignBytes() []byte {
  bz := ModuleCdc.MustMarshalJSON(msg)
  return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteAccPerm) ValidateBasic() error {
  _, err := sdk.AccAddressFromBech32(msg.Creator)
  if err != nil {
    return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
  }
  return nil
}