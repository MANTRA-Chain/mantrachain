package keeper

import (
	"context"

	"github.com/LimeChain/mantrachain/x/mdb/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Nft(c context.Context, req *types.QueryGetNftRequest) (*types.QueryGetNftResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	collectionCreator, err := sdk.AccAddressFromBech32(req.CollectionCreator)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.CollectionId == "" {
		return nil, status.Error(codes.InvalidArgument, "empty collection id")
	}

	var collIndex []byte

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "empty nft id")
	} else {
		collIndex = types.GetNftCollectionIndex(collectionCreator, req.CollectionId)

		if !k.HasNftCollection(ctx, collectionCreator, collIndex) {
			return nil, status.Error(codes.InvalidArgument, "collection not exists")
		}
	}

	index := types.GetNftIndex(collIndex, req.Id)

	meta, found := k.GetNft(
		ctx,
		collIndex,
		index,
	)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)
	nft, found := nftExecutor.GetNft(ctx, string(collIndex), string(index))
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetNftResponse{
		Id:           nft.UriHash,
		Title:        meta.Title,
		Did:          meta.Did,
		Description:  meta.Description,
		Images:       meta.Images,
		Url:          meta.Url,
		Links:        meta.Links,
		Attributes:   meta.Attributes,
		Creator:      meta.Creator.String(),
		Owner:        meta.Owner.String(),
		Resellable:   meta.Resellable,
		Data:         nft.Data,
		CollectionId: meta.CollectionId,
	}, nil
}
