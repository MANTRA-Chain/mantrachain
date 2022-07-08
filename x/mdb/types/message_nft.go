package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgMintNfts       = "mint_nfts"
	TypeMsgBurnNfts       = "burn_nfts"
	TypeMsgTransferNfts   = "transfer_nfts"
	TypeMsgApproveNfts    = "approve_nfts"
	TypeMsgApproveAllNfts = "approve_all_nfts"
	TypeMsgMintNft        = "mint_nft"
	TypeMsgBurnNft        = "burn_nft"
	TypeMsgTransferNft    = "transfer_nft"
	TypeMsgApproveNft     = "approve_nft"
)

var (
	_ sdk.Msg = &MsgMintNfts{}
	_ sdk.Msg = &MsgBurnNfts{}
	_ sdk.Msg = &MsgTransferNfts{}
	_ sdk.Msg = &MsgApproveNfts{}
	_ sdk.Msg = &MsgApproveAllNfts{}
	_ sdk.Msg = &MsgMintNft{}
	_ sdk.Msg = &MsgBurnNft{}
	_ sdk.Msg = &MsgTransferNft{}
	_ sdk.Msg = &MsgApproveNft{}
)

func NewMsgMintNfts(creator string, collectionCreator string, collectionId string,
	nfts *MsgNftsMetadata,
	receiver string,
	pubKeyHex string,
	pubKeyType string,
	strictCollection bool,
) *MsgMintNfts {
	return &MsgMintNfts{
		Creator:           creator,
		CollectionCreator: collectionCreator,
		CollectionId:      collectionId,
		Nfts:              nfts,
		Receiver:          receiver,
		PubKeyHex:         pubKeyHex,
		PubKeyType:        pubKeyType,
		StrictCollection:  strictCollection,
	}
}

func (msg *MsgMintNfts) Route() string {
	return RouterKey
}

func (msg *MsgMintNfts) Type() string {
	return TypeMsgMintNfts
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
	if msg.Receiver != "" {
		_, err = sdk.AccAddressFromBech32(msg.Receiver)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receiver address (%s)", err)
		}
	}
	if msg.StrictCollection {
		if msg.CollectionCreator == "" {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "collection creator should not be empty with strict-collection flag")
		}
		if msg.CollectionId == "" {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "collection id should not be empty with strict-collection flag")
		}
	}
	if msg.CollectionCreator != "" {
		_, err := sdk.AccAddressFromBech32(msg.CollectionCreator)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid collection creator address (%s)", err)
		}
	}
	return nil
}

func NewMsgBurnNfts(creator string, collectionCreator string, collectionId string,
	nftsIds *MsgNftsIds,
	pubKeyHex string,
	pubKeyType string,
	strictCollection bool,
) *MsgBurnNfts {
	return &MsgBurnNfts{
		Creator:           creator,
		CollectionCreator: collectionCreator,
		CollectionId:      collectionId,
		Nfts:              nftsIds,
		PubKeyHex:         pubKeyHex,
		PubKeyType:        pubKeyType,
		StrictCollection:  strictCollection,
	}
}

func (msg *MsgBurnNfts) Route() string {
	return RouterKey
}

func (msg *MsgBurnNfts) Type() string {
	return TypeMsgBurnNfts
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
	if msg.StrictCollection {
		if msg.CollectionCreator == "" {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "collection creator should not be empty with strict-collection flag")
		}
		if msg.CollectionId == "" {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "collection id should not be empty with strict-collection flag")
		}
	}
	if msg.CollectionCreator != "" {
		_, err := sdk.AccAddressFromBech32(msg.CollectionCreator)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid collection creator address (%s)", err)
		}
	}
	return nil
}

func NewMsgTransferNfts(creator string, collectionCreator string, collectionId string,
	nftsIds *MsgNftsIds,
	owner string,
	receiver string,
	pubKeyHex string,
	pubKeyType string,
	strictCollection bool,
) *MsgTransferNfts {
	return &MsgTransferNfts{
		Creator:           creator,
		Owner:             owner,
		Receiver:          receiver,
		CollectionCreator: collectionCreator,
		CollectionId:      collectionId,
		Nfts:              nftsIds,
		PubKeyHex:         pubKeyHex,
		PubKeyType:        pubKeyType,
		StrictCollection:  strictCollection,
	}
}

func (msg *MsgTransferNfts) Route() string {
	return RouterKey
}

func (msg *MsgTransferNfts) Type() string {
	return TypeMsgTransferNfts
}

func (msg *MsgTransferNfts) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgTransferNfts) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgTransferNfts) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receiver owner (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receiver address (%s)", err)
	}
	if msg.StrictCollection {
		if msg.CollectionCreator == "" {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "collection creator should not be empty with strict-collection flag")
		}
		if msg.CollectionId == "" {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "collection id should not be empty with strict-collection flag")
		}
	}
	if msg.CollectionCreator != "" {
		_, err := sdk.AccAddressFromBech32(msg.CollectionCreator)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid collection creator address (%s)", err)
		}
	}
	return nil
}

func NewMsgApproveNfts(creator string, receiver string, collectionCreator string, collectionId string,
	nftsIds *MsgNftsIds,
	approved bool,
	pubKeyHex string,
	pubKeyType string,
	strictCollection bool,
) *MsgApproveNfts {
	return &MsgApproveNfts{
		Creator:           creator,
		Receiver:          receiver,
		CollectionCreator: collectionCreator,
		CollectionId:      collectionId,
		Nfts:              nftsIds,
		Approved:          approved,
		PubKeyHex:         pubKeyHex,
		PubKeyType:        pubKeyType,
		StrictCollection:  strictCollection,
	}
}

func (msg *MsgApproveNfts) Route() string {
	return RouterKey
}

func (msg *MsgApproveNfts) Type() string {
	return TypeMsgApproveNfts
}

func (msg *MsgApproveNfts) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgApproveNfts) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgApproveNfts) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receiver address (%s)", err)
	}
	if msg.StrictCollection {
		if msg.CollectionCreator == "" {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "collection creator should not be empty with strict-collection flag")
		}
		if msg.CollectionId == "" {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "collection id should not be empty with strict-collection flag")
		}
	}
	if msg.CollectionCreator != "" {
		_, err := sdk.AccAddressFromBech32(msg.CollectionCreator)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid collection creator address (%s)", err)
		}
	}
	return nil
}

func NewMsgApproveAllNfts(creator string, receiver string,
	approved bool, pubKeyHex string,
	pubKeyType string) *MsgApproveAllNfts {
	return &MsgApproveAllNfts{
		Creator:    creator,
		Receiver:   receiver,
		Approved:   approved,
		PubKeyHex:  pubKeyHex,
		PubKeyType: pubKeyType,
	}
}

func (msg *MsgApproveAllNfts) Route() string {
	return RouterKey
}

func (msg *MsgApproveAllNfts) Type() string {
	return TypeMsgApproveAllNfts
}

func (msg *MsgApproveAllNfts) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgApproveAllNfts) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgApproveAllNfts) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receiver address (%s)", err)
	}
	return nil
}

func NewMsgMintNft(creator string, collectionCreator string, collectionId string,
	nft *MsgNftMetadata,
	receiver string,
	pubKeyHex string,
	pubKeyType string,
	strictCollection bool,
) *MsgMintNft {
	return &MsgMintNft{
		Creator:           creator,
		CollectionCreator: collectionCreator,
		CollectionId:      collectionId,
		Nft:               nft,
		Receiver:          receiver,
		PubKeyHex:         pubKeyHex,
		PubKeyType:        pubKeyType,
		StrictCollection:  strictCollection,
	}
}

func (msg *MsgMintNft) Route() string {
	return RouterKey
}

func (msg *MsgMintNft) Type() string {
	return TypeMsgMintNft
}

func (msg *MsgMintNft) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgMintNft) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgMintNft) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.Receiver != "" {
		_, err = sdk.AccAddressFromBech32(msg.Receiver)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receiver address (%s)", err)
		}
	}
	if msg.StrictCollection {
		if msg.CollectionCreator == "" {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "collection creator should not be empty with strict-collection flag")
		}
		if msg.CollectionId == "" {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "collection id should not be empty with strict-collection flag")
		}
	}
	if msg.CollectionCreator != "" {
		_, err := sdk.AccAddressFromBech32(msg.CollectionCreator)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid collection creator address (%s)", err)
		}
	}
	return nil
}

func NewMsgBurnNft(creator string, collectionCreator string, collectionId string,
	nftId string,
	pubKeyHex string,
	pubKeyType string,
	strictCollection bool,
) *MsgBurnNft {
	return &MsgBurnNft{
		Creator:           creator,
		CollectionCreator: collectionCreator,
		CollectionId:      collectionId,
		NftId:             nftId,
		PubKeyHex:         pubKeyHex,
		PubKeyType:        pubKeyType,
		StrictCollection:  strictCollection,
	}
}

func (msg *MsgBurnNft) Route() string {
	return RouterKey
}

func (msg *MsgBurnNft) Type() string {
	return TypeMsgBurnNft
}

func (msg *MsgBurnNft) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgBurnNft) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgBurnNft) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.StrictCollection {
		if msg.CollectionCreator == "" {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "collection creator should not be empty with strict-collection flag")
		}
		if msg.CollectionId == "" {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "collection id should not be empty with strict-collection flag")
		}
	}
	if msg.CollectionCreator != "" {
		_, err := sdk.AccAddressFromBech32(msg.CollectionCreator)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid collection creator address (%s)", err)
		}
	}
	return nil
}

func NewMsgTransferNft(creator string, collectionCreator string, collectionId string,
	nftId string,
	owner string,
	receiver string,
	pubKeyHex string,
	pubKeyType string,
	strictCollection bool,
) *MsgTransferNft {
	return &MsgTransferNft{
		Creator:           creator,
		Owner:             owner,
		Receiver:          receiver,
		CollectionCreator: collectionCreator,
		CollectionId:      collectionId,
		NftId:             nftId,
		PubKeyHex:         pubKeyHex,
		PubKeyType:        pubKeyType,
		StrictCollection:  strictCollection,
	}
}

func (msg *MsgTransferNft) Route() string {
	return RouterKey
}

func (msg *MsgTransferNft) Type() string {
	return TypeMsgTransferNft
}

func (msg *MsgTransferNft) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgTransferNft) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgTransferNft) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receiver address (%s)", err)
	}
	if msg.StrictCollection {
		if msg.CollectionCreator == "" {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "collection creator should not be empty with strict-collection flag")
		}
		if msg.CollectionId == "" {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "collection id should not be empty with strict-collection flag")
		}
	}
	if msg.CollectionCreator != "" {
		_, err := sdk.AccAddressFromBech32(msg.CollectionCreator)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid collection creator address (%s)", err)
		}
	}
	return nil
}

func NewMsgApproveNft(creator string, receiver string, collectionCreator string, collectionId string,
	nftId string,
	approved bool,
	pubKeyHex string,
	pubKeyType string,
	strictCollection bool,
) *MsgApproveNft {
	return &MsgApproveNft{
		Creator:           creator,
		Receiver:          receiver,
		CollectionCreator: collectionCreator,
		CollectionId:      collectionId,
		NftId:             nftId,
		Approved:          approved,
		PubKeyHex:         pubKeyHex,
		PubKeyType:        pubKeyType,
		StrictCollection:  strictCollection,
	}
}

func (msg *MsgApproveNft) Route() string {
	return RouterKey
}

func (msg *MsgApproveNft) Type() string {
	return TypeMsgApproveNft
}

func (msg *MsgApproveNft) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgApproveNft) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgApproveNft) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receiver address (%s)", err)
	}
	if msg.StrictCollection {
		if msg.CollectionCreator == "" {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "collection creator should not be empty with strict-collection flag")
		}
		if msg.CollectionId == "" {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "collection id should not be empty with strict-collection flag")
		}
	}
	if msg.CollectionCreator != "" {
		_, err := sdk.AccAddressFromBech32(msg.CollectionCreator)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid collection creator address (%s)", err)
		}
	}
	return nil
}
