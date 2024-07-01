package types

import (
	"time"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
)

var (
	_ sdk.Msg = &MsgCreatePrivatePlan{}
	_ sdk.Msg = &MsgTerminatePrivatePlan{}
	_ sdk.Msg = &MsgFarm{}
	_ sdk.Msg = &MsgUnfarm{}
	_ sdk.Msg = &MsgHarvest{}

	_ legacytx.LegacyMsg = &MsgCreatePrivatePlan{}
	_ legacytx.LegacyMsg = &MsgTerminatePrivatePlan{}
	_ legacytx.LegacyMsg = &MsgFarm{}
	_ legacytx.LegacyMsg = &MsgHarvest{}
)

// Message types for the module
const (
	TypeMsgCreatePrivatePlan    = "create_private_plan"
	TypeMsgTerminatePrivatePlan = "terminate_private_plan"
	TypeMsgFarm                 = "farm"
	TypeMsgUnfarm               = "unfarm"
	TypeMsgHarvest              = "harvest"
)

// NewMsgCreatePrivatePlan creates a new MsgCreatePrivatePlan.
func NewMsgCreatePrivatePlan(
	creatorAddr sdk.AccAddress, description string, rewardAllocations []RewardAllocation,
	startTime, endTime time.Time) *MsgCreatePrivatePlan {
	return &MsgCreatePrivatePlan{
		Creator:           creatorAddr.String(),
		Description:       description,
		RewardAllocations: rewardAllocations,
		StartTime:         startTime,
		EndTime:           endTime,
	}
}

func (msg MsgCreatePrivatePlan) Route() string { return RouterKey }
func (msg MsgCreatePrivatePlan) Type() string  { return TypeMsgCreatePrivatePlan }

func (msg *MsgCreatePrivatePlan) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgCreatePrivatePlan) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid creator address: %v", err)
	}
	// Create a dummy plan with valid fields and utilize Validate() method
	// for user-provided data.
	validAddr := RewardsPoolAddress // Chose random valid address
	dummyPlan := NewPlan(
		1, msg.Description, validAddr, validAddr,
		msg.RewardAllocations, msg.StartTime, msg.EndTime, true)
	if err := dummyPlan.Validate(); err != nil {
		return errors.Wrap(errorstypes.ErrInvalidRequest, err.Error())
	}
	return nil
}

func (msg MsgCreatePrivatePlan) GetCreatorAddress() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgTerminatePrivatePlan creates a new MsgTerminatePrivatePlan.
func NewMsgTerminatePrivatePlan(
	creatorAddr sdk.AccAddress, planId uint64) *MsgTerminatePrivatePlan {
	return &MsgTerminatePrivatePlan{
		Creator: creatorAddr.String(),
		PlanId:  planId,
	}
}

func (msg MsgTerminatePrivatePlan) Route() string { return RouterKey }
func (msg MsgTerminatePrivatePlan) Type() string  { return TypeMsgTerminatePrivatePlan }

func (msg *MsgTerminatePrivatePlan) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgTerminatePrivatePlan) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid creator address: %v", err)
	}
	if msg.PlanId == 0 {
		return errors.Wrap(errorstypes.ErrInvalidRequest, "plan id must not be 0")
	}
	return nil
}

func (msg MsgTerminatePrivatePlan) GetCreatorAddress() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgFarm creates a new MsgFarm.
func NewMsgFarm(farmerAddr sdk.AccAddress, coin sdk.Coin) *MsgFarm {
	return &MsgFarm{
		Farmer: farmerAddr.String(),
		Coin:   coin,
	}
}

func (msg MsgFarm) Route() string { return RouterKey }
func (msg MsgFarm) Type() string  { return TypeMsgFarm }

func (msg *MsgFarm) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgFarm) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Farmer); err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid farmer address: %v", err)
	}
	if err := msg.Coin.Validate(); err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidRequest, "invalid coin: %v", err)
	}
	if !msg.Coin.IsPositive() {
		return errors.Wrapf(errorstypes.ErrInvalidRequest, "non-positive coin: %s", msg.Coin)
	}
	return nil
}

func (msg MsgFarm) GetFarmerAddress() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Farmer)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgUnfarm creates a new MsgUnfarm.
func NewMsgUnfarm(farmerAddr sdk.AccAddress, coin sdk.Coin) *MsgUnfarm {
	return &MsgUnfarm{
		Farmer: farmerAddr.String(),
		Coin:   coin,
	}
}

func (msg MsgUnfarm) Route() string { return RouterKey }
func (msg MsgUnfarm) Type() string  { return TypeMsgUnfarm }

func (msg *MsgUnfarm) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgUnfarm) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Farmer); err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid farmer address: %v", err)
	}
	if err := msg.Coin.Validate(); err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidRequest, "invalid coin: %v", err)
	}
	if !msg.Coin.IsPositive() {
		return errors.Wrapf(errorstypes.ErrInvalidRequest, "non-positive coin: %s", msg.Coin)
	}
	return nil
}

func (msg MsgUnfarm) GetFarmerAddress() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Farmer)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgHarvest creates a new MsgHarvest.
func NewMsgHarvest(farmerAddr sdk.AccAddress, denom string) *MsgHarvest {
	return &MsgHarvest{
		Farmer: farmerAddr.String(),
		Denom:  denom,
	}
}

func (msg MsgHarvest) Route() string { return RouterKey }
func (msg MsgHarvest) Type() string  { return TypeMsgHarvest }

func (msg *MsgHarvest) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgHarvest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Farmer); err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid farmer address: %v", err)
	}
	if err := sdk.ValidateDenom(msg.Denom); err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidRequest, "invalid denom: %v", err)
	}
	return nil
}

func (msg MsgHarvest) GetFarmerAddress() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Farmer)
	if err != nil {
		panic(err)
	}
	return addr
}
