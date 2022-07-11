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

func (c *NftExecutor) SetClass(nftClass nfttypes.Class) error {
	return c.nftKeeper.SaveClass(c.ctx, nftClass)
}

func (c *NftExecutor) SetDefaultClass(collIndex []byte) error {
	return c.nftKeeper.SaveClass(c.ctx, nfttypes.Class{
		Id:      string(collIndex),
		Name:    types.DefaultParams().NftCollectionDefaultName,
		Uri:     types.ModuleName,
		UriHash: types.DefaultParams().NftCollectionDefaultId,
	})
}

func (c *NftExecutor) GetClass(classId string) (nfttypes.Class, bool) {
	return c.nftKeeper.GetClass(c.ctx, classId)
}

func (c *NftExecutor) GetClasses(classesIds []string) []nfttypes.Class {
	return c.nftKeeper.GetClassesByIds(c.ctx, classesIds)
}

func (c *NftExecutor) MintNft(nft nfttypes.NFT, receiver sdk.AccAddress) error {
	return c.nftKeeper.Mint(c.ctx, nft, receiver)
}

func (c *NftExecutor) MintNftBatch(nfts []nfttypes.NFT, receiver sdk.AccAddress) error {
	return c.nftKeeper.MintBatch(c.ctx, nfts, receiver)
}

func (c *NftExecutor) BurnNftBatch(classId string, nftsIds []string) error {
	return c.nftKeeper.BurnBatch(c.ctx, classId, nftsIds)
}

func (c *NftExecutor) BurnNft(classId string, nftId string) error {
	return c.nftKeeper.Burn(c.ctx, classId, nftId)
}

func (c *NftExecutor) GetNft(classId string, nftId string) (nfttypes.NFT, bool) {
	return c.nftKeeper.GetNFT(c.ctx, classId, nftId)
}

func (c *NftExecutor) GetNfts(classId string, nftsIds []string) []nfttypes.NFT {
	return c.nftKeeper.GetNFTsByIds(c.ctx, classId, nftsIds)
}

func (c *NftExecutor) GetNftOwner(classId string, nftId string) sdk.AccAddress {
	return c.nftKeeper.GetOwner(c.ctx, classId, nftId)
}

func (c *NftExecutor) GetNftBalance(classId string, owner sdk.AccAddress) uint64 {
	return c.nftKeeper.GetBalance(c.ctx, classId, owner)
}

func (c *NftExecutor) TransferNft(classId string, nftId string, receiver sdk.AccAddress) error {
	return c.nftKeeper.Transfer(c.ctx, classId, nftId, receiver)
}

func (c *NftExecutor) TransferNftBatch(classId string, nftsIds []string, receiver sdk.AccAddress) error {
	return c.nftKeeper.TransferBatch(c.ctx, classId, nftsIds, receiver)
}
