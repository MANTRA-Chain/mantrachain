package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgUpdateRequiredPrivileges             = "update_required_privileges"
	TypeMsgUpdateRequiredPrivilegesBatch        = "update_required_privileges_batch"
	TypeMsgUpdateRequiredPrivilegesGroupedBatch = "update_required_privileges_grouped_batch"
)

var _ sdk.Msg = &MsgUpdateRequiredPrivileges{}

func NewMsgUpdateRequiredPrivileges(
	creator string,
	index []byte,
	privileges []byte,
	kind string,

) *MsgUpdateRequiredPrivileges {
	return &MsgUpdateRequiredPrivileges{
		Creator:    creator,
		Index:      index,
		Privileges: privileges,
		Kind:       kind,
	}
}

func (msg *MsgUpdateRequiredPrivileges) Route() string {
	return RouterKey
}

func (msg *MsgUpdateRequiredPrivileges) Type() string {
	return TypeMsgUpdateRequiredPrivileges
}

func (msg *MsgUpdateRequiredPrivileges) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateRequiredPrivileges) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateRequiredPrivileges) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if len(msg.Index) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "index should not be empty")
	}
	if strings.TrimSpace(msg.Kind) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "kind should not be empty")
	}
	_, err = ParseRequiredPrivilegesKind(msg.Kind)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "kind is invalid")
	}
	if msg.Privileges != nil && len(msg.Privileges) > 0 && len(msg.Privileges) != 32 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid privileges length (%d)", len(msg.Privileges))
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateRequiredPrivilegesBatch{}

func NewMsgUpdateRequiredPrivilegesBatch(
	creator string,
	requiredPrivileges MsgRequiredPrivileges,
	kind string,
) *MsgUpdateRequiredPrivilegesBatch {
	return &MsgUpdateRequiredPrivilegesBatch{
		Creator:            creator,
		RequiredPrivileges: &requiredPrivileges,
		Kind:               kind,
	}
}

func (msg *MsgUpdateRequiredPrivilegesBatch) Route() string {
	return RouterKey
}

func (msg *MsgUpdateRequiredPrivilegesBatch) Type() string {
	return TypeMsgUpdateRequiredPrivilegesBatch
}

func (msg *MsgUpdateRequiredPrivilegesBatch) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateRequiredPrivilegesBatch) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateRequiredPrivilegesBatch) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if strings.TrimSpace(msg.Kind) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "kind should not be empty")
	}
	_, err = ParseRequiredPrivilegesKind(msg.Kind)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "kind is invalid")
	}
	if msg.RequiredPrivileges == nil || len(msg.RequiredPrivileges.Indexes) == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "indexes and/or privileges are empty")
	}
	if msg.RequiredPrivileges.Privileges == nil || len(msg.RequiredPrivileges.Indexes) != len(msg.RequiredPrivileges.Privileges) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "indexes and privileges length is not equal")
	}
	for i, index := range msg.RequiredPrivileges.Indexes {
		if len(index) == 0 {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid index (%s)", index)
		}
		if msg.RequiredPrivileges.Privileges[i] != nil && len(msg.RequiredPrivileges.Privileges[i]) > 0 && len(msg.RequiredPrivileges.Privileges[i]) != 32 {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid privileges length (%d)", len(msg.RequiredPrivileges.Privileges[i]))
		}
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateRequiredPrivilegesGroupedBatch{}

func NewMsgUpdateRequiredPrivilegesGroupedBatch(
	creator string,
	requiredPrivilegesGrouped MsgRequiredPrivilegesGrouped,
	kind string,
) *MsgUpdateRequiredPrivilegesGroupedBatch {
	return &MsgUpdateRequiredPrivilegesGroupedBatch{
		Creator:                   creator,
		RequiredPrivilegesGrouped: &requiredPrivilegesGrouped,
		Kind:                      kind,
	}
}

func (msg *MsgUpdateRequiredPrivilegesGroupedBatch) Route() string {
	return RouterKey
}

func (msg *MsgUpdateRequiredPrivilegesGroupedBatch) Type() string {
	return TypeMsgUpdateRequiredPrivilegesGroupedBatch
}

func (msg *MsgUpdateRequiredPrivilegesGroupedBatch) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateRequiredPrivilegesGroupedBatch) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateRequiredPrivilegesGroupedBatch) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if strings.TrimSpace(msg.Kind) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "kind should not be empty")
	}
	_, err = ParseRequiredPrivilegesKind(msg.Kind)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "kind is invalid")
	}
	if msg.RequiredPrivilegesGrouped == nil || len(msg.RequiredPrivilegesGrouped.Indexes) == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "grouped indexes and/or privileges are empty")
	}
	if msg.RequiredPrivilegesGrouped.Privileges == nil || len(msg.RequiredPrivilegesGrouped.Indexes) != len(msg.RequiredPrivilegesGrouped.Privileges) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "indexes and privileges length is not equal")
	}
	for i := range msg.RequiredPrivilegesGrouped.Indexes {
		for k, index := range msg.RequiredPrivilegesGrouped.Indexes[i].Indexes {
			if len(index) == 0 {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid index (%s)", index)
			}
			if msg.RequiredPrivilegesGrouped.Privileges[k] != nil && len(msg.RequiredPrivilegesGrouped.Privileges[k]) > 0 && len(msg.RequiredPrivilegesGrouped.Privileges[k]) != 32 {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid privileges length (%d)", len(msg.RequiredPrivilegesGrouped.Privileges[k]))
			}
		}
	}
	return nil
}
