package ante

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type GuardTransferCoinsDecorator struct {
	guardKeeper GuardKeeper
	nftKeeper   NFTKeeper
}

func NewGuardTransferCoinsDecorator(gk GuardKeeper, nk NFTKeeper) GuardTransferCoinsDecorator {
	return GuardTransferCoinsDecorator{guardKeeper: gk, nftKeeper: nk}
}

func (gtd GuardTransferCoinsDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	switch tx.(type) {
	case sdk.Tx:
		for _, msg := range tx.GetMsgs() {
			switch msg := msg.(type) {
			case *banktypes.MsgSend:
				from, err := sdk.AccAddressFromBech32(msg.FromAddress)

				if err != nil {
					return ctx, sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", err)
				}

				to, err := sdk.AccAddressFromBech32(msg.ToAddress)

				if err != nil {
					return ctx, sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", err)
				}

				ok, err := gtd.guardKeeper.CheckCanTransferCoins(ctx, gtd.nftKeeper, []sdk.AccAddress{from, to}, msg.Amount)

				if err != nil {
					return ctx, errors.Wrap(err, "guard send coins: fail")
				}

				if !ok {
					return ctx, errors.Wrap(err, "guard send coins: cannot transfer")
				}
			}
		}
	}

	return next(ctx, tx, simulate)
}
