package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgWithdrawNftReward = "withdraw_nft_reward"

var _ sdk.Msg = &MsgWithdrawNftReward{}

func NewMsgWithdrawNftReward(creator string, marketplaceCreator string, marketplaceId string,
	collectionCreator string, collectionId string, nftId string, receiver string) *MsgWithdrawNftReward {
	return &MsgWithdrawNftReward{
		Creator:            creator,
		MarketplaceCreator: marketplaceCreator,
		MarketplaceId:      marketplaceId,
		CollectionCreator:  collectionCreator,
		CollectionId:       collectionId,
		NftId:              nftId,
		Receiver:           receiver,
	}
}

func (msg *MsgWithdrawNftReward) Route() string {
	return RouterKey
}

func (msg *MsgWithdrawNftReward) Type() string {
	return TypeMsgWithdrawNftReward
}

func (msg *MsgWithdrawNftReward) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgWithdrawNftReward) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgWithdrawNftReward) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if strings.TrimSpace(msg.Receiver) != "" {
		_, err = sdk.AccAddressFromBech32(msg.Receiver)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receiver address (%s)", err)
		}
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
	if strings.TrimSpace(msg.NftId) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "nft id should not be empty")
	}
	return nil
}
