package keeper

import (
	"github.com/LimeChain/mantrachain/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type NftExecutor struct {
	ctx       sdk.Context
	nftKeeper types.NFTKeeper
}

func NewNftExecutor(ctx sdk.Context, nftKeeper types.NFTKeeper) *NftExecutor {
	return &NftExecutor{
		ctx:       ctx,
		nftKeeper: nftKeeper,
	}
}

func (c *NftExecutor) GetNftOwner(classId string, nftId string) sdk.AccAddress {
	return c.nftKeeper.GetOwner(c.ctx, classId, nftId)
}
