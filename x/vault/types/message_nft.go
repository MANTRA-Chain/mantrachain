package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgWithdrawNftRewards   = "withdraw_nft_rewards"
	TypeMsgUpdateNftStakeStaked = "update_nft_stake_staked"
)

var (
	_ sdk.Msg = &MsgWithdrawNftRewards{}
	_ sdk.Msg = &MsgUpdateNftStakeStaked{}
)

func NewMsgWithdrawNftRewards(creator string, marketplaceCreator string, marketplaceId string,
	collectionCreator string, collectionId string, nftId string, receiver string, stakingChain string, stakingValidator string) *MsgWithdrawNftRewards {
	return &MsgWithdrawNftRewards{
		Creator:            creator,
		MarketplaceCreator: marketplaceCreator,
		MarketplaceId:      marketplaceId,
		CollectionCreator:  collectionCreator,
		CollectionId:       collectionId,
		NftId:              nftId,
		Receiver:           receiver,
		StakingChain:       stakingChain,
		StakingValidator:   stakingValidator,
	}
}

func (msg *MsgWithdrawNftRewards) Route() string {
	return RouterKey
}

func (msg *MsgWithdrawNftRewards) Type() string {
	return TypeMsgWithdrawNftRewards
}

func (msg *MsgWithdrawNftRewards) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgWithdrawNftRewards) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgWithdrawNftRewards) ValidateBasic() error {
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

func NewMsgUpdateNftStakeStaked(creator string, marketplaceCreator string, marketplaceId string,
	collectionCreator string, collectionId string, nftId string,
	stakingChain string, stakingValidator string, blockHeight int64, stakedIndex int64, shares string,
) *MsgUpdateNftStakeStaked {
	return &MsgUpdateNftStakeStaked{
		Creator:            creator,
		MarketplaceCreator: marketplaceCreator,
		MarketplaceId:      marketplaceId,
		CollectionCreator:  collectionCreator,
		CollectionId:       collectionId,
		NftId:              nftId,
		StakingChain:       stakingChain,
		StakingValidator:   stakingValidator,
		BlockHeight:        blockHeight,
		StakedIndex:        stakedIndex,
		Shares:             shares,
	}
}

func (msg *MsgUpdateNftStakeStaked) Route() string {
	return RouterKey
}

func (msg *MsgUpdateNftStakeStaked) Type() string {
	return TypeMsgUpdateNftStakeStaked
}

func (msg *MsgUpdateNftStakeStaked) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateNftStakeStaked) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateNftStakeStaked) ValidateBasic() error {
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
	if strings.TrimSpace(msg.NftId) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "nft id should not be empty")
	}
	if strings.TrimSpace(msg.StakingChain) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "staking chain should not be empty")
	}
	if strings.TrimSpace(msg.StakingValidator) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "staking validator should not be empty")
	}
	if msg.BlockHeight <= 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "block height should be positive")
	}
	if msg.StakedIndex < 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "staked index should not be negative")
	}
	if strings.TrimSpace(msg.Shares) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "shares should not be empty")
	}
	return nil
}
