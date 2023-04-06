package ante

import (
	"github.com/MANTRA-Finance/mantrachain/x/token/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

type TokenSoulBondedNftsCollectionDecorator struct {
	tokenKeeper TokenKeeper
}

func NewTokenSoulBondedNftsCollectionDecorator(tk TokenKeeper) TokenSoulBondedNftsCollectionDecorator {
	return TokenSoulBondedNftsCollectionDecorator{tokenKeeper: tk}
}

func (ttd TokenSoulBondedNftsCollectionDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	switch tx.(type) {
	case sdk.Tx:
		for _, msg := range tx.GetMsgs() {
			switch msg := msg.(type) {
			case *types.MsgApproveNft:
				if err := ttd.tokenKeeper.CheckSoulBondedNftsCollection(ctx, msg.CollectionCreator, msg.CollectionId); err != nil {
					return ctx, errors.Wrap(err, "token soul bonded collection: fail")
				}
			case *types.MsgTransferNft:
				if err := ttd.tokenKeeper.CheckSoulBondedNftsCollection(ctx, msg.CollectionCreator, msg.CollectionId); err != nil {
					return ctx, errors.Wrap(err, "token soul bonded collection: fail")
				}
			case *types.MsgTransferNfts:
				if err := ttd.tokenKeeper.CheckSoulBondedNftsCollection(ctx, msg.CollectionCreator, msg.CollectionId); err != nil {
					return ctx, errors.Wrap(err, "token soul bonded collection: fail")
				}
			}
		}
	}

	return next(ctx, tx, simulate)
}
