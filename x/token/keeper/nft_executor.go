package keeper

import (
	"github.com/MANTRA-Finance/mantrachain/x/token/types"
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

func (c *NftExecutor) MintNft(nft nft.NFT, receiver sdk.AccAddress) error {
	return c.nftKeeper.Mint(c.ctx, nft, receiver)
}

func (c *NftExecutor) BatchMintNft(nfts []nft.NFT, receiver sdk.AccAddress) error {
	return c.nftKeeper.BatchMint(c.ctx, nfts, receiver)
}

func (c *NftExecutor) BatchBurnNft(classId string, nftsIds []string) error {
	return c.nftKeeper.BatchBurn(c.ctx, classId, nftsIds)
}

func (c *NftExecutor) BurnNft(classId string, nftId string) error {
	return c.nftKeeper.Burn(c.ctx, classId, nftId)
}

func (c *NftExecutor) GetNftOwner(classId string, nftId string) sdk.AccAddress {
	return c.nftKeeper.GetOwner(c.ctx, classId, nftId)
}

func (c *NftExecutor) TransferNft(classId string, nftId string, receiver sdk.AccAddress) error {
	return c.nftKeeper.Transfer(c.ctx, classId, nftId, receiver)
}

func (c *NftExecutor) BatchTransferNft(classId string, nftsIds []string, receiver sdk.AccAddress) error {
	return c.nftKeeper.BatchTransfer(c.ctx, classId, nftsIds, receiver)
}
