package post

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	xfeemarkettypes "github.com/MANTRA-Chain/mantrachain/x/xfeemarket/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/skip-mev/feemarket/x/feemarket/ante"
	feemarkettypes "github.com/skip-mev/feemarket/x/feemarket/types"
)

// BankSendGasConsumption is the gas consumption of the bank sends that occur during feemarket handler execution.
// 37325 is the additional gas consumed for burning the fees instead of just transferring them.
const BankSendGasConsumption = 12490 + 37325

// FeeMarketDeductDecorator deducts fees from the fee payer based off of the current state of the feemarket.
// The fee payer is the fee granter (if specified) or first signer of the tx.
// If the fee payer does not have the funds to pay for the fees, return an InsufficientFunds error.
// Excess between the given fee and the on-chain min base fee is refunded to payer.
// Call next PostHandler if fees successfully deducted.
// CONTRACT: Tx must implement FeeTx interface
type FeeMarketDeductDecorator struct {
	bankKeeper      BankKeeper
	feemarketKeeper FeeMarketKeeper
}

func NewFeeMarketDeductDecorator(bk BankKeeper, fmk FeeMarketKeeper) FeeMarketDeductDecorator {
	return FeeMarketDeductDecorator{
		bankKeeper:      bk,
		feemarketKeeper: fmk,
	}
}

// PostHandle deducts the fee from the fee payer based on the min base fee and the gas consumed in the gasmeter.
// If there is a difference between the provided fee and the min-base fee, the difference is paid as a tip.
// Fees are sent to the x/feemarket fee-collector address.
func (dfd FeeMarketDeductDecorator) PostHandle(ctx sdk.Context, tx sdk.Tx, simulate, success bool, next sdk.PostHandler) (sdk.Context, error) {
	// GenTx consume no fee
	if ctx.BlockHeight() == 0 {
		return next(ctx, tx, simulate, success)
	}

	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	if !simulate && ctx.BlockHeight() > 0 && feeTx.GetGas() == 0 {
		return ctx, errorsmod.Wrap(sdkerrors.ErrInvalidGasLimit, "must provide positive gas")
	}

	// update fee market params
	params, err := dfd.feemarketKeeper.GetParams(ctx)
	if err != nil {
		return ctx, errorsmod.Wrapf(err, "unable to get fee market params")
	}

	// return if disabled
	if !params.Enabled {
		return next(ctx, tx, simulate, success)
	}

	enabledHeight, err := dfd.feemarketKeeper.GetEnabledHeight(ctx)
	if err != nil {
		return ctx, errorsmod.Wrapf(err, "unable to get fee market enabled height")
	}

	// if the current height is that which enabled the feemarket or lower, skip deduction
	if ctx.BlockHeight() <= enabledHeight {
		return next(ctx, tx, simulate, success)
	}

	// update fee market state
	state, err := dfd.feemarketKeeper.GetState(ctx)
	if err != nil {
		return ctx, errorsmod.Wrapf(err, "unable to get fee market state")
	}

	feeCoins := feeTx.GetFee()
	gas := ctx.GasMeter().GasConsumed() // use context gas consumed

	if len(feeCoins) == 0 && !simulate {
		return ctx, errorsmod.Wrapf(feemarkettypes.ErrNoFeeCoins, "got length %d", len(feeCoins))
	}
	if len(feeCoins) > 1 {
		return ctx, errorsmod.Wrapf(feemarkettypes.ErrTooManyFeeCoins, "got length %d", len(feeCoins))
	}

	// if simulating and user did not provider a fee - create a dummy value for them
	var (
		tip     = sdk.NewCoin(params.FeeDenom, math.ZeroInt())
		payCoin = sdk.NewCoin(params.FeeDenom, math.ZeroInt())
	)
	if !simulate {
		payCoin = feeCoins[0]
	}

	feeGas := int64(feeTx.GetGas())

	minGasPrice, err := dfd.feemarketKeeper.GetMinGasPrice(ctx, payCoin.GetDenom())
	if err != nil {
		return ctx, errorsmod.Wrapf(err, "unable to get min gas price for denom %s", payCoin.GetDenom())
	}

	ctx.Logger().Info("fee deduct post handle",
		"min gas prices", minGasPrice,
		"gas consumed", gas,
	)

	if !simulate {
		payCoin, tip, err = ante.CheckTxFee(ctx, minGasPrice, payCoin, feeGas, false)
		if err != nil {
			return ctx, err
		}
	}

	ctx.Logger().Info("fee deduct post handle",
		"fee", payCoin,
		"tip", tip,
	)

	feePayer := feeTx.FeePayer()

	if err := dfd.BurnFeeAndRefund(ctx, payCoin, tip, feePayer, params.FeeDenom); err != nil {
		return ctx, err
	}

	err = state.Update(gas, params)
	if err != nil {
		return ctx, errorsmod.Wrapf(err, "unable to update fee market state")
	}

	err = dfd.feemarketKeeper.SetState(ctx, state)
	if err != nil {
		return ctx, errorsmod.Wrapf(err, "unable to set fee market state")
	}

	if simulate {
		// consume the gas that would be consumed during normal execution
		ctx.GasMeter().ConsumeGas(BankSendGasConsumption, "simulation send gas consumption")
	}

	return next(ctx, tx, simulate, success)
}

// BurnFeeAndRefund burns the fees and refunds the extra/tip to the fee payer.
func (dfd FeeMarketDeductDecorator) BurnFeeAndRefund(ctx sdk.Context, fee, tip sdk.Coin, feePayer sdk.AccAddress, defaultFeeDenom string) error {
	var events sdk.Events

	// burn the fees if it is the default fee denom
	if !fee.IsNil() && !fee.IsZero() && fee.Denom == defaultFeeDenom {
		err := BurnCoins(dfd.bankKeeper, ctx, sdk.NewCoins(fee))
		if err != nil {
			return err
		}

		events = append(events, sdk.NewEvent(
			feemarkettypes.EventTypeFeePay,
			sdk.NewAttribute(sdk.AttributeKeyFee, fee.String()),
		))
	}

	// refund the tip if it is not nil and non zero
	if !tip.IsNil() && !tip.IsZero() {
		err := RefundTip(dfd.bankKeeper, ctx, feePayer, sdk.NewCoins(tip))
		if err != nil {
			return err
		}

		events = append(events, sdk.NewEvent(
			xfeemarkettypes.EventTypeTipRefund,
			sdk.NewAttribute(feemarkettypes.AttributeKeyTip, tip.String()),
			sdk.NewAttribute(feemarkettypes.AttributeKeyTipPayee, feePayer.String()),
		))
	}

	ctx.EventManager().EmitEvents(events)
	return nil
}

// BurnCoins burns coins from the fee collector account.
func BurnCoins(bankKeeper BankKeeper, ctx sdk.Context, coins sdk.Coins) error {
	err := bankKeeper.BurnCoins(ctx, feemarkettypes.FeeCollectorName, coins)
	if err != nil {
		return err
	}
	return nil
}

// RefundTip sends a tip to the txfee payer.
func RefundTip(bankKeeper BankKeeper, ctx sdk.Context, feePayer sdk.AccAddress, coins sdk.Coins) error {
	err := bankKeeper.SendCoinsFromModuleToAccount(ctx, feemarkettypes.FeeCollectorName, feePayer, coins)
	if err != nil {
		return err
	}

	return nil
}
