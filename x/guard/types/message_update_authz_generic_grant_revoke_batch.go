package types

import (
	"strings"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
)

const TypeMsgUpdateAuthzGenericGrantRevokeBatch = "update_authz_generic_grant_revoke_batch"

var (
	_ legacytx.LegacyMsg = &MsgUpdateAuthzGenericGrantRevokeBatch{}
	_ sdk.Msg            = &MsgUpdateAuthzGenericGrantRevokeBatch{}
)

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

func (msg *MsgUpdateAuthzGenericGrantRevokeBatch) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateAuthzGenericGrantRevokeBatch) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.AuthzGrantRevokeMsgsTypes == nil || len(msg.AuthzGrantRevokeMsgsTypes.Msgs) == 0 {
		return errors.Wrapf(errorstypes.ErrKeyNotFound, "authz grant revoke msgs types are empty")
	}
	for _, msg := range msg.AuthzGrantRevokeMsgsTypes.Msgs {
		if strings.TrimSpace(msg.TypeUrl) == "" {
			return errors.Wrap(errorstypes.ErrNotFound, "empty type url")
		}
	}
	return nil
}
