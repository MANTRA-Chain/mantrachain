package keeper

import (
	"github.com/LimeChain/mantrachain/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type DistributionExecutor struct {
	ctx sdk.Context
	ac  types.AccountKeeper
	sk  types.StakingKeeper
	dk  types.DistrKeeper
}

func NewDistributionExecutor(ctx sdk.Context, ac types.AccountKeeper, sk types.StakingKeeper, dk types.DistrKeeper) *DistributionExecutor {
	return &DistributionExecutor{
		ctx: ctx,
		ac:  ac,
		sk:  sk,
		dk:  dk,
	}
}

func (c *DistributionExecutor) GetDelegationRewards(
	validator string,
) (sdk.DecCoins, error) {
	valAdr, err := sdk.ValAddressFromBech32(validator)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "invalid validator address")
	}

	val := c.sk.Validator(c.ctx, valAdr)
	if val == nil {
		return nil, sdkerrors.Wrap(err, "validator not exists")
	}

	del := c.sk.Delegation(c.ctx, c.ac.GetModuleAddress(types.ModuleName), valAdr)
	if del == nil {
		return nil, sdkerrors.Wrap(err, "delegation not exists")
	}

	endingPeriod := c.dk.IncrementValidatorPeriod(c.ctx, val)
	rewards := c.dk.CalculateDelegationRewards(c.ctx, val, del, endingPeriod)

	return rewards, nil
}

func (c *DistributionExecutor) WithdrawDelegationRewards(
	validator string,
) (sdk.Coins, error) {
	valAdr, err := sdk.ValAddressFromBech32(validator)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "invalid validator address")
	}

	rewards, err := c.GetDelegationRewards(validator)

	if err != nil {
		return nil, err
	}

	if len(rewards) == 0 {
		return nil, nil
	}

	hasReachedMinTreshold := false

	for _, v := range rewards {
		// checks minimum threshold for withdrawal delegation rewards
		// TODO: make this configurable
		if v.Amount.GTE(sdk.NewDecFromInt(sdk.NewInt(1))) {
			hasReachedMinTreshold = true
			break
		}
	}

	if !hasReachedMinTreshold {
		return nil, nil
	}

	withdrawn, err := c.dk.WithdrawDelegationRewards(c.ctx, c.ac.GetModuleAddress(types.ModuleName), valAdr)

	if err != nil {
		return nil, sdkerrors.Wrap(err, "cannot witdraw rewards")
	}

	return withdrawn, nil
}
