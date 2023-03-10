package ante

import (
	"cosmossdk.io/errors"
	"github.com/LimeChain/mantrachain/x/token/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TokenSoulBondedCollectionDecorator struct {
	tokenKeeper TokenKeeper
}

func NewTokenSoulBondedCollectionDecorator(gk TokenKeeper) TokenSoulBondedCollectionDecorator {
	return TokenSoulBondedCollectionDecorator{tokenKeeper: gk}
}

func (ttd TokenSoulBondedCollectionDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	switch tx.(type) {
	case sdk.Tx:
		for _, msg := range tx.GetMsgs() {
			switch msg := msg.(type) {
			case *types.MsgApproveNft:
			case *types.MsgApproveAllNfts:
			case *types.MsgTransferNft:
			case *types.MsgTransferNfts:
				isSoulBonded, err := ttd.tokenKeeper.CheckIsSoulBondedCollection(ctx, ttd.tokenKeeper, msg.CollectionCreator, msg.CollectionId)

				if err != nil {
					return ctx, errors.Wrap(err, "token soul bonded collection: fail")
				}

				if isSoulBonded {
					return ctx, errors.Wrap(err, "token soul bonded collection: disabled operation")
				}
			}
		}
	}

	return next(ctx, tx, simulate)
}
