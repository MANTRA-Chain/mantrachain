package ante

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	nft "github.com/cosmos/cosmos-sdk/x/nft"
)

type TokenTransferNftDecorator struct {
	tokenKeeper TokenKeeper
}

func NewTokenTransferNftDecorator(gk TokenKeeper) TokenTransferNftDecorator {
	return TokenTransferNftDecorator{tokenKeeper: gk}
}

func (ttd TokenTransferNftDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	switch tx.(type) {
	case sdk.Tx:
		for _, msg := range tx.GetMsgs() {
			switch msg := msg.(type) {
			case *nft.MsgSend:
				ok, err := ttd.tokenKeeper.CheckCanTransferNft(ctx, ttd.tokenKeeper, msg.ClassId)

				if err != nil {
					return ctx, errors.Wrap(err, "token send nft: fail")
				}

				if !ok {
					return ctx, errors.Wrap(err, "token send nft: cannot transfer")
				}
			}
		}
	}

	return next(ctx, tx, simulate)
}
