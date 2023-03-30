package ante

import (
	"cosmossdk.io/errors"
	coinfactorytypes "github.com/MANTRA-Finance/mantrachain/x/coinfactory/types"
	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
	liquiditytypes "github.com/MANTRA-Finance/mantrachain/x/liquidity/types"
	lpfarmtypes "github.com/MANTRA-Finance/mantrachain/x/lpfarm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GuardAdminAuthzDecorator struct {
	guardKeeper GuardKeeper
}

func NewGuardAdminAuthzDecorator(gk GuardKeeper) GuardAdminAuthzDecorator {
	return GuardAdminAuthzDecorator{guardKeeper: gk}
}

func (gtd GuardAdminAuthzDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	switch tx.(type) {
	case sdk.Tx:
		for _, msg := range tx.GetMsgs() {
			switch msg := msg.(type) {
			// Guard module
			case *types.MsgUpdateAccountPrivileges:
				if err := gtd.guardKeeper.CheckIsAdmin(ctx, msg.GetCreator()); err != nil {
					return ctx, errors.Wrap(err, "unauthorized")
				}
			case *types.MsgUpdateAccountPrivilegesBatch:
				if err := gtd.guardKeeper.CheckIsAdmin(ctx, msg.GetCreator()); err != nil {
					return ctx, errors.Wrap(err, "unauthorized")
				}
			case *types.MsgUpdateAccountPrivilegesGroupedBatch:
				if err := gtd.guardKeeper.CheckIsAdmin(ctx, msg.GetCreator()); err != nil {
					return ctx, errors.Wrap(err, "unauthorized")
				}
			case *types.MsgUpdateGuardTransferCoins:
				if err := gtd.guardKeeper.CheckIsAdmin(ctx, msg.GetCreator()); err != nil {
					return ctx, errors.Wrap(err, "unauthorized")
				}
			case *types.MsgUpdateLocked:
				if err := gtd.guardKeeper.CheckIsAdmin(ctx, msg.GetCreator()); err != nil {
					return ctx, errors.Wrap(err, "unauthorized")
				}
			case *types.MsgUpdateRequiredPrivileges:
				if err := gtd.guardKeeper.CheckIsAdmin(ctx, msg.GetCreator()); err != nil {
					return ctx, errors.Wrap(err, "unauthorized")
				}
			case *types.MsgUpdateRequiredPrivilegesBatch:
				if err := gtd.guardKeeper.CheckIsAdmin(ctx, msg.GetCreator()); err != nil {
					return ctx, errors.Wrap(err, "unauthorized")
				}
			case *types.MsgUpdateRequiredPrivilegesGroupedBatch:
				if err := gtd.guardKeeper.CheckIsAdmin(ctx, msg.GetCreator()); err != nil {
					return ctx, errors.Wrap(err, "unauthorized")
				}
			case *types.MsgUpdateAuthzGenericGrantRevokeBatch:
				if err := gtd.guardKeeper.CheckIsAdmin(ctx, msg.GetCreator()); err != nil {
					return ctx, errors.Wrap(err, "unauthorized")
				}
			// Coin factory module
			case *coinfactorytypes.MsgCreateDenom:
				if err := gtd.guardKeeper.CheckIsAdmin(ctx, msg.Sender); err != nil {
					return ctx, errors.Wrap(err, "unauthorized")
				}
			// Liquidity module
			case *liquiditytypes.MsgCreatePair:
				if err := gtd.guardKeeper.CheckIsAdmin(ctx, msg.Creator); err != nil {
					return ctx, errors.Wrap(err, "unauthorized")
				}
			case *liquiditytypes.MsgCreatePool:
				if err := gtd.guardKeeper.CheckIsAdmin(ctx, msg.Creator); err != nil {
					return ctx, errors.Wrap(err, "unauthorized")
				}
			case *liquiditytypes.MsgCreateRangedPool:
				if err := gtd.guardKeeper.CheckIsAdmin(ctx, msg.Creator); err != nil {
					return ctx, errors.Wrap(err, "unauthorized")
				}
			// case *liquiditytypes.MsgDeposit:
			// 	if err := gtd.guardKeeper.CheckIsAdmin(ctx, msg.Depositor); err != nil {
			// 		return ctx, errors.Wrap(err, "unauthorized")
			// 	}
			// case *liquiditytypes.MsgWithdraw:
			// 	if err := gtd.guardKeeper.CheckIsAdmin(ctx, msg.Withdrawer); err != nil {
			// 		return ctx, errors.Wrap(err, "unauthorized")
			// 	}
			// case *liquiditytypes.MsgCancelAllOrders:
			// 	if err := gtd.guardKeeper.CheckIsAdmin(ctx, msg.Orderer); err != nil {
			// 		return ctx, errors.Wrap(err, "unauthorized")
			// 	}
			// Lpfarm module
			case *lpfarmtypes.MsgCreatePrivatePlan:
				if err := gtd.guardKeeper.CheckIsAdmin(ctx, msg.Creator); err != nil {
					return ctx, errors.Wrap(err, "unauthorized")
				}
			}
		}
	}

	return next(ctx, tx, simulate)
}
