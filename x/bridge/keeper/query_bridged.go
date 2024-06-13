package keeper

import (
	"context"

	"github.com/MANTRA-Finance/mantrachain/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) BridgedAll(goCtx context.Context, req *types.QueryAllBridgedRequest) (*types.QueryAllBridgedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var bridgeds []types.Bridged
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	bridgedStore := prefix.NewStore(store, types.KeyPrefix(types.BridgedKeyPrefix))

	pageRes, err := query.Paginate(bridgedStore, req.Pagination, func(key []byte, value []byte) error {
		var bridged types.Bridged
		if err := k.cdc.Unmarshal(value, &bridged); err != nil {
			return err
		}

		bridgeds = append(bridgeds, bridged)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllBridgedResponse{Bridged: bridgeds, Pagination: pageRes}, nil
}

func (k Keeper) Bridged(goCtx context.Context, req *types.QueryGetBridgedRequest) (*types.QueryGetBridgedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetBridged(
		ctx,
		req.EthTxHash,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetBridgedResponse{Bridged: val}, nil
}
