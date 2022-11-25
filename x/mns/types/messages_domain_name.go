package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateDomainName = "create_domain_name"
)

var _ sdk.Msg = &MsgCreateDomainName{}

func NewMsgCreateDomainName(
	creator string,
	domain string,
	domainName string,
	pubKeyHex string,
	pubKeyType string,

) *MsgCreateDomainName {
	return &MsgCreateDomainName{
		Creator:    creator,
		Domain:     domain,
		DomainName: domainName,
		PubKeyHex:  pubKeyHex,
		PubKeyType: pubKeyType,
	}
}

func (msg *MsgCreateDomainName) Route() string {
	return RouterKey
}

func (msg *MsgCreateDomainName) Type() string {
	return TypeMsgCreateDomainName
}

func (msg *MsgCreateDomainName) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateDomainName) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateDomainName) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if strings.TrimSpace(msg.Domain) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "domain should not be empty")
	}
	if strings.TrimSpace(msg.DomainName) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "domain name should not be empty")
	}
	return nil
}
