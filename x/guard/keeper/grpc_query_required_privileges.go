package keeper

import (
	"context"

	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) RequiredPrivilegesAll(goCtx context.Context, req *types.QueryAllRequiredPrivilegesRequest) (*types.QueryAllRequiredPrivilegesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var indexes [][]byte
	var privileges [][]byte
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)

	kind, err := types.ParseRequiredPrivilegesKind(req.Kind)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "kind is invalid")
	}

	requiredPrivilegesStore := prefix.NewStore(store, types.RequiredPrivilegesStoreKey(kind.Bytes()))

	pageRes, err := query.Paginate(requiredPrivilegesStore, req.Pagination, func(key []byte, value []byte) error {
		indexes = append(indexes, key)
		privileges = append(privileges, value)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllRequiredPrivilegesResponse{
		Indexes:    indexes,
		Privileges: privileges,
		Kind:       req.Kind,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) RequiredPrivileges(goCtx context.Context, req *types.QueryGetRequiredPrivilegesRequest) (*types.QueryGetRequiredPrivilegesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	kind, err := types.ParseRequiredPrivilegesKind(req.Kind)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "kind is invalid")
	}

	val, found := k.GetRequiredPrivileges(
		ctx,
		req.Index,
		kind,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetRequiredPrivilegesResponse{
		Index:      req.Index,
		Privileges: val,
		Kind:       req.Kind,
	}, nil
}
