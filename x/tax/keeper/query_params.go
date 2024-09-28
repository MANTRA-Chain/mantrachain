package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	"github.com/MANTRA-Chain/mantrachain/x/tax/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	params, err := q.k.Params.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}
	// set the max tax rate to the hardcoded value
	params.MaxMcaTax = types.MaxMcaTax
	return &types.QueryParamsResponse{Params: params}, nil
}
