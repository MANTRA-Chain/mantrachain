package keeper

import (
	"context"

	"github.com/LimeChain/mantrachain/x/mdb/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) NftCollections(c context.Context, req *types.QueryGetNftCollectionsRequest) (*types.QueryGetNftCollectionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	creator, err := sdk.AccAddressFromBech32(req.Creator)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	colls := k.GetAllNftCollection(
		ctx,
		sdk.AccAddress(creator),
	)

	var collsIndexes []string

	for _, coll := range colls {
		collsIndexes = append(collsIndexes, string(coll.Index))
	}

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)
	nftColls := nftExecutor.GetClasses(collsIndexes)

	var nftCollections []*types.QueryGetNftCollectionResponse

	for i, nftColl := range nftColls {
		meta := colls[i]
		nftCollections = append(nftCollections, &types.QueryGetNftCollectionResponse{
			Id:          nftColl.UriHash,
			Name:        nftColl.Name,
			Symbol:      nftColl.Symbol,
			Description: nftColl.Description,
			Did:         meta.Did,
			Images:      meta.Images,
			Url:         meta.Url,
			Links:       meta.Links,
			Category:    meta.Category,
			Options:     meta.Options,
			Creator:     meta.Creator.String(),
			Owner:       meta.Owner.String(),
			Opened:      meta.Opened,
			Data:        nftColl.Data,
		})
	}

	return &types.QueryGetNftCollectionsResponse{
		NftCollections: nftCollections,
	}, nil
}
