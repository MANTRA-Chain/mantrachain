package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateNftCollection = "create_nft_collection"
)

var (
	_ sdk.Msg = &MsgCreateNftCollection{}
)

func NewMsgCreateNftCollection(creator string, collection *MsgCreateNftCollectionMetadata) *MsgCreateNftCollection {
	return &MsgCreateNftCollection{
		Creator:    creator,
		Collection: collection,
	}
}

func (msg *MsgCreateNftCollection) Route() string {
	return RouterKey
}

func (msg *MsgCreateNftCollection) Type() string {
	return TypeMsgCreateNftCollection
}

func (msg *MsgCreateNftCollection) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateNftCollection) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateNftCollection) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.Collection == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "nft collection is empty")
	}
	if strings.TrimSpace(msg.Collection.Id) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "collection id should not be empty")
	}
	return nil
}
