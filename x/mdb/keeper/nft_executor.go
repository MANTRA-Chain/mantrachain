package keeper

import (
	"github.com/LimeChain/mantrachain/x/mdb/types"
	nfttypes "github.com/LimeChain/mantrachain/x/nft/types"
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

func (c *NftExecutor) SetClass(nftClass nfttypes.Class) (bool, error) {
	err := c.nftKeeper.SaveClass(c.ctx, nftClass)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *NftExecutor) SetDefaultClass(collIndex []byte) (bool, error) {
	err := c.nftKeeper.SaveClass(c.ctx, nfttypes.Class{
		Id:      string(collIndex),
		Name:    types.DefaultParams().NftCollectionDefaultName,
		Uri:     types.ModuleName,
		UriHash: types.DefaultParams().NftCollectionDefaultId,
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *NftExecutor) GetClass(ctx sdk.Context, classId string) (nfttypes.Class, bool) {
	return c.nftKeeper.GetClass(c.ctx, classId)
}

func (c *NftExecutor) MintNftBatch(nfts []nfttypes.NFT, receiver sdk.AccAddress) (bool, error) {
	for _, nft := range nfts {
		err := c.nftKeeper.Mint(c.ctx, nft, receiver)

		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func (c *NftExecutor) GetNft(ctx sdk.Context, classId string, nftId string) (nfttypes.NFT, bool) {
	return c.nftKeeper.GetNFT(c.ctx, classId, nftId)
}
