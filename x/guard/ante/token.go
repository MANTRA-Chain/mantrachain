package ante

import (
	"cosmossdk.io/errors"
	tokentypes "github.com/LimeChain/mantrachain/x/token/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GuardTokenAuthzDecorator struct {
	guardKeeper GuardKeeper
}

func NewGuardTokenAuthzDecorator(gk GuardKeeper) GuardTokenAuthzDecorator {
	return GuardTokenAuthzDecorator{guardKeeper: gk}
}

func (gtd GuardTokenAuthzDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	switch tx.(type) {
	case sdk.Tx:
		for _, msg := range tx.GetMsgs() {
			switch msg := msg.(type) {
			case *tokentypes.MsgCreateNftCollection:
				if err := gtd.guardKeeper.CheckNewRestrictedNftsCollection(ctx, msg.Collection.RestrictedNfts, msg.GetCreator()); err != nil {
					return ctx, errors.Wrap(err, "guard token: fail")
				}
			case *tokentypes.MsgMintNft:
				if err := gtd.guardKeeper.CheckRestrictedNftsCollection(ctx, msg.CollectionCreator, msg.CollectionId, msg.GetCreator()); err != nil {
					return ctx, errors.Wrap(err, "guard token: fail")
				}
			case *tokentypes.MsgMintNfts:
				if err := gtd.guardKeeper.CheckRestrictedNftsCollection(ctx, msg.CollectionCreator, msg.CollectionId, msg.GetCreator()); err != nil {
					return ctx, errors.Wrap(err, "guard token: fail")
				}
			case *tokentypes.MsgBurnNft:
				if err := gtd.guardKeeper.CheckRestrictedNftsCollection(ctx, msg.CollectionCreator, msg.CollectionId, msg.GetCreator()); err != nil {
					return ctx, errors.Wrap(err, "guard token: fail")
				}
			case *tokentypes.MsgBurnNfts:
				if err := gtd.guardKeeper.CheckRestrictedNftsCollection(ctx, msg.CollectionCreator, msg.CollectionId, msg.GetCreator()); err != nil {
					return ctx, errors.Wrap(err, "guard token: fail")
				}
			case *tokentypes.MsgApproveNft:
				if err := gtd.guardKeeper.CheckRestrictedNftsCollection(ctx, msg.CollectionCreator, msg.CollectionId, msg.GetCreator()); err != nil {
					return ctx, errors.Wrap(err, "guard token: fail")
				}
			case *tokentypes.MsgApproveNfts:
				if err := gtd.guardKeeper.CheckRestrictedNftsCollection(ctx, msg.CollectionCreator, msg.CollectionId, msg.GetCreator()); err != nil {
					return ctx, errors.Wrap(err, "guard token: fail")
				}
			case *tokentypes.MsgTransferNft:
				if err := gtd.guardKeeper.CheckRestrictedNftsCollection(ctx, msg.CollectionCreator, msg.CollectionId, msg.GetCreator()); err != nil {
					return ctx, errors.Wrap(err, "guard token: fail")
				}
			case *tokentypes.MsgTransferNfts:
				if err := gtd.guardKeeper.CheckRestrictedNftsCollection(ctx, msg.CollectionCreator, msg.CollectionId, msg.GetCreator()); err != nil {
					return ctx, errors.Wrap(err, "guard token: fail")
				}
			}
		}
	}

	return next(ctx, tx, simulate)
}
