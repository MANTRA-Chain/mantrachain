package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgStartEpoch = "start_epoch"

var _ sdk.Msg = &MsgStartEpoch{}

func NewMsgStartEpoch(creator string, blockStart int64, reward string,
	stakingChain string, stakingValidator string) *MsgStartEpoch {
	return &MsgStartEpoch{
		Creator:          creator,
		BlockStart:       blockStart,
		Reward:           reward,
		StakingChain:     stakingChain,
		StakingValidator: stakingValidator,
	}
}

func (msg *MsgStartEpoch) Route() string {
	return RouterKey
}

func (msg *MsgStartEpoch) Type() string {
	return TypeMsgStartEpoch
}

func (msg *MsgStartEpoch) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgStartEpoch) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgStartEpoch) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.BlockStart <= 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "block start should be positive")
	}
	if strings.TrimSpace(msg.StakingChain) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "staking chain should not be empty")
	}
	if strings.TrimSpace(msg.StakingValidator) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "staking validator should not be empty")
	}
	return nil
}
