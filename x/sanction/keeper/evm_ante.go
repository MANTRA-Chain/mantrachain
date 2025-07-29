package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/MANTRA-Chain/mantrachain/v5/x/sanction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

// EVMMsgCheckDecorator is a decorator that checks if the transaction contains
// exactly one EVM message. If it contains more than one, it returns an error.
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
