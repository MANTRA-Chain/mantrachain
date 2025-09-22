package ante

import (
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
)

// MultiChainIDDecorator is a custom decorator that allows an additional chain ID for signature verification.
// It wraps the standard SigVerificationDecorator.
type MultiChainIDDecorator struct {
	svd ante.SigVerificationDecorator
}

// NewMultiChainIDDecorator creates a new MultiChainIDDecorator.
func NewMultiChainIDDecorator(svd ante.SigVerificationDecorator) MultiChainIDDecorator {
	return MultiChainIDDecorator{
		svd: svd,
	}
}

// AnteHandle performs signature verification. If the signature verification fails with the context's
// chain ID, it retries with a derived EVM-compatible chain ID (e.g. "mantra-1" -> "mantra-evm-1").
func (mcd MultiChainIDDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	// A no-op next handler to isolate the SigVerificationDecorator's logic from the rest of the ante chain.
	noOpNext := func(c sdk.Context, t sdk.Tx, s bool) (sdk.Context, error) { return c, nil }

	// Try to verify with the default chain ID from the context.
	_, err1 := mcd.svd.AnteHandle(ctx, tx, simulate, noOpNext)
	if err1 == nil {
		// If it succeeded, then we call the real `next` handler.
		return next(ctx, tx, simulate)
	}

	// If the error is not an unauthorized error, return the error.
	if !errorsmod.IsOf(err1, errortypes.ErrUnauthorized) {
		return ctx, err1
	}

	// If signature verification failed with ErrUnauthorized, try again with the derived EVM chain ID.
	chainID := ctx.ChainID()
	lastHyphen := strings.LastIndex(chainID, "-")
	if lastHyphen == -1 {
		// The chain ID is not in the expected <string>-<number> format.
		// Return the original unauthorized error.
		return ctx, err
	}

	name := chainID[:lastHyphen]
	version := chainID[lastHyphen+1:]
	evmCosmosChainId := name + "-evm-" + version
	extraCtx := ctx.WithChainID(evmCosmosChainId)

	// verify with the "-evm-" chain ID.
	_, err = mcd.svd.AnteHandle(extraCtx, tx, simulate, noOpNext)
	if err != nil {
		return ctx, fmt.Errorf("%s, %s", err1, err)
	}
	// If it succeeded, then we call the real `next` handler with the default chain ID.
	return next(ctx, tx, simulate)
}
