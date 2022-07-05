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

func (c *NftExecutor) GetClass(classId string) (nfttypes.Class, bool) {
	return c.nftKeeper.GetClass(c.ctx, classId)
}

func (c *NftExecutor) GetClasses(classesIds []string) []nfttypes.Class {
	// TODO: use async iterator
	var classes []nfttypes.Class
	for _, classId := range classesIds {
		class, _ := c.nftKeeper.GetClass(c.ctx, classId)

		classes = append(classes, class)
	}
	return classes
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

func (c *NftExecutor) BurnNftBatch(classId string, nftsIds []string) (bool, error) {
	// TODO: use async iterator
	for _, id := range nftsIds {
		err := c.nftKeeper.Burn(c.ctx, classId, id)

		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func (c *NftExecutor) GetNft(classId string, nftId string) (nfttypes.NFT, bool) {
	return c.nftKeeper.GetNFT(c.ctx, classId, nftId)
}

func (c *NftExecutor) GetNfts(classId string, nftsIds []string) []nfttypes.NFT {
	// TODO: use async iterator
	var nfts []nfttypes.NFT
	for _, nftId := range nftsIds {
		nft, _ := c.nftKeeper.GetNFT(c.ctx, classId, nftId)

		nfts = append(nfts, nft)
	}
	return nfts
}

func (c *NftExecutor) GetNftOwner(classId string, nftId string) sdk.AccAddress {
	return c.nftKeeper.GetOwner(c.ctx, classId, nftId)
}
