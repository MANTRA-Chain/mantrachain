package keeper

import (
	"context"

	"github.com/LimeChain/mantrachain/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Address(c context.Context, req *types.QueryAddressRequest) (*types.QueryAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryAddressResponse{Address: k.GetAddress(ctx).String()}, nil
}
