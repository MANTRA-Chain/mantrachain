package ante

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// EVMMsgCheckDecorator is a decorator that checks if the transaction contains
// exactly one EVM message. If it contains more than one, it returns an error.
type EVMMsgCheckDecorator struct{}

func (emcd EVMMsgCheckDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	// only allow for a single EVM message in the transaction
	if len(tx.GetMsgs()) > 1 {
		return ctx, errorsmod.Wrapf(
			errortypes.ErrInvalidRequest,
			"expected only one EVM message, got %d",
			len(tx.GetMsgs()),
		)
	}

	return next(ctx, tx, simulate)
}
