package ante

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	nft "github.com/cosmos/cosmos-sdk/x/nft"
)

type TokenTransferDecorator struct {
	tokenKeeper TokenKeeper
}

func NewTokenTransferDecorator(gk TokenKeeper) TokenTransferDecorator {
	return TokenTransferDecorator{tokenKeeper: gk}
}

func (ttd TokenTransferDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	switch tx.(type) {
	case sdk.Tx:
		for _, msg := range tx.GetMsgs() {
			switch msg := msg.(type) {
			case *nft.MsgSend:
				ok, err := ttd.tokenKeeper.CheckCanTransfer(ctx, ttd.tokenKeeper, msg.ClassId)

				if err != nil {
					return ctx, errors.Wrap(err, "send token: fail")
				}

				if !ok {
					return ctx, errors.Wrap(err, "send token: cannot transfer")
				}
			}
		}
	}

	return next(ctx, tx, simulate)
}
