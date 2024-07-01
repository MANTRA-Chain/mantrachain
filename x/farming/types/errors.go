package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/farming module sentinel errors
var (
	ErrInvalidSigner                   = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrInvalidPlanType                 = errors.Register(ModuleName, 1101, "invalid plan type")
	ErrInvalidPlanName                 = errors.Register(ModuleName, 1102, "invalid plan name")
	ErrInvalidPlanEndTime              = errors.Register(ModuleName, 1103, "invalid plan end time")
	ErrInvalidStakingCoinWeights       = errors.Register(ModuleName, 1104, "invalid staking coin weights")
	ErrInvalidTotalEpochRatio          = errors.Register(ModuleName, 1105, "invalid total epoch ratio")
	ErrStakingNotExists                = errors.Register(ModuleName, 1106, "staking not exists")
	ErrConflictPrivatePlanFarmingPool  = errors.Register(ModuleName, 1107, "the address is already in use, please use a different plan name")
	ErrInvalidStakingReservedAmount    = errors.Register(ModuleName, 1108, "staking reserved amount invariant broken")
	ErrInvalidRemainingRewardsAmount   = errors.Register(ModuleName, 1109, "remaining rewards amount invariant broken")
	ErrInvalidOutstandingRewardsAmount = errors.Register(ModuleName, 1110, "outstanding rewards amount invariant broken")
	ErrNumPrivatePlansLimit            = errors.Register(ModuleName, 1111, "cannot create private plans more than the limit")
	ErrNumMaxDenomsLimit               = errors.Register(ModuleName, 1112, "number of denoms cannot exceed the limit")
	ErrInvalidEpochAmount              = errors.Register(ModuleName, 1113, "invalid epoch amount")
	ErrRatioPlanDisabled               = errors.Register(ModuleName, 1114, "creation of ratio plans is disabled")
	ErrInvalidUnharvestedRewardsAmount = errors.Register(ModuleName, 1115, "invalid unharvested rewards amount")
	ErrModuleDisabled                  = errors.Register(ModuleName, 1116, "farming module has been disabled")
)
