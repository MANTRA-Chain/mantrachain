package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
)

const (
	TypeMsgUpdateGuardTransferCoins = "update_guard_transfer_coins"

	AttributeKeyCreator            = "creator"
	AttributeKeyAccount            = "account"
	AttributeKeyAccounts           = "accounts"
	AttributeKeyIndex              = "index"
	AttributeKeyIndexes            = "indexes"
	AttributeKeyKind               = "kind"
	AttributeKeyGuardTransferCoins = "guard_transfer_coins"
)

var _ legacytx.LegacyMsg = &MsgUpdateGuardTransferCoins{}
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

func (msg *MsgUpdateGuardTransferCoins) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateGuardTransferCoins) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
