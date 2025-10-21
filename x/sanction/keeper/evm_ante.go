package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/MANTRA-Chain/mantrachain/v7/x/sanction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

// EVMBlacklistCheckDecorator is a decorator that checks if EVM transaction senders
// are blacklisted. If a sender is found in the blacklist, it returns an error.
type EVMBlacklistCheckDecorator struct {
	sanctionKeeper Keeper
}

func NewEVMBlacklistCheckDecorator(sk Keeper) EVMBlacklistCheckDecorator {
	return EVMBlacklistCheckDecorator{
		sanctionKeeper: sk,
	}
}

func (ebcd EVMBlacklistCheckDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	for _, msg := range tx.GetMsgs() {
		ethMsg, _, err := evmtypes.UnpackEthMsg(msg)
		if err != nil {
			return ctx, err
		}

		if has, err := ebcd.sanctionKeeper.BlacklistAccounts.Has(ctx, ethMsg.GetFrom().String()); err != nil {
			return ctx, err
		} else if has {
			return ctx, errorsmod.Wrapf(types.ErrAccountBlacklisted, "%s is blacklisted", ethcommon.BytesToAddress(ethMsg.GetFrom()))
		}
	}

	return next(ctx, tx, simulate)
}
