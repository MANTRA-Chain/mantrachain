package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRegisterMarketplace = "register_marketplace"

var _ sdk.Msg = &MsgRegisterMarketplace{}

func NewMsgRegisterMarketplace(creator string, marketplace *MsgMarketplaceMetadata) *MsgRegisterMarketplace {
	return &MsgRegisterMarketplace{
		Creator:     creator,
		Marketplace: marketplace,
	}
}

func (msg *MsgRegisterMarketplace) Route() string {
	return RouterKey
}

func (msg *MsgRegisterMarketplace) Type() string {
	return TypeMsgRegisterMarketplace
}

func (msg *MsgRegisterMarketplace) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRegisterMarketplace) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRegisterMarketplace) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if strings.TrimSpace(msg.Marketplace.Id) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "marketplace id should not be empty")
	}
	return nil
}
