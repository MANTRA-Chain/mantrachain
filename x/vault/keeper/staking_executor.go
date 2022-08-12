package keeper

import (
	"github.com/LimeChain/mantrachain/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sktypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type StakingExecutor struct {
	ctx sdk.Context
	ac  types.AccountKeeper
	bk  types.BankKeeper
	sk  types.StakingKeeper
}

func NewStakingExecutor(ctx sdk.Context, ac types.AccountKeeper, bk types.BankKeeper, sk types.StakingKeeper) *StakingExecutor {
	return &StakingExecutor{
		ctx: ctx,
		ac:  ac,
		bk:  bk,
		sk:  sk,
	}
}

func (c *StakingExecutor) Delegate(
	creator sdk.AccAddress,
	amount sdk.Coin,
	val string,
) (sdk.Dec, error) {
	shares := sdk.NewDec(0)
	valAddress, err := sdk.ValAddressFromBech32(val)

	if err != nil {
		return shares, err
	}

	err = c.bk.SendCoinsFromAccountToModule(c.ctx, creator, types.ModuleName, []sdk.Coin{amount})

	if err != nil {
		return shares, err
	}

	validator, found := c.sk.GetValidator(c.ctx, valAddress)

	if !found {
		return shares, sdkerrors.Wrapf(types.ErrValidatorDoesNotExist, "invalid or non-existent validator %s", valAddress)
	}

	shares, err = c.sk.Delegate(
		c.ctx,
		c.ac.GetModuleAddress(types.ModuleName),
		amount.Amount,
		sktypes.Unbonded,
		validator,
		true,
	)

	if err != nil {
		return shares, err
	}

	return shares, nil
}
