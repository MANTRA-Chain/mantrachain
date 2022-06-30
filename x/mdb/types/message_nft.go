package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgMintNft = "mint_nfts"

var _ sdk.Msg = &MsgMintNfts{}

func NewMsgMintNfts(creator string, collectionCreator string, collectionId string,
	nfts *MsgMintNftsMetadata,
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
