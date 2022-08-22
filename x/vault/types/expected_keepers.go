package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sktypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type NFTKeeper interface {
	GetOwner(ctx sdk.Context, classID string, nftID string) sdk.AccAddress
}

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	// Methods imported from bank should be defined here
}

type StakingKeeper interface {
	Delegate(
		ctx sdk.Context, delAddr sdk.AccAddress, bondAmt sdk.Int, tokenSrc sktypes.BondStatus,
		validator sktypes.Validator, subtractAccount bool,
	) (newShares sdk.Dec, err error)
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator sktypes.Validator, found bool)
	Delegation(sdk.Context, sdk.AccAddress, sdk.ValAddress) sktypes.DelegationI
	Validator(sdk.Context, sdk.ValAddress) sktypes.ValidatorI
	GetDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (delegation sktypes.Delegation, found bool)
	// Methods imported from bank should be defined here
}

type DistrKeeper interface {
	IncrementValidatorPeriod(ctx sdk.Context, val sktypes.ValidatorI) uint64
	CalculateDelegationRewards(ctx sdk.Context, val sktypes.ValidatorI, del sktypes.DelegationI, endingPeriod uint64) (rewards sdk.DecCoins)
	WithdrawDelegationRewards(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdk.Coins, error)
	// Methods imported from bank should be defined here
}