package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/LimeChain/mantrachain/x/bridge/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Cw20Contract(c context.Context, req *types.QueryGetCw20ContractRequest) (*types.QueryGetCw20ContractResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetCw20Contract(ctx)
	if !found {
	    return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetCw20ContractResponse{Cw20Contract: val}, nil
}