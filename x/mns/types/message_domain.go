package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgCreateDomain = "create_domain"

var _ sdk.Msg = &MsgCreateDomain{}

func NewMsgCreateDomain(
	creator string,
	domain string,
	domainType string,
	pubKeyHex string,
	pubKeyType string) *MsgCreateDomain {

	return &MsgCreateDomain{
		Creator:    creator,
		Domain:     domain,
		DomainType: domainType,
		PubKeyHex:  pubKeyHex,
		PubKeyType: pubKeyType,
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
	if strings.TrimSpace(msg.Domain) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "domain should not be empty")
	}
	if err := ValidateDomainType(DomainType(msg.DomainType)); err != nil {
		return err
	}
	return nil
}
