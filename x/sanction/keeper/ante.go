package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/MANTRA-Chain/mantrachain/v7/x/sanction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
)

type BlacklistCheckDecorator struct {
	sanctionKeeper Keeper
}

func NewBlacklistCheckDecorator(sk Keeper) BlacklistCheckDecorator {
	return BlacklistCheckDecorator{
		sanctionKeeper: sk,
	}
}

// AnteHandle checks if the tx signer is blacklisted
func (d BlacklistCheckDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	sigTx, ok := tx.(authsigning.Tx)
	if !ok {
		return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type")
	}
	signers, err := sigTx.GetSigners()
	if err != nil {
		return ctx, err
	}
	for _, signer := range signers {
		signerAcc := sdk.AccAddress(signer)
		has, err := d.sanctionKeeper.BlacklistAccounts.Has(ctx, signerAcc.String())
		if err != nil {
			return ctx, err
		}
		if has {
			return ctx, errorsmod.Wrapf(types.ErrAccountBlacklisted, "%s is blacklisted", signerAcc.String())
		}
	}
	return next(ctx, tx, simulate)
}
