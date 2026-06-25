package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/MANTRA-Chain/mantrachain/v8/x/sanction/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authz "github.com/cosmos/cosmos-sdk/x/authz"
)

type BlacklistCheckDecorator struct {
	sanctionKeeper Keeper
	cdc            codec.Codec
}

func NewBlacklistCheckDecorator(sk Keeper, cdc codec.Codec) BlacklistCheckDecorator {
	return BlacklistCheckDecorator{
		sanctionKeeper: sk,
		cdc:            cdc,
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

	// Check fee granter
	if feeTx, ok := tx.(sdk.FeeTx); ok {
		if granter := feeTx.FeeGranter(); granter != nil {
			granterAddr := sdk.AccAddress(granter)
			has, err := d.sanctionKeeper.BlacklistAccounts.Has(ctx, granterAddr.String())
			if err != nil {
				return ctx, err
			}
			if has {
				return ctx, errorsmod.Wrapf(types.ErrAccountBlacklisted, "fee granter %s is blacklisted", granterAddr)
			}
		}
	}

	// Check authz granters (inner message signers of MsgExec).
	// Only a single MsgExec per tx is allowed; multiple could obscure blacklisted granters.
	// Nested MsgExec is also rejected outright — it cannot be safely inspected
	// with a flat check and could be used to hide a blacklisted granter.
	var execCount int
	for _, msg := range tx.GetMsgs() {
		if _, ok := msg.(*authz.MsgExec); ok {
			execCount++
			if execCount > 1 {
				return ctx, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "only a single MsgExec is allowed per transaction")
			}
		}
	}
	for _, msg := range tx.GetMsgs() {
		execMsg, ok := msg.(*authz.MsgExec)
		if !ok {
			continue
		}
		innerMsgs, err := execMsg.GetMessages()
		if err != nil {
			return ctx, err
		}
		for _, innerMsg := range innerMsgs {
			if _, ok := innerMsg.(*authz.MsgExec); ok {
				return ctx, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "nested MsgExec is not allowed")
			}
			innerSigners, _, err := d.cdc.GetMsgV1Signers(innerMsg)
			if err != nil {
				return ctx, err
			}
			if len(innerSigners) == 0 {
				continue
			}
			granterAddr := sdk.AccAddress(innerSigners[0])
			has, err := d.sanctionKeeper.BlacklistAccounts.Has(ctx, granterAddr.String())
			if err != nil {
				return ctx, err
			}
			if has {
				return ctx, errorsmod.Wrapf(types.ErrAccountBlacklisted, "authz granter %s is blacklisted", granterAddr)
			}
		}
	}

	return next(ctx, tx, simulate)
}
