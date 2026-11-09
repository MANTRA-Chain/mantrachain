package types

import (
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewMsgUpdateParams(authority string, params Params) *MsgUpdateParams {
	return &MsgUpdateParams{Authority: authority, Params: params}
}

func (m MsgUpdateParams) Validate() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	if err := m.Params.Validate(); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	return nil
}

func NewMsgUpdateSanctioners(authority string, sanctioners []Sanctioner) *MsgUpdateSanctioners {
	return &MsgUpdateSanctioners{Authority: authority, Sanctioners: sanctioners}
}

func (m MsgUpdateSanctioners) Validate() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	if err := ValidateSanctioners(m.Sanctioners); err != nil {
		return errorsmod.Wrap(ErrInvalidSanctioner, err.Error())
	}
	return nil
}

func NewMsgSanctionAddresses(signer string, addresses []string, reason string) *MsgSanctionAddresses {
	return &MsgSanctionAddresses{Signer: signer, Addresses: addresses, Reason: reason}
}

func (m MsgSanctionAddresses) Validate() error {
	return validateSanctionerMessage(m.Signer, m.Addresses, m.Reason)
}

func NewMsgUnsanctionAddresses(signer string, addresses []string, reason string) *MsgUnsanctionAddresses {
	return &MsgUnsanctionAddresses{Signer: signer, Addresses: addresses, Reason: reason}
}

func (m MsgUnsanctionAddresses) Validate() error {
	return validateSanctionerMessage(m.Signer, m.Addresses, m.Reason)
}

func NewMsgRevokeSanctioner(signer, address, reason string) *MsgRevokeSanctioner {
	return &MsgRevokeSanctioner{Signer: signer, Address: address, Reason: reason}
}

func (m MsgRevokeSanctioner) Validate() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid signer address (%s)", err)
	}
	if _, err := sdk.AccAddressFromBech32(m.Address); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sanctioner address (%s)", err)
	}
	return validateReason(m.Reason)
}

// validateSanctionerMessage checks the stateless parts of a sanctioner
// message. The caps on address count and reason length are params, so the
// keeper checks those.
func validateSanctionerMessage(signer string, addresses []string, reason string) error {
	if _, err := sdk.AccAddressFromBech32(signer); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid signer address (%s)", err)
	}
	if len(addresses) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "addresses cannot be empty")
	}
	for _, addr := range addresses {
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address: %s", addr)
		}
	}
	return validateReason(reason)
}

// validateReason requires a reason on every sanctioner action, so each
// change on the fast path can be audited.
func validateReason(reason string) error {
	if strings.TrimSpace(reason) == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "reason cannot be empty")
	}
	return nil
}
