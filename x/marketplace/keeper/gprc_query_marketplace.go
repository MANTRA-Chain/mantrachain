package keeper

import (
	"context"

	"github.com/LimeChain/mantrachain/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Marketplace(c context.Context, req *types.QueryGetMarketplaceRequest) (*types.QueryGetMarketplaceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	creator, err := sdk.AccAddressFromBech32(req.Creator)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	conf := k.GetParams(ctx)
	err = types.ValidateMarketplaceId(conf.ValidMarketplaceId, req.Id, nil)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	index := types.GetMarketplaceIndex(creator, req.Id)

	meta, found := k.GetMarketplace(
		ctx,
		sdk.AccAddress(creator),
		index,
	)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetMarketplaceResponse{
		Id:          meta.Id,
		Name:        meta.Name,
		Description: meta.Description,
		Url:         meta.Url,
		Creator:     meta.Creator.String(),
		Owner:       meta.Owner.String(),
		Opened:      meta.Opened,
		Options:     meta.Options,
		Attributes:  meta.Attributes,
		Images:      meta.Images,
		Links:       meta.Links,
		Data:        meta.Data,
	}, nil
}
