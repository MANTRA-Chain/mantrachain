package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgUpdateGuardTransferCoins = "update_guard_transfer_coins"

	AttributeKeyCreator            = "creator"
	AttributeKeyAccount            = "account"
	AttributeKeyAccounts           = "accounts"
	AttributeKeyLocked             = "locked"
	AttributeKeyIndex              = "index"
	AttributeKeyIndexes            = "indexes"
	AttributeKeyKind               = "kind"
	AttributeKeyGuardTransferCoins = "guard_transfer_coins"
)

var _ sdk.Msg = &MsgUpdateGuardTransferCoins{}

func NewMsgUpdateGuardTransferCoins(creator string, enabled bool) *MsgUpdateGuardTransferCoins {
	return &MsgUpdateGuardTransferCoins{
		Creator: creator,
		Enabled: enabled,
	}
}

func (msg *MsgUpdateGuardTransferCoins) Route() string {
	return RouterKey
}

func (msg *MsgUpdateGuardTransferCoins) Type() string {
	return TypeMsgUpdateGuardTransferCoins
}

func (msg *MsgUpdateGuardTransferCoins) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateGuardTransferCoins) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateGuardTransferCoins) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
