package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgUpdateAccountPrivileges             = "update_account_privileges"
	TypeMsgUpdateAccountPrivilegesBatch        = "update_account_privileges_batch"
	TypeMsgUpdateAccountPrivilegesGroupedBatch = "update_account_privileges_grouped_batch"
)

var _ sdk.Msg = &MsgUpdateAccountPrivileges{}

func NewMsgUpdateAccountPrivileges(
	creator string,
	account string,
	privileges []byte,

) *MsgUpdateAccountPrivileges {
	return &MsgUpdateAccountPrivileges{
		Creator:    creator,
		Account:    account,
		Privileges: privileges,
	}
}

func (msg *MsgUpdateAccountPrivileges) Route() string {
	return RouterKey
}

func (msg *MsgUpdateAccountPrivileges) Type() string {
	return TypeMsgUpdateAccountPrivileges
}

func (msg *MsgUpdateAccountPrivileges) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateAccountPrivileges) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateAccountPrivileges) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Account)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid account address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateAccountPrivilegesBatch{}

func NewMsgUpdateAccountPrivilegesBatch(
	creator string,
	accountsPrivileges MsgAccountsPrivileges,
) *MsgUpdateAccountPrivilegesBatch {
	return &MsgUpdateAccountPrivilegesBatch{
		Creator:            creator,
		AccountsPrivileges: &accountsPrivileges,
	}
}

func (msg *MsgUpdateAccountPrivilegesBatch) Route() string {
	return RouterKey
}

func (msg *MsgUpdateAccountPrivilegesBatch) Type() string {
	return TypeMsgUpdateAccountPrivilegesBatch
}

func (msg *MsgUpdateAccountPrivilegesBatch) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateAccountPrivilegesBatch) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateAccountPrivilegesBatch) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.AccountsPrivileges == nil || len(msg.AccountsPrivileges.Accounts) == 0 {
		return errors.Wrapf(sdkerrors.ErrKeyNotFound, "accounts and/or privileges are empty")
	}
	if msg.AccountsPrivileges.Privileges == nil || len(msg.AccountsPrivileges.Accounts) != len(msg.AccountsPrivileges.Privileges) {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "accounts and privileges length is not equal")
	}
	for _, acc := range msg.AccountsPrivileges.Accounts {
		_, err = sdk.AccAddressFromBech32(acc)
		if err != nil {
			return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid account address (%s)", err)
		}
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateAccountPrivilegesGroupedBatch{}

func NewMsgUpdateAccountPrivilegesGroupedBatch(
	creator string,
	accountsPrivilegesGrouped MsgAccountsPrivilegesGrouped,
) *MsgUpdateAccountPrivilegesGroupedBatch {
	return &MsgUpdateAccountPrivilegesGroupedBatch{
		Creator:                   creator,
		AccountsPrivilegesGrouped: &accountsPrivilegesGrouped,
	}
}

func (msg *MsgUpdateAccountPrivilegesGroupedBatch) Route() string {
	return RouterKey
}

func (msg *MsgUpdateAccountPrivilegesGroupedBatch) Type() string {
	return TypeMsgUpdateAccountPrivilegesGroupedBatch
}

func (msg *MsgUpdateAccountPrivilegesGroupedBatch) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateAccountPrivilegesGroupedBatch) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateAccountPrivilegesGroupedBatch) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.AccountsPrivilegesGrouped == nil || len(msg.AccountsPrivilegesGrouped.Accounts) == 0 {
		return errors.Wrapf(sdkerrors.ErrKeyNotFound, "grouped accounts and/or privileges are empty")
	}
	if msg.AccountsPrivilegesGrouped.Privileges == nil || len(msg.AccountsPrivilegesGrouped.Accounts) != len(msg.AccountsPrivilegesGrouped.Privileges) {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "accounts and privileges length is not equal")
	}
	for i := range msg.AccountsPrivilegesGrouped.Accounts {
		for _, acc := range msg.AccountsPrivilegesGrouped.Accounts[i].Accounts {
			_, err = sdk.AccAddressFromBech32(acc)
			if err != nil {
				return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid account address (%s)", err)
			}
		}
	}
	return nil
}
