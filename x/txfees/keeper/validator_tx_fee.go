package keeper

import (
	"cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/MANTRA-Finance/mantrachain/x/txfees/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorstypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// checkTxFeeWithValidatorMinGasPrices implements the default fee logic, where the minimum price per
// unit of gas is fixed and set by each validator, can the tx priority is computed from the gas price.
func checkTxFeeWithValidatorMinGasPrices(ctx sdk.Context, tx sdk.Tx, baseDenom string, txfeesKeeper types.TxfeesKeeper, liquidityKeeper types.LiquidityKeeper) (sdk.Coins, error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return nil, errors.Wrap(errorstypes.ErrTxDecode, "Tx must be a FeeTx")
	}

	feeCoins := feeTx.GetFee()

	if len(feeCoins) > 1 {
		return nil, errors.Wrap(types.ErrTooManyFeeCoins, "Only accepts fees in one denom")
	}

	var pairId uint64
	var feeDenomNotBaseDenom bool

	if len(feeCoins) == 1 {
		feeDenom := feeCoins.GetDenomByIndex(0)
		feeDenomNotBaseDenom = feeDenom != baseDenom

		if feeDenomNotBaseDenom {
			feeToken, foundFeeToken := txfeesKeeper.GetFeeToken(ctx, feeDenom)
			if !foundFeeToken {
				return nil, errors.Wrap(types.ErrInvalidFeeDenom, "Invalid fee denom")
			}

			pairId = feeToken.PairId
		}
	}

	gas := feeTx.GetGas()

	// Ensure that the provided fees meet a minimum threshold for the validator,
	// if this is a CheckTx. This is only for local mempool purposes, and thus
	// is only ran on check tx.
	if ctx.IsCheckTx() {
		minGasPrices := ctx.MinGasPrices()

		if len(minGasPrices) > 1 {
			return nil, errors.Wrap(types.ErrTooManyGasPricesCoins, "Only accepts min gas prices in one denom")
		}

		if !minGasPrices.IsZero() {
			requiredFees := make(sdk.Coins, len(minGasPrices))

			// Determine the required fees by multiplying each required minimum gas
			// price by the gas limit, where fee = ceil(minGasPrice * gasLimit).
			glDec := sdkmath.LegacyNewDec(int64(gas))
			for i, gp := range minGasPrices {
				fee := gp.Amount.Mul(glDec).Ceil().RoundInt()

				if feeDenomNotBaseDenom {
					offerCoin, _, err := liquidityKeeper.GetSwapAmount(ctx, pairId, sdk.NewCoin(gp.Denom, fee))
					if err != nil {
						return nil, err
					}

					if offerCoin.IsZero() {
						return nil, errors.Wrapf(types.ErrZeroFee, "zero fees; required fees: %s", offerCoin)
					}

					requiredFees[i] = offerCoin
				} else {
					requiredFees[i] = sdk.NewCoin(gp.Denom, fee)
				}
			}

			if !feeCoins.IsAnyGTE(requiredFees) {
				return nil, errors.Wrapf(errorstypes.ErrInsufficientFee, "insufficient fees; got: %s required: %s", feeCoins, requiredFees)
			}
		}
	}

	return feeCoins, nil
}
