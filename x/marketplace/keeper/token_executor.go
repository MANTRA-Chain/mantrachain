package keeper

import (
	"github.com/LimeChain/mantrachain/x/marketplace/types"
	tokentypes "github.com/LimeChain/mantrachain/x/token/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TokenExecutor struct {
	ctx         sdk.Context
	tokenKeeper types.TokenKeeper
}

func NewTokenExecutor(ctx sdk.Context, tokenKeeper types.TokenKeeper) *TokenExecutor {
	return &TokenExecutor{
		ctx:         ctx,
		tokenKeeper: tokenKeeper,
	}
}

func (c *TokenExecutor) GetNftCollection(creator sdk.AccAddress, id string) (tokentypes.NftCollection, bool) {
	index := tokentypes.GetNftCollectionIndex(creator, id)
	return c.tokenKeeper.GetNftCollection(c.ctx, creator, index)
}

func (c *TokenExecutor) HasNftCollection(creator sdk.AccAddress, id string) bool {
	index := tokentypes.GetNftCollectionIndex(creator, id)
	return c.tokenKeeper.HasNftCollection(c.ctx, creator, index)
}

func (c *TokenExecutor) GetNft(collectionCreator sdk.AccAddress, collectionId string, id string) (tokentypes.Nft, bool) {
	collectionIndex := tokentypes.GetNftCollectionIndex(collectionCreator, collectionId)
	index := tokentypes.GetNftIndex(collectionIndex, id)
	return c.tokenKeeper.GetNft(c.ctx, collectionIndex, index)
}
