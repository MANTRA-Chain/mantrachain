package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	coinfactorytypes "mantrachain/x/coinfactory/types"
	"mantrachain/x/guard/types"
)

func (k Keeper) LockedAll(goCtx context.Context, req *types.QueryAllLockedRequest) (*types.QueryAllLockedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var indexes [][]byte
	var locked []bool
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)

	kind, err := types.ParseLockedKind(req.Kind)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "kind is invalid")
	}

	lockedStore := prefix.NewStore(store, types.LockedStoreKey(kind.Bytes()))

	pageRes, err := query.Paginate(lockedStore, req.Pagination, func(key []byte, value []byte) error {
		indexes = append(indexes, key)
		locked = append(locked, true)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllLockedResponse{
		Indexes:    indexes,
		Locked:     locked,
		Kind:       req.Kind,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) Locked(goCtx context.Context, req *types.QueryGetLockedRequest) (*types.QueryGetLockedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	kind, err := types.ParseLockedKind(req.Kind)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "kind is invalid")
	}

	denom := string(req.Index)

	if kind == types.LockedCoin {
		_, _, err := coinfactorytypes.DeconstructDenom(denom)
		if err != nil {
			return nil, err
		}

		_, found := k.bk.GetDenomMetaData(ctx, denom)
		if !found {
			return nil, status.Error(codes.InvalidArgument, "invalid denom")
		}
	}

	_, found := k.GetLocked(
		ctx,
		req.Index,
		kind,
	)

	return &types.QueryGetLockedResponse{
		Index:  req.Index,
		Locked: found,
		Kind:   req.Kind,
	}, nil
}
