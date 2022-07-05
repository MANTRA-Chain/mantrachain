package keeper

import (
	"context"

	"github.com/LimeChain/mantrachain/x/mdb/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) CollectionNfts(goCtx context.Context, req *types.QueryGetCollectionNftsRequest) (*types.QueryGetCollectionNftsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	collectionCreator, err := sdk.AccAddressFromBech32(req.CollectionCreator)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.CollectionId == "" {
		return nil, status.Error(codes.InvalidArgument, "empty collection id")
	}

	collIndex := types.GetNftCollectionIndex(collectionCreator, req.CollectionId)

	if !k.HasNftCollection(ctx, collectionCreator, collIndex) {
		return nil, status.Error(codes.InvalidArgument, "collection not exists")
	}

	nftsMeta := k.GetAllNft(
		ctx,
		collIndex,
	)

	var nftsIndexes []string

	for _, nftMeta := range nftsMeta {
		nftsIndexes = append(nftsIndexes, string(nftMeta.Index))
	}

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)
	nfts := nftExecutor.GetNfts(string(collIndex), nftsIndexes)

	var nftsRes []*types.QueryGetNftResponse

	for i, nft := range nfts {
		meta := nftsMeta[i]

		nftsRes = append(nftsRes, &types.QueryGetNftResponse{
			Id:           nft.UriHash,
			Title:        meta.Title,
			Did:          meta.Did,
			Description:  meta.Description,
			Images:       meta.Images,
			Url:          meta.Url,
			Links:        meta.Links,
			Attributes:   meta.Attributes,
			Creator:      meta.Creator.String(),
			Owner:        nftExecutor.GetNftOwner(string(collIndex), string(meta.Index)).String(),
			Resellable:   meta.Resellable,
			Data:         nft.Data,
			CollectionId: meta.CollectionId,
		})
	}

	return &types.QueryGetCollectionNftsResponse{
		Nfts: nftsRes,
	}, nil
}
