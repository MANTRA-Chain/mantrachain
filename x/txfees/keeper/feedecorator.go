package keeper

import (
	"bytes"
	"fmt"

	"github.com/MANTRA-Finance/mantrachain/x/txfees/types"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// TxFeeChecker check if the provided fee is enough and returns the effective fee and tx priority,
// the effective fee should be deducted later, and the priority should be returned in abci response.
type TxFeeChecker func(ctx sdk.Context, tx sdk.Tx, baseDenom string, txfeesKeeper types.TxfeesKeeper, liquidityKeeper types.LiquidityKeeper) (sdk.Coins, error)

// DeductFeeDecorator deducts fees from the first signer of the tx
// If the first signer does not have the funds to pay for the fees, return with InsufficientFunds error
// Call next AnteHandler if fees successfully deducted
// CONTRACT: Tx must implement FeeTx interface to use DeductFeeDecorator
type DeductFeeDecorator struct {
	accountKeeper   types.AccountKeeper
	bankKeeper      authtypes.BankKeeper
	feegrantKeeper  types.FeegrantKeeper
	txFeeChecker    TxFeeChecker
	guardKeeper     types.GuardKeeper
	liquidityKeeper types.LiquidityKeeper
	txfeesKeeper    types.TxfeesKeeper
}

func NewDeductFeeDecorator(ak types.AccountKeeper, bk authtypes.BankKeeper, fk types.FeegrantKeeper, tfc TxFeeChecker, gk types.GuardKeeper, lk types.LiquidityKeeper, tfk types.TxfeesKeeper) DeductFeeDecorator {
	if tfc == nil {
		tfc = checkTxFeeWithValidatorMinGasPrices
	}

	return DeductFeeDecorator{
		accountKeeper:   ak,
		bankKeeper:      bk,
		feegrantKeeper:  fk,
		txFeeChecker:    tfc,
		guardKeeper:     gk,
		liquidityKeeper: lk,
		txfeesKeeper:    tfk,
	}
}

func (dfd DeductFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	conf := dfd.txfeesKeeper.GetParams(ctx)

	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, errors.Wrap(errorstypes.ErrTxDecode, "Tx must be a FeeTx")
	}

	if !simulate && ctx.BlockHeight() > 0 && feeTx.GetGas() == 0 {
		return ctx, errors.Wrap(errorstypes.ErrInvalidGasLimit, "must provide positive gas")
	}

	var (
		err error
	)

	fee := feeTx.GetFee()
	if !simulate {
		fee, err = dfd.txFeeChecker(ctx, tx, conf.BaseDenom, dfd.txfeesKeeper, dfd.liquidityKeeper)
		if err != nil {
			return ctx, err
		}
	}
	if err := dfd.checkDeductFee(ctx, tx, fee); err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate)
}

func (dfd DeductFeeDecorator) checkDeductFee(ctx sdk.Context, sdkTx sdk.Tx, fee sdk.Coins) error {
	feeTx, ok := sdkTx.(sdk.FeeTx)
	if !ok {
		return errors.Wrap(errorstypes.ErrTxDecode, "Tx must be a FeeTx")
	}

	if addr := dfd.accountKeeper.GetModuleAddress(authtypes.FeeCollectorName); addr == nil {
		return fmt.Errorf("fee collector module account (%s) has not been set", authtypes.FeeCollectorName)
	}

	feePayer := feeTx.FeePayer()
	feeGranter := feeTx.FeeGranter()
	deductFeesFrom := feePayer

	// if feegranter set deduct fee from feegranter account.
	// this works with only when feegrant enabled.
	if feeGranter != nil {
		if dfd.feegrantKeeper == nil {
			return errorstypes.ErrInvalidRequest.Wrap("fee grants are not enabled")
		} else if !bytes.Equal(feeGranter, feePayer) {
			err := dfd.feegrantKeeper.UseGrantedFees(ctx, feeGranter, feePayer, fee, sdkTx.GetMsgs())
			if err != nil {
				return errors.Wrapf(err, "%s does not allow to pay fees for %s", feeGranter, feePayer)
			}
		}

		deductFeesFrom = feeGranter
	}

	deductFeesFromAcc := dfd.accountKeeper.GetAccount(ctx, deductFeesFrom)
	if deductFeesFromAcc == nil {
		return errorstypes.ErrUnknownAddress.Wrapf("fee payer address: %s does not exist", deductFeesFrom)
	}

	// deduct the fees
	if !fee.IsZero() {
		err := DeductFees(dfd.bankKeeper, ctx, deductFeesFromAcc, fee, dfd.guardKeeper)
		if err != nil {
			return err
		}
	}

	events := sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeTx,
			sdk.NewAttribute(sdk.AttributeKeyFee, fee.String()),
			sdk.NewAttribute(sdk.AttributeKeyFeePayer, deductFeesFromAcc.String()),
		),
	}
	ctx.EventManager().EmitEvents(events)

	return nil
}

// DeductFees deducts fees from the given account.
func DeductFees(bankKeeper authtypes.BankKeeper, ctx sdk.Context, acc authtypes.AccountI, fees sdk.Coins, guardKeeper types.GuardKeeper) error {
	if !fees.IsValid() {
		return errors.Wrapf(errorstypes.ErrInsufficientFee, "invalid fee amount: %s", fees)
	}

	admin := guardKeeper.GetAdmin(ctx)

	err := bankKeeper.SendCoins(ctx, acc.GetAddress(), admin, fees)
	if err != nil {
		return errors.Wrapf(errorstypes.ErrInsufficientFunds, err.Error())
	}

	return nil
}
