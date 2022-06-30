package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgCreateNftCollection = "create_nft_collection"

var _ sdk.Msg = &MsgCreateNftCollection{}

func NewMsgCreateNftCollection(creator string, collection *MsgCreateNftCollectionMetadata,
	pubKeyHex string,
	pubKeyType string) *MsgCreateNftCollection {
	return &MsgCreateNftCollection{
		Creator:    creator,
		Collection: collection,
		PubKeyHex:  pubKeyHex,
		PubKeyType: pubKeyType,
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
	return nil
}
