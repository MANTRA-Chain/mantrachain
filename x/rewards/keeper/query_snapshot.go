package keeper

import (
	"context"

	"github.com/MANTRA-Finance/aumega/x/rewards/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) SnapshotAll(goCtx context.Context, req *types.QueryAllSnapshotRequest) (*types.QueryAllSnapshotResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var snapshots []types.Snapshot
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	snapshotStore := prefix.NewStore(store, types.SnapshotStoreKey(req.PairId))

	pageRes, err := query.Paginate(snapshotStore, req.Pagination, func(key []byte, value []byte) error {
		var snapshot types.Snapshot
		if err := k.cdc.Unmarshal(value, &snapshot); err != nil {
			return err
		}

		snapshots = append(snapshots, snapshot)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllSnapshotResponse{Snapshot: snapshots, Pagination: pageRes}, nil
}

func (k Keeper) Snapshot(goCtx context.Context, req *types.QueryGetSnapshotRequest) (*types.QueryGetSnapshotResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	snapshot, found := k.GetSnapshot(ctx, req.PairId, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryGetSnapshotResponse{Snapshot: snapshot}, nil
}
