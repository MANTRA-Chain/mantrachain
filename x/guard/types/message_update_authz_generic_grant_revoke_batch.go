package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateAuthzGenericGrantRevokeBatch = "update_authz_generic_grant_revoke_batch"

var _ sdk.Msg = &MsgUpdateAuthzGenericGrantRevokeBatch{}

func NewMsgUpdateAuthzGenericGrantRevokeBatch(creator string, grantee string, authzGrantRevokeMsgsTypes *AuthzGrantRevokeMsgsTypes) *MsgUpdateAuthzGenericGrantRevokeBatch {
	return &MsgUpdateAuthzGenericGrantRevokeBatch{
		Creator:                   creator,
		Grantee:                   grantee,
		AuthzGrantRevokeMsgsTypes: authzGrantRevokeMsgsTypes,
	}
}

func (msg *MsgUpdateAuthzGenericGrantRevokeBatch) Route() string {
	return RouterKey
}

func (msg *MsgUpdateAuthzGenericGrantRevokeBatch) Type() string {
	return TypeMsgUpdateAuthzGenericGrantRevokeBatch
}

func (msg *MsgUpdateAuthzGenericGrantRevokeBatch) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateAuthzGenericGrantRevokeBatch) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateAuthzGenericGrantRevokeBatch) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.AuthzGrantRevokeMsgsTypes == nil || len(msg.AuthzGrantRevokeMsgsTypes.Msgs) == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "authz grant revoke msgs types are empty")
	}
	for _, msg := range msg.AuthzGrantRevokeMsgsTypes.Msgs {
		if strings.TrimSpace(msg.TypeUrl) == "" {
			return sdkerrors.Wrap(sdkerrors.ErrNotFound, "empty type url")
		}
	}
	return nil
}
