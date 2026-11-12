package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/MANTRA-Chain/mantrachain/v8/x/sanction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

// EVMBlacklistCheckDecorator rejects an EVM transaction whose sender is
// blacklisted, or that carries an EIP-7702 authorization signed by a
// blacklisted account. Without the second check, a relayer could submit a
// sanctioned account's authorization and run code as that account.
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

		if err := ebcd.checkAddress(ctx, ethcommon.BytesToAddress(ethMsg.GetFrom()), "sender"); err != nil {
			return ctx, err
		}

		for _, auth := range ethMsg.AsTransaction().SetCodeAuthorizations() {
			authority, err := auth.Authority()
			if err != nil {
				// The EVM skips an authorization whose signer can't be
				// recovered, so it can't install code for anyone.
				continue
			}
			if err := ebcd.checkAddress(ctx, authority, "authorization signer"); err != nil {
				return ctx, err
			}
		}
	}

	return next(ctx, tx, simulate)
}

func (ebcd EVMBlacklistCheckDecorator) checkAddress(ctx sdk.Context, addr ethcommon.Address, what string) error {
	has, err := ebcd.sanctionKeeper.BlacklistAccounts.Has(ctx, sdk.AccAddress(addr.Bytes()).String())
	if err != nil {
		return err
	}
	if has {
		return errorsmod.Wrapf(types.ErrAccountBlacklisted, "%s %s is blacklisted", what, addr)
	}
	return nil
}
