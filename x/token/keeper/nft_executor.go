package keeper

import (
	"github.com/LimeChain/mantrachain/x/token/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	nft "github.com/cosmos/cosmos-sdk/x/nft"
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

func (c *NftExecutor) SetClass(nftClass nft.Class) error {
	return c.nftKeeper.SaveClass(c.ctx, nftClass)
}

func (c *NftExecutor) SetDefaultClass(collIndex []byte) error {
	return c.nftKeeper.SaveClass(c.ctx, nft.Class{
		Id:      string(collIndex),
		Name:    types.DefaultParams().NftCollectionDefaultName,
		Uri:     types.ModuleName,
		UriHash: types.DefaultParams().NftCollectionDefaultId,
	})
}

func (c *NftExecutor) GetClass(classId string) (nft.Class, bool) {
	return c.nftKeeper.GetClass(c.ctx, classId)
}

func (c *NftExecutor) GetClassSupply(classId string) uint64 {
	return c.nftKeeper.GetTotalSupply(c.ctx, classId)
}

// TODO: fix ASAP
func (c *NftExecutor) GetClasses(classesIds []string) []nft.Class {
	return []nft.Class{}
	// return c.nftKeeper.GetClassesByIds(c.ctx, classesIds)
}

func (c *NftExecutor) MintNft(nft nft.NFT, receiver sdk.AccAddress) error {
	return c.nftKeeper.Mint(c.ctx, nft, receiver)
}

func (c *NftExecutor) MintNftBatch(nfts []nft.NFT, receiver sdk.AccAddress) error {
	return c.nftKeeper.BatchMint(c.ctx, nfts, receiver)
}

func (c *NftExecutor) BurnNftBatch(classId string, nftsIds []string) error {
	return c.nftKeeper.BatchBurn(c.ctx, classId, nftsIds)
}

func (c *NftExecutor) BurnNft(classId string, nftId string) error {
	return c.nftKeeper.Burn(c.ctx, classId, nftId)
}

func (c *NftExecutor) GetNft(classId string, nftId string) (nft.NFT, bool) {
	return c.nftKeeper.GetNFT(c.ctx, classId, nftId)
}

// TODO: fix ASAP
func (c *NftExecutor) GetNfts(classId string, nftsIds []string) []nft.NFT {
	return []nft.NFT{}
	// return c.nftKeeper.GetNFTsByIds(c.ctx, classId, nftsIds)
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
	return c.nftKeeper.BatchTransfer(c.ctx, classId, nftsIds, receiver)
}

func (c *NftExecutor) GetNftsOfClassByOwner(classId string, owner sdk.AccAddress) []nft.NFT {
	return c.nftKeeper.GetNFTsOfClassByOwner(c.ctx, classId, owner)
}
