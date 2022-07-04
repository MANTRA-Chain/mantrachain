package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Mint Nfts
const TypeMsgMintNft = "mint_nfts"

var _ sdk.Msg = &MsgMintNfts{}

func NewMsgMintNfts(creator string, collectionCreator string, collectionId string,
	nfts *MsgNftsMetadata,
	pubKeyHex string,
	pubKeyType string) *MsgMintNfts {
	return &MsgMintNfts{
		Creator:           creator,
		CollectionCreator: collectionCreator,
		CollectionId:      collectionId,
		Nfts:              nfts,
		PubKeyHex:         pubKeyHex,
		PubKeyType:        pubKeyType,
	}
}

func (msg *MsgMintNfts) Route() string {
	return RouterKey
}

func (msg *MsgMintNfts) Type() string {
	return TypeMsgMintNft
}

func (msg *MsgMintNfts) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgMintNfts) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgMintNfts) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

// Burn Nfts
const TypeMsgBurnNft = "burn_nfts"

var _ sdk.Msg = &MsgBurnNfts{}

func NewMsgBurnNfts(creator string, collectionCreator string, collectionId string,
	nftsIds *MsgNftsIds,
	pubKeyHex string,
	pubKeyType string) *MsgBurnNfts {
	return &MsgBurnNfts{
		Creator:           creator,
		CollectionCreator: collectionCreator,
		CollectionId:      collectionId,
		Nfts:              nftsIds,
		PubKeyHex:         pubKeyHex,
		PubKeyType:        pubKeyType,
	}
}

func (msg *MsgBurnNfts) Route() string {
	return RouterKey
}

func (msg *MsgBurnNfts) Type() string {
	return TypeMsgBurnNft
}

func (msg *MsgBurnNfts) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgBurnNfts) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgBurnNfts) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
