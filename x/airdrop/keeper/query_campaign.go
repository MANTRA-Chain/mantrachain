package keeper

import (
	"context"
	"strconv"

	"cosmossdk.io/store/prefix"
	"github.com/MANTRA-Finance/mantrachain/x/airdrop/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) QueryCampaignAll(goCtx context.Context, req *types.QueryAllCampaignRequest) (*types.QueryAllCampaignResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var campaigns []types.Campaign
	ctx := sdk.UnwrapSDKContext(goCtx)

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.CampaignStoreKey())

	pageRes, err := query.Paginate(store, req.Pagination, func(key []byte, value []byte) error {
		var campaign types.Campaign
		if err := k.cdc.Unmarshal(value, &campaign); err != nil {
			return err
		}

		campaigns = append(campaigns, campaign)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllCampaignResponse{Campaign: campaigns, Pagination: pageRes}, nil
}

func (k Keeper) QueryCampaign(goCtx context.Context, req *types.QueryGetCampaignRequest) (*types.QueryGetCampaignResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetCampaign(
		ctx,
		types.GetCampaignIndex(strconv.FormatUint(req.Id, 10)),
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetCampaignResponse{Campaign: val}, nil
}
