package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"github.com/MANTRA-Chain/mantrachain/v8/x/sanction/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) Sanctioners(ctx context.Context, req *types.QuerySanctionersRequest) (*types.QuerySanctionersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	sanctioners := []types.Sanctioner{}
	_, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Sanctioners,
		req.Pagination,
		func(_ string, s types.Sanctioner) (bool, error) {
			sanctioners = append(sanctioners, s)
			return true, nil
		},
	)
	if err != nil {
		return nil, err
	}
	return &types.QuerySanctionersResponse{Sanctioners: sanctioners, Pagination: pageRes}, nil
}

func (q queryServer) SanctionedAddresses(ctx context.Context, req *types.QuerySanctionedAddressesRequest) (*types.QuerySanctionedAddressesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	entries := []types.SanctionedAddress{}
	_, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.BlacklistAccounts,
		req.Pagination,
		func(addr string, _ collections.NoValue) (bool, error) {
			_, record, err := q.k.IsSanctioned(ctx, addr)
			if err != nil {
				return false, err
			}
			entries = append(entries, types.SanctionedAddress{Address: addr, Record: record})
			return true, nil
		},
	)
	if err != nil {
		return nil, err
	}
	return &types.QuerySanctionedAddressesResponse{SanctionedAddresses: entries, Pagination: pageRes}, nil
}

func (q queryServer) IsSanctioned(ctx context.Context, req *types.QueryIsSanctionedRequest) (*types.QueryIsSanctionedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	addr, err := types.NormalizeAddress(req.Address)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid address: %s", req.Address)
	}
	sanctioned, record, err := q.k.IsSanctioned(ctx, addr)
	if err != nil {
		return nil, err
	}
	return &types.QueryIsSanctionedResponse{Sanctioned: sanctioned, Record: record}, nil
}
