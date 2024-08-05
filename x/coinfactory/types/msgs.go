package types

import (
	"strings"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// constants
const (
	TypeMsgCreateDenom      = "create_denom"
	TypeMsgMint             = "cf_mint"
	TypeMsgBurn             = "cf_burn"
	TypeMsgForceTransfer    = "force_transfer"
	TypeMsgChangeAdmin      = "change_admin"
	TypeMsgSetDenomMetadata = "set_denom_metadata"
)

var (
	_ legacytx.LegacyMsg = &MsgCreateDenom{}
	_ sdk.Msg            = &MsgCreateDenom{}
)

// NewMsgCreateDenom creates a msg to create a new denom
func NewMsgCreateDenom(sender, subdenom string) *MsgCreateDenom {
	return &MsgCreateDenom{
		Sender:   sender,
		Subdenom: subdenom,
	}
}

func (msg MsgCreateDenom) Route() string { return RouterKey }
func (msg MsgCreateDenom) Type() string  { return TypeMsgCreateDenom }
func (msg MsgCreateDenom) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	_, err = GetTokenDenom(msg.Sender, msg.Subdenom)
	if err != nil {
		return errors.Wrap(ErrInvalidDenom, err.Error())
	}

	return nil
}

func (msg *MsgCreateDenom) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

var (
	_ legacytx.LegacyMsg = &MsgMint{}
	_ sdk.Msg            = &MsgMint{}
)

// NewMsgMint creates a message to mint tokens
func NewMsgMint(sender string, amount sdk.Coin) *MsgMint {
	return &MsgMint{
		Sender: sender,
		Amount: amount,
	}
}

func NewMsgMintTo(sender string, amount sdk.Coin, mintToAddress string) *MsgMint {
	return &MsgMint{
		Sender:        sender,
		Amount:        amount,
		MintToAddress: mintToAddress,
	}
}

func (msg MsgMint) Route() string { return RouterKey }
func (msg MsgMint) Type() string  { return TypeMsgMint }
func (msg MsgMint) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	if !msg.Amount.IsValid() || msg.Amount.Amount.Equal(math.ZeroInt()) {
		return errors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	if strings.TrimSpace(msg.MintToAddress) != "" {
		_, err = sdk.AccAddressFromBech32(msg.MintToAddress)
		if err != nil {
			return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid mint to address (%s)", err)
		}
	}

	return nil
}

func (msg *MsgMint) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

var (
	_ legacytx.LegacyMsg = &MsgBurn{}
	_ sdk.Msg            = &MsgBurn{}
)

// NewMsgBurn creates a message to burn tokens
func NewMsgBurn(sender string, amount sdk.Coin) *MsgBurn {
	return &MsgBurn{
		Sender: sender,
		Amount: amount,
	}
}

func NewMsgBurnFrom(sender string, amount sdk.Coin, burnFromAddress string) *MsgBurn {
	return &MsgBurn{
		Sender:          sender,
		Amount:          amount,
		BurnFromAddress: burnFromAddress,
	}
}

func (msg MsgBurn) Route() string { return RouterKey }
func (msg MsgBurn) Type() string  { return TypeMsgBurn }
func (msg MsgBurn) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	if !msg.Amount.IsValid() || msg.Amount.Amount.Equal(math.ZeroInt()) {
		return errors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	if strings.TrimSpace(msg.BurnFromAddress) != "" {
		_, err = sdk.AccAddressFromBech32(msg.BurnFromAddress)
		if err != nil {
			return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid burn from address (%s)", err)
		}
	}

	return nil
}

func (msg *MsgBurn) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

var (
	_ legacytx.LegacyMsg = &MsgForceTransfer{}
	_ sdk.Msg            = &MsgForceTransfer{}
)

// NewMsgForceTransfer creates a transfer funds from one account to another
func NewMsgForceTransfer(sender string, amount sdk.Coin, fromAddr, toAddr string) *MsgForceTransfer {
	return &MsgForceTransfer{
		Sender:              sender,
		Amount:              amount,
		TransferFromAddress: fromAddr,
		TransferToAddress:   toAddr,
	}
}

func (msg MsgForceTransfer) Route() string { return RouterKey }
func (msg MsgForceTransfer) Type() string  { return TypeMsgForceTransfer }
func (msg MsgForceTransfer) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.TransferFromAddress)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.TransferToAddress)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid address (%s)", err)
	}

	if !msg.Amount.IsValid() {
		return errors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	return nil
}

func (msg *MsgForceTransfer) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

var (
	_ legacytx.LegacyMsg = &MsgChangeAdmin{}
	_ sdk.Msg            = &MsgChangeAdmin{}
)

// NewMsgChangeAdmin creates a message to burn tokens
func NewMsgChangeAdmin(sender, denom, newAdmin string) *MsgChangeAdmin {
	return &MsgChangeAdmin{
		Sender:   sender,
		Denom:    denom,
		NewAdmin: newAdmin,
	}
}

func (msg MsgChangeAdmin) Route() string { return RouterKey }
func (msg MsgChangeAdmin) Type() string  { return TypeMsgChangeAdmin }
func (msg MsgChangeAdmin) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.NewAdmin)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid address (%s)", err)
	}

	_, _, err = DeconstructDenom(msg.Denom)
	if err != nil {
		return err
	}

	return nil
}

func (msg *MsgChangeAdmin) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

var (
	_ legacytx.LegacyMsg = &MsgSetDenomMetadata{}
	_ sdk.Msg            = &MsgSetDenomMetadata{}
)

// NewMsgChangeAdmin creates a message to burn tokens
func NewMsgSetDenomMetadata(sender string, metadata banktypes.Metadata) *MsgSetDenomMetadata {
	return &MsgSetDenomMetadata{
		Sender:   sender,
		Metadata: metadata,
	}
}

func (msg MsgSetDenomMetadata) Route() string { return RouterKey }
func (msg MsgSetDenomMetadata) Type() string  { return TypeMsgSetDenomMetadata }
func (msg MsgSetDenomMetadata) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	err = msg.Metadata.Validate()
	if err != nil {
		return err
	}

	_, _, err = DeconstructDenom(msg.Metadata.Base)
	if err != nil {
		return err
	}

	return nil
}

func (msg *MsgSetDenomMetadata) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
