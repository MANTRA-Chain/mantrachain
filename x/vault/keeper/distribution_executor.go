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
	denom string,
) (sdk.Dec, error) {
	rewards := sdk.NewDec(0)

	valAdr, err := sdk.ValAddressFromBech32(validator)
	if err != nil {
		return rewards, sdkerrors.Wrap(err, "invalid validator address")
	}

	val := c.sk.Validator(c.ctx, valAdr)
	if val == nil {
		return rewards, sdkerrors.Wrap(err, "validator not exists")
	}

	del := c.sk.Delegation(c.ctx, c.ac.GetModuleAddress(types.ModuleName), valAdr)
	if del == nil {
		return rewards, sdkerrors.Wrap(err, "delegation not exists")
	}

	endingPeriod := c.dk.IncrementValidatorPeriod(c.ctx, val)
	res := c.dk.CalculateDelegationRewards(c.ctx, val, del, endingPeriod)

	for _, v := range res {
		if v.Denom == denom {
			rewards = rewards.Add(v.Amount)
		}
	}

	return rewards, nil
}

func (c *DistributionExecutor) WithdrawDelegationRewards(
	validator string,
	denom string,
) (sdk.Coin, error) {
	withdrawn := sdk.NewCoin(denom, sdk.NewInt(0))

	valAdr, err := sdk.ValAddressFromBech32(validator)
	if err != nil {
		return withdrawn, sdkerrors.Wrap(err, "invalid validator address")
	}

	amount, err := c.dk.WithdrawDelegationRewards(c.ctx, c.ac.GetModuleAddress(types.ModuleName), valAdr)

	if err != nil {
		return withdrawn, sdkerrors.Wrap(err, "cannot witdraw reward")
	}

	for _, v := range amount {
		if v.Denom == denom {
			withdrawn = withdrawn.AddAmount(v.Amount)
		}
	}

	return withdrawn, nil
}
