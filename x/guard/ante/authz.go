package ante

import (
	"cosmossdk.io/errors"
	coinfactorytypes "github.com/LimeChain/mantrachain/x/coinfactory/types"
	"github.com/LimeChain/mantrachain/x/guard/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GuardAuthzDecorator struct {
	guardKeeper GuardKeeper
}

func NewGuardAuthzDecorator(gk GuardKeeper) GuardAuthzDecorator {
	return GuardAuthzDecorator{guardKeeper: gk}
}

func (gtd GuardAuthzDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	switch tx.(type) {
	case sdk.Tx:
		for _, msg := range tx.GetMsgs() {
			switch msg := msg.(type) {
			case *types.MsgUpdateAccountPrivileges:
			case *types.MsgUpdateAccountPrivilegesBatch:
			case *types.MsgUpdateAccountPrivilegesGroupedBatch:
			case *types.MsgUpdateGuardTransferCoins:
			case *types.MsgUpdateLocked:
			case *types.MsgUpdateRequiredPrivileges:
			case *types.MsgUpdateRequiredPrivilegesBatch:
			case *types.MsgUpdateRequiredPrivilegesGroupedBatch:
				err := gtd.guardKeeper.CheckHasPerm(ctx, msg.GetCreator())

				if err != nil {
					return ctx, errors.Wrap(err, "fail")
				}
			case *coinfactorytypes.MsgCreateDenom:
			case *coinfactorytypes.MsgMint:
			case *coinfactorytypes.MsgBurn:
			case *coinfactorytypes.MsgChangeAdmin:
			case *coinfactorytypes.MsgForceTransfer:
			case *coinfactorytypes.MsgSetDenomMetadata:
				err := gtd.guardKeeper.CheckHasPerm(ctx, msg.Sender)

				if err != nil {
					return ctx, errors.Wrap(err, "fail")
				}
			}
			// TODO: add rest of the modules relevant messages
		}
	}

	return next(ctx, tx, simulate)
}
