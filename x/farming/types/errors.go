package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// farming module sentinel errors
var (
	ErrInvalidPlanType                 = sdkerrors.Register(ModuleName, 1301, "invalid plan type")
	ErrInvalidPlanName                 = sdkerrors.Register(ModuleName, 1302, "invalid plan name")
	ErrInvalidPlanEndTime              = sdkerrors.Register(ModuleName, 1303, "invalid plan end time")
	ErrInvalidStakingCoinWeights       = sdkerrors.Register(ModuleName, 1304, "invalid staking coin weights")
	ErrInvalidTotalEpochRatio          = sdkerrors.Register(ModuleName, 1305, "invalid total epoch ratio")
	ErrStakingNotExists                = sdkerrors.Register(ModuleName, 1306, "staking not exists")
	ErrConflictPrivatePlanFarmingPool  = sdkerrors.Register(ModuleName, 1307, "the address is already in use, please use a different plan name")
	ErrInvalidStakingReservedAmount    = sdkerrors.Register(ModuleName, 1308, "staking reserved amount invariant broken")
	ErrInvalidRemainingRewardsAmount   = sdkerrors.Register(ModuleName, 1309, "remaining rewards amount invariant broken")
	ErrInvalidOutstandingRewardsAmount = sdkerrors.Register(ModuleName, 1310, "outstanding rewards amount invariant broken")
	ErrNumPrivatePlansLimit            = sdkerrors.Register(ModuleName, 1311, "cannot create private plans more than the limit")
	ErrNumMaxDenomsLimit               = sdkerrors.Register(ModuleName, 1312, "number of denoms cannot exceed the limit")
	ErrInvalidEpochAmount              = sdkerrors.Register(ModuleName, 1313, "invalid epoch amount")
	ErrRatioPlanDisabled               = sdkerrors.Register(ModuleName, 1314, "creation of ratio plans is disabled")
	ErrInvalidUnharvestedRewardsAmount = sdkerrors.Register(ModuleName, 1315, "invalid unharvested rewards amount")
	ErrModuleDisabled                  = sdkerrors.Register(ModuleName, 1316, "farming module has been disabled")
)
