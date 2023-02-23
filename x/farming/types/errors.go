package types

import (
	"cosmossdk.io/errors"
)

// farming module sentinel errors
var (
	ErrInvalidPlanType                 = errors.Register(ModuleName, 2, "invalid plan type")
	ErrInvalidPlanName                 = errors.Register(ModuleName, 3, "invalid plan name")
	ErrInvalidPlanEndTime              = errors.Register(ModuleName, 4, "invalid plan end time")
	ErrInvalidStakingCoinWeights       = errors.Register(ModuleName, 5, "invalid staking coin weights")
	ErrInvalidTotalEpochRatio          = errors.Register(ModuleName, 6, "invalid total epoch ratio")
	ErrStakingNotExists                = errors.Register(ModuleName, 7, "staking not exists")
	ErrConflictPrivatePlanFarmingPool  = errors.Register(ModuleName, 8, "the address is already in use, please use a different plan name")
	ErrInvalidStakingReservedAmount    = errors.Register(ModuleName, 9, "staking reserved amount invariant broken")
	ErrInvalidRemainingRewardsAmount   = errors.Register(ModuleName, 10, "remaining rewards amount invariant broken")
	ErrInvalidOutstandingRewardsAmount = errors.Register(ModuleName, 11, "outstanding rewards amount invariant broken")
	ErrNumPrivatePlansLimit            = errors.Register(ModuleName, 12, "cannot create private plans more than the limit")
	ErrNumMaxDenomsLimit               = errors.Register(ModuleName, 13, "number of denoms cannot exceed the limit")
	ErrInvalidEpochAmount              = errors.Register(ModuleName, 14, "invalid epoch amount")
	ErrRatioPlanDisabled               = errors.Register(ModuleName, 15, "creation of ratio plans is disabled")
	ErrInvalidUnharvestedRewardsAmount = errors.Register(ModuleName, 16, "invalid unharvested rewards amount")
	ErrModuleDisabled                  = errors.Register(ModuleName, 17, "farming module has been disabled")
)
