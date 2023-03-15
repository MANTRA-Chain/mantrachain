package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GuardAuthzDecorator struct {
	guardKeeper GuardKeeper
}

func NewGuardAuthzDecorator(gk GuardKeeper) GuardAuthzDecorator {
	return GuardAuthzDecorator{guardKeeper: gk}
}

func (gtd GuardAuthzDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	// switch tx.(type) {
	// case sdk.Tx:
	// 	for _, msg := range tx.GetMsgs() {
	// 		switch msg := msg.(type) {
	// 		case *banktypes.MsgSend:
	// 			ok, err := gtd.guardKeeper.CheckHasPerm(ctx, msg)

	// 			if err != nil {
	// 				return ctx, errors.Wrap(err, "fail")
	// 			}

	// 			if !ok {
	// 				return ctx, errors.Wrap(err, "cannot execute")
	// 			}
	// 		}
	// 	}
	// }

	return next(ctx, tx, simulate)
}
