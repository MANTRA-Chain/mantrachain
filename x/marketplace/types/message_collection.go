package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgImportCollection = "import_collection"

var _ sdk.Msg = &MsgImportCollection{}

func NewMsgImportCollection(creator string, marketplaceCreator string, marketplaceId string,
	collectionCreator string, collectionId string, settings *MsgCollectionSettings) *MsgImportCollection {
	return &MsgImportCollection{
		Creator:            creator,
		MarketplaceCreator: marketplaceCreator,
		MarketplaceId:      marketplaceId,
		CollectionCreator:  collectionCreator,
		CollectionId:       collectionId,
		Settings:           settings,
	}
}

func (msg *MsgImportCollection) Route() string {
	return RouterKey
}

func (msg *MsgImportCollection) Type() string {
	return TypeMsgImportCollection
}

func (msg *MsgImportCollection) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgImportCollection) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgImportCollection) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.MarketplaceCreator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid marketplace creator address (%s)", err)
	}
	if strings.TrimSpace(msg.MarketplaceId) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "marketplace id should not be empty")
	}
	_, err = sdk.AccAddressFromBech32(msg.CollectionCreator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid collection creator address (%s)", err)
	}
	if strings.TrimSpace(msg.CollectionId) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "collection id should not be empty")
	}
	return nil
}
