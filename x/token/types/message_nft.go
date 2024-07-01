package types

import (
	"strings"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
)

const (
	TypeMsgMintNfts                                       = "mint_nfts"
	TypeMsgBurnNfts                                       = "burn_nfts"
	TypeMsgTransferNfts                                   = "transfer_nfts"
	TypeMsgApproveNfts                                    = "approve_nfts"
	TypeMsgApproveAllNfts                                 = "approve_all_nfts"
	TypeMsgMintNft                                        = "mint_nft"
	TypeMsgBurnNft                                        = "burn_nft"
	TypeMsgTransferNft                                    = "transfer_nft"
	TypeMsgApproveNft                                     = "approve_nft"
	TypeMsgUpdateGuardSoulBondNftImage                    = "update_guard_soul_bond_nft_image"
	TypeMsgUpdateRestrictedCollectionNftImage             = "update_restricted_collection_nft_image"
	TypeMsgUpdateRestrictedCollectionNftImageBatch        = "update_restricted_collection_nft_image_batch"
	TypeMsgUpdateRestrictedCollectionNftImageGroupedBatch = "update_restricted_collection_nft_image_grouped_batch"
)

var (
	_ legacytx.LegacyMsg = &MsgMintNfts{}
	_ legacytx.LegacyMsg = &MsgBurnNfts{}
	_ legacytx.LegacyMsg = &MsgTransferNfts{}
	_ legacytx.LegacyMsg = &MsgApproveNfts{}
	_ legacytx.LegacyMsg = &MsgApproveAllNfts{}
	_ legacytx.LegacyMsg = &MsgMintNft{}
	_ legacytx.LegacyMsg = &MsgBurnNft{}
	_ legacytx.LegacyMsg = &MsgTransferNft{}
	_ legacytx.LegacyMsg = &MsgApproveNft{}
	_ legacytx.LegacyMsg = &MsgUpdateGuardSoulBondNftImage{}
	_ legacytx.LegacyMsg = &MsgUpdateRestrictedCollectionNftImage{}
	_ legacytx.LegacyMsg = &MsgUpdateRestrictedCollectionNftImageBatch{}
	_ legacytx.LegacyMsg = &MsgUpdateRestrictedCollectionNftImageGroupedBatch{}

	_ sdk.Msg = &MsgMintNfts{}
	_ sdk.Msg = &MsgBurnNfts{}
	_ sdk.Msg = &MsgTransferNfts{}
	_ sdk.Msg = &MsgApproveNfts{}
	_ sdk.Msg = &MsgApproveAllNfts{}
	_ sdk.Msg = &MsgMintNft{}
	_ sdk.Msg = &MsgBurnNft{}
	_ sdk.Msg = &MsgTransferNft{}
	_ sdk.Msg = &MsgApproveNft{}
	_ sdk.Msg = &MsgUpdateGuardSoulBondNftImage{}
	_ sdk.Msg = &MsgUpdateRestrictedCollectionNftImage{}
	_ sdk.Msg = &MsgUpdateRestrictedCollectionNftImageBatch{}
	_ sdk.Msg = &MsgUpdateRestrictedCollectionNftImageGroupedBatch{}
)

func NewMsgUpdateRestrictedCollectionNftImageGroupedBatch(
	creator string, collectionCreator string, collectionId string, nftsImagesGrouped MsgNftsImagesGroupedMetadata,
) *MsgUpdateRestrictedCollectionNftImageGroupedBatch {
	return &MsgUpdateRestrictedCollectionNftImageGroupedBatch{
		Creator:           creator,
		CollectionCreator: collectionCreator,
		CollectionId:      collectionId,
		NftsImagesGrouped: &nftsImagesGrouped,
	}
}

func (msg *MsgUpdateRestrictedCollectionNftImageGroupedBatch) Route() string {
	return RouterKey
}

func (msg *MsgUpdateRestrictedCollectionNftImageGroupedBatch) Type() string {
	return TypeMsgUpdateRestrictedCollectionNftImageGroupedBatch
}

func (msg *MsgUpdateRestrictedCollectionNftImageGroupedBatch) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateRestrictedCollectionNftImageGroupedBatch) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if strings.TrimSpace(msg.CollectionCreator) == "" {
		return errors.Wrap(errorstypes.ErrInvalidRequest, "collection creator should not be empty")
	}
	if strings.TrimSpace(msg.CollectionId) == "" {
		return errors.Wrap(errorstypes.ErrInvalidRequest, "collection id should not be empty")
	}
	if msg.NftsImagesGrouped == nil ||
		len(msg.NftsImagesGrouped.NftsIdsGrouped) == 0 ||
		msg.NftsImagesGrouped.Images == nil ||
		len(msg.NftsImagesGrouped.Images) == 0 {
		return errors.Wrapf(errorstypes.ErrKeyNotFound, "nfts ids/nfts images are empty")
	}
	if len(msg.NftsImagesGrouped.NftsIdsGrouped) != len(msg.NftsImagesGrouped.Images) {
		return errors.Wrapf(errorstypes.ErrInvalidRequest, "nfts ids/nfts images length is not equal")
	}
	for i := range msg.NftsImagesGrouped.NftsIdsGrouped {
		if msg.NftsImagesGrouped.NftsIdsGrouped[i] == nil || len(msg.NftsImagesGrouped.NftsIdsGrouped[i].NftsIds) == 0 {
			return errors.Wrapf(errorstypes.ErrKeyNotFound, "nft id is empty")
		}
		if msg.NftsImagesGrouped.Images[i] == nil {
			return errors.Wrapf(errorstypes.ErrKeyNotFound, "image is empty")
		}
		if strings.TrimSpace(msg.NftsImagesGrouped.Images[i].Type) == "" {
			return errors.Wrapf(errorstypes.ErrKeyNotFound, "image type is empty")
		}
		if strings.TrimSpace(msg.NftsImagesGrouped.Images[i].Url) == "" {
			return errors.Wrapf(errorstypes.ErrKeyNotFound, "image url is empty")
		}
	}
	return nil
}

func NewMsgUpdateRestrictedCollectionNftImageBatch(
	creator string, collectionCreator string, collectionId string, nftsImages MsgNftsImagesMetadata,
) *MsgUpdateRestrictedCollectionNftImageBatch {
	return &MsgUpdateRestrictedCollectionNftImageBatch{
		Creator:           creator,
		CollectionCreator: collectionCreator,
		CollectionId:      collectionId,
		NftsImages:        &nftsImages,
	}
}

func (msg *MsgUpdateRestrictedCollectionNftImageBatch) Route() string {
	return RouterKey
}

func (msg *MsgUpdateRestrictedCollectionNftImageBatch) Type() string {
	return TypeMsgUpdateRestrictedCollectionNftImageBatch
}

func (msg *MsgUpdateRestrictedCollectionNftImageBatch) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateRestrictedCollectionNftImageBatch) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if strings.TrimSpace(msg.CollectionCreator) == "" {
		return errors.Wrap(errorstypes.ErrInvalidRequest, "collection creator should not be empty")
	}
	if strings.TrimSpace(msg.CollectionId) == "" {
		return errors.Wrap(errorstypes.ErrInvalidRequest, "collection id should not be empty")
	}
	if msg.NftsImages == nil ||
		len(msg.NftsImages.NftsIds) == 0 ||
		msg.NftsImages.Images == nil ||
		len(msg.NftsImages.Images) == 0 {
		return errors.Wrapf(errorstypes.ErrKeyNotFound, "nfts ids/nfts images are empty")
	}
	if len(msg.NftsImages.NftsIds) != len(msg.NftsImages.Images) {
		return errors.Wrapf(errorstypes.ErrInvalidRequest, "nfts ids/nfts images length is not equal")
	}
	for i := range msg.NftsImages.NftsIds {
		if strings.TrimSpace(msg.NftsImages.NftsIds[i]) == "" {
			return errors.Wrapf(errorstypes.ErrKeyNotFound, "nft id is empty")
		}
		if msg.NftsImages.Images[i] == nil {
			return errors.Wrapf(errorstypes.ErrKeyNotFound, "image is empty")
		}
		if strings.TrimSpace(msg.NftsImages.Images[i].Type) == "" {
			return errors.Wrapf(errorstypes.ErrKeyNotFound, "image type is empty")
		}
		if strings.TrimSpace(msg.NftsImages.Images[i].Url) == "" {
			return errors.Wrapf(errorstypes.ErrKeyNotFound, "image url is empty")
		}
	}
	return nil
}

func NewMsgUpdateRestrictedCollectionNftImage(
	creator string,
	owner string,
	collectionCreator string,
	collectionId string,
	nftId string,
	index uint64,
	image *MsgNftImageMetadata,
) *MsgUpdateRestrictedCollectionNftImage {
	return &MsgUpdateRestrictedCollectionNftImage{
		Creator:           creator,
		Owner:             owner,
		NftId:             nftId,
		Index:             index,
		Image:             image,
		CollectionCreator: collectionCreator,
		CollectionId:      collectionId,
	}
}

func (msg *MsgUpdateRestrictedCollectionNftImage) Route() string {
	return RouterKey
}

func (msg *MsgUpdateRestrictedCollectionNftImage) Type() string {
	return TypeMsgUpdateRestrictedCollectionNftImage
}

func (msg *MsgUpdateRestrictedCollectionNftImage) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateRestrictedCollectionNftImage) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid owner address (%s)", err)
	}
	if strings.TrimSpace(msg.CollectionCreator) == "" {
		return errors.Wrap(errorstypes.ErrInvalidRequest, "collection creator should not be empty")
	}
	if strings.TrimSpace(msg.CollectionId) == "" {
		return errors.Wrap(errorstypes.ErrInvalidRequest, "collection id should not be empty")
	}
	if msg.Image == nil {
		return errors.Wrapf(errorstypes.ErrKeyNotFound, "image is empty")
	}
	if msg.Image.Image == nil {
		return errors.Wrapf(errorstypes.ErrKeyNotFound, "image is empty")
	}
	if strings.TrimSpace(msg.Image.Image.Type) == "" {
		return errors.Wrapf(errorstypes.ErrKeyNotFound, "image type is empty")
	}
	if strings.TrimSpace(msg.Image.Image.Url) == "" {
		return errors.Wrapf(errorstypes.ErrKeyNotFound, "image url is empty")
	}
	if strings.TrimSpace(msg.NftId) == "" {
		return errors.Wrapf(errorstypes.ErrKeyNotFound, "nft id is empty")
	}
	return nil
}

func NewMsgUpdateGuardSoulBondNftImage(creator string, owner string, nftId string, index uint64,
	image *MsgNftImageMetadata,
) *MsgUpdateGuardSoulBondNftImage {
	return &MsgUpdateGuardSoulBondNftImage{
		Creator: creator,
		Owner:   owner,
		NftId:   nftId,
		Index:   index,
		Image:   image,
	}
}

func (msg *MsgUpdateGuardSoulBondNftImage) Route() string {
	return RouterKey
}

func (msg *MsgUpdateGuardSoulBondNftImage) Type() string {
	return TypeMsgUpdateGuardSoulBondNftImage
}

func (msg *MsgUpdateGuardSoulBondNftImage) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateGuardSoulBondNftImage) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid owner address (%s)", err)
	}
	if msg.Image == nil {
		return errors.Wrapf(errorstypes.ErrKeyNotFound, "image is empty")
	}
	if msg.Image.Image == nil {
		return errors.Wrapf(errorstypes.ErrKeyNotFound, "image is empty")
	}
	if strings.TrimSpace(msg.Image.Image.Type) == "" {
		return errors.Wrapf(errorstypes.ErrKeyNotFound, "image type is empty")
	}
	if strings.TrimSpace(msg.Image.Image.Url) == "" {
		return errors.Wrapf(errorstypes.ErrKeyNotFound, "image url is empty")
	}
	if strings.TrimSpace(msg.NftId) == "" {
		return errors.Wrapf(errorstypes.ErrKeyNotFound, "nft id is empty")
	}
	return nil
}

func NewMsgMintNfts(creator string, collectionCreator string, collectionId string,
	nfts *MsgNftsMetadata,
	receiver string,
	strict bool,
	did bool,
) *MsgMintNfts {
	return &MsgMintNfts{
		Creator:           creator,
		CollectionCreator: collectionCreator,
		CollectionId:      collectionId,
		Nfts:              nfts,
		Receiver:          receiver,
		Strict:            strict,
		Did:               did,
	}
}

func (msg *MsgMintNfts) Route() string {
	return RouterKey
}

func (msg *MsgMintNfts) Type() string {
	return TypeMsgMintNfts
}

func (msg *MsgMintNfts) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgMintNfts) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if strings.TrimSpace(msg.Receiver) != "" {
		_, err = sdk.AccAddressFromBech32(msg.Receiver)
		if err != nil {
			return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid receiver address (%s)", err)
		}
	}
	if msg.Strict {
		if strings.TrimSpace(msg.CollectionCreator) == "" {
			return errors.Wrap(errorstypes.ErrInvalidRequest, "collection creator should not be empty with strict flag")
		}
		if strings.TrimSpace(msg.CollectionId) == "" {
			return errors.Wrap(errorstypes.ErrInvalidRequest, "collection id should not be empty with strict flag")
		}
	}
	if strings.TrimSpace(msg.CollectionCreator) != "" {
		_, err := sdk.AccAddressFromBech32(msg.CollectionCreator)
		if err != nil {
			return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid collection creator address (%s)", err)
		}
	}
	if msg.Nfts == nil || len(msg.Nfts.Nfts) == 0 {
		return errors.Wrapf(errorstypes.ErrKeyNotFound, "nfts are empty")
	}
	return nil
}

func NewMsgBurnNfts(creator string, collectionCreator string, collectionId string,
	nftsIds *MsgNftsIds,
	strict bool,
) *MsgBurnNfts {
	return &MsgBurnNfts{
		Creator:           creator,
		CollectionCreator: collectionCreator,
		CollectionId:      collectionId,
		Nfts:              nftsIds,
		Strict:            strict,
	}
}

func (msg *MsgBurnNfts) Route() string {
	return RouterKey
}

func (msg *MsgBurnNfts) Type() string {
	return TypeMsgBurnNfts
}

func (msg *MsgBurnNfts) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgBurnNfts) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.Strict {
		if strings.TrimSpace(msg.CollectionCreator) == "" {
			return errors.Wrap(errorstypes.ErrInvalidRequest, "collection creator should not be empty with strict flag")
		}
		if strings.TrimSpace(msg.CollectionId) == "" {
			return errors.Wrap(errorstypes.ErrInvalidRequest, "collection id should not be empty with strict flag")
		}
	}
	if strings.TrimSpace(msg.CollectionCreator) != "" {
		_, err := sdk.AccAddressFromBech32(msg.CollectionCreator)
		if err != nil {
			return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid collection creator address (%s)", err)
		}
	}
	if msg.Nfts == nil || len(msg.Nfts.NftsIds) == 0 {
		return errors.Wrapf(errorstypes.ErrKeyNotFound, "nfts ids are empty")
	}
	return nil
}

func NewMsgTransferNfts(creator string, collectionCreator string, collectionId string,
	nftsIds *MsgNftsIds,
	owner string,
	receiver string,
	strict bool,
) *MsgTransferNfts {
	return &MsgTransferNfts{
		Creator:           creator,
		Owner:             owner,
		Receiver:          receiver,
		CollectionCreator: collectionCreator,
		CollectionId:      collectionId,
		Nfts:              nftsIds,
		Strict:            strict,
	}
}

func (msg *MsgTransferNfts) Route() string {
	return RouterKey
}

func (msg *MsgTransferNfts) Type() string {
	return TypeMsgTransferNfts
}

func (msg *MsgTransferNfts) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgTransferNfts) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid receiver owner (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid receiver address (%s)", err)
	}
	if msg.Strict {
		if strings.TrimSpace(msg.CollectionCreator) == "" {
			return errors.Wrap(errorstypes.ErrInvalidRequest, "collection creator should not be empty with strict flag")
		}
		if strings.TrimSpace(msg.CollectionId) == "" {
			return errors.Wrap(errorstypes.ErrInvalidRequest, "collection id should not be empty with strict flag")
		}
	}
	if strings.TrimSpace(msg.CollectionCreator) != "" {
		_, err := sdk.AccAddressFromBech32(msg.CollectionCreator)
		if err != nil {
			return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid collection creator address (%s)", err)
		}
	}
	if msg.Nfts == nil || len(msg.Nfts.NftsIds) == 0 {
		return errors.Wrapf(errorstypes.ErrKeyNotFound, "nfts ids are empty")
	}
	return nil
}

func NewMsgApproveNfts(creator string, receiver string, collectionCreator string, collectionId string,
	nftsIds *MsgNftsIds,
	approved bool,
	strict bool,
) *MsgApproveNfts {
	return &MsgApproveNfts{
		Creator:           creator,
		Receiver:          receiver,
		CollectionCreator: collectionCreator,
		CollectionId:      collectionId,
		Nfts:              nftsIds,
		Approved:          approved,
		Strict:            strict,
	}
}

func (msg *MsgApproveNfts) Route() string {
	return RouterKey
}

func (msg *MsgApproveNfts) Type() string {
	return TypeMsgApproveNfts
}

func (msg *MsgApproveNfts) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgApproveNfts) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid receiver address (%s)", err)
	}
	if msg.Strict {
		if strings.TrimSpace(msg.CollectionCreator) == "" {
			return errors.Wrap(errorstypes.ErrInvalidRequest, "collection creator should not be empty with strict flag")
		}
		if strings.TrimSpace(msg.CollectionId) == "" {
			return errors.Wrap(errorstypes.ErrInvalidRequest, "collection id should not be empty with strict flag")
		}
	}
	if strings.TrimSpace(msg.CollectionCreator) != "" {
		_, err := sdk.AccAddressFromBech32(msg.CollectionCreator)
		if err != nil {
			return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid collection creator address (%s)", err)
		}
	}
	if msg.Nfts == nil || len(msg.Nfts.NftsIds) == 0 {
		return errors.Wrapf(errorstypes.ErrKeyNotFound, "nfts ids are empty")
	}
	return nil
}

func NewMsgApproveAllNfts(creator string, receiver string,
	approved bool) *MsgApproveAllNfts {
	return &MsgApproveAllNfts{
		Creator:  creator,
		Receiver: receiver,
		Approved: approved,
	}
}

func (msg *MsgApproveAllNfts) Route() string {
	return RouterKey
}

func (msg *MsgApproveAllNfts) Type() string {
	return TypeMsgApproveAllNfts
}

func (msg *MsgApproveAllNfts) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgApproveAllNfts) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid receiver address (%s)", err)
	}
	return nil
}

func NewMsgMintNft(creator string, collectionCreator string, collectionId string,
	nft *MsgNftMetadata,
	receiver string,
	strict bool,
	did bool,
) *MsgMintNft {
	return &MsgMintNft{
		Creator:           creator,
		CollectionCreator: collectionCreator,
		CollectionId:      collectionId,
		Nft:               nft,
		Receiver:          receiver,
		Strict:            strict,
		Did:               did,
	}
}

func (msg *MsgMintNft) Route() string {
	return RouterKey
}

func (msg *MsgMintNft) Type() string {
	return TypeMsgMintNft
}

func (msg *MsgMintNft) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgMintNft) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if strings.TrimSpace(msg.Receiver) != "" {
		_, err = sdk.AccAddressFromBech32(msg.Receiver)
		if err != nil {
			return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid receiver address (%s)", err)
		}
	}
	if msg.Strict {
		if strings.TrimSpace(msg.CollectionCreator) == "" {
			return errors.Wrap(errorstypes.ErrInvalidRequest, "collection creator should not be empty with strict flag")
		}
		if strings.TrimSpace(msg.CollectionId) == "" {
			return errors.Wrap(errorstypes.ErrInvalidRequest, "collection id should not be empty with strict flag")
		}
	}
	if strings.TrimSpace(msg.CollectionCreator) != "" {
		_, err := sdk.AccAddressFromBech32(msg.CollectionCreator)
		if err != nil {
			return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid collection creator address (%s)", err)
		}
	}
	if msg.Nft == nil {
		return errors.Wrapf(errorstypes.ErrKeyNotFound, "nft is empty")
	}
	return nil
}

func NewMsgBurnNft(creator string, collectionCreator string, collectionId string,
	nftId string,
	strict bool,
) *MsgBurnNft {
	return &MsgBurnNft{
		Creator:           creator,
		CollectionCreator: collectionCreator,
		CollectionId:      collectionId,
		NftId:             nftId,
		Strict:            strict,
	}
}

func (msg *MsgBurnNft) Route() string {
	return RouterKey
}

func (msg *MsgBurnNft) Type() string {
	return TypeMsgBurnNft
}

func (msg *MsgBurnNft) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgBurnNft) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.Strict {
		if strings.TrimSpace(msg.CollectionCreator) == "" {
			return errors.Wrap(errorstypes.ErrInvalidRequest, "collection creator should not be empty with strict flag")
		}
		if strings.TrimSpace(msg.CollectionId) == "" {
			return errors.Wrap(errorstypes.ErrInvalidRequest, "collection id should not be empty with strict flag")
		}
	}
	if strings.TrimSpace(msg.CollectionCreator) != "" {
		_, err := sdk.AccAddressFromBech32(msg.CollectionCreator)
		if err != nil {
			return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid collection creator address (%s)", err)
		}
	}
	if strings.TrimSpace(msg.NftId) == "" {
		return errors.Wrapf(errorstypes.ErrKeyNotFound, "nft id is empty")
	}
	return nil
}

func NewMsgTransferNft(creator string, collectionCreator string, collectionId string,
	nftId string,
	owner string,
	receiver string,
	strict bool,
) *MsgTransferNft {
	return &MsgTransferNft{
		Creator:           creator,
		Owner:             owner,
		Receiver:          receiver,
		CollectionCreator: collectionCreator,
		CollectionId:      collectionId,
		NftId:             nftId,
		Strict:            strict,
	}
}

func (msg *MsgTransferNft) Route() string {
	return RouterKey
}

func (msg *MsgTransferNft) Type() string {
	return TypeMsgTransferNft
}

func (msg *MsgTransferNft) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgTransferNft) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid receiver owner (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid receiver address (%s)", err)
	}
	if msg.Strict {
		if strings.TrimSpace(msg.CollectionCreator) == "" {
			return errors.Wrap(errorstypes.ErrInvalidRequest, "collection creator should not be empty with strict flag")
		}
		if strings.TrimSpace(msg.CollectionId) == "" {
			return errors.Wrap(errorstypes.ErrInvalidRequest, "collection id should not be empty with strict flag")
		}
	}
	if strings.TrimSpace(msg.CollectionCreator) != "" {
		_, err := sdk.AccAddressFromBech32(msg.CollectionCreator)
		if err != nil {
			return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid collection creator address (%s)", err)
		}
	}
	if strings.TrimSpace(msg.NftId) == "" {
		return errors.Wrapf(errorstypes.ErrKeyNotFound, "nft id is empty")
	}
	return nil
}

func NewMsgApproveNft(creator string, receiver string, collectionCreator string, collectionId string,
	nftId string,
	approved bool,
	strict bool,
) *MsgApproveNft {
	return &MsgApproveNft{
		Creator:           creator,
		Receiver:          receiver,
		CollectionCreator: collectionCreator,
		CollectionId:      collectionId,
		NftId:             nftId,
		Approved:          approved,
		Strict:            strict,
	}
}

func (msg *MsgApproveNft) Route() string {
	return RouterKey
}

func (msg *MsgApproveNft) Type() string {
	return TypeMsgApproveNft
}

func (msg *MsgApproveNft) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgApproveNft) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid receiver address (%s)", err)
	}
	if msg.Strict {
		if strings.TrimSpace(msg.CollectionCreator) == "" {
			return errors.Wrap(errorstypes.ErrInvalidRequest, "collection creator should not be empty with strict flag")
		}
		if strings.TrimSpace(msg.CollectionId) == "" {
			return errors.Wrap(errorstypes.ErrInvalidRequest, "collection id should not be empty with strict flag")
		}
	}
	if strings.TrimSpace(msg.CollectionCreator) != "" {
		_, err := sdk.AccAddressFromBech32(msg.CollectionCreator)
		if err != nil {
			return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid collection creator address (%s)", err)
		}
	}
	if strings.TrimSpace(msg.NftId) == "" {
		return errors.Wrapf(errorstypes.ErrKeyNotFound, "nft id is empty")
	}
	return nil
}
