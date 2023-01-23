package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/LimeChain/mantrachain/x/guard/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AccPermAll(c context.Context, req *types.QueryAllAccPermRequest) (*types.QueryAllAccPermResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var accPerms []types.AccPerm
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	accPermStore := prefix.NewStore(store, types.KeyPrefix(types.AccPermKeyPrefix))

	pageRes, err := query.Paginate(accPermStore, req.Pagination, func(key []byte, value []byte) error {
		var accPerm types.AccPerm
		if err := k.cdc.Unmarshal(value, &accPerm); err != nil {
			return err
		}

		accPerms = append(accPerms, accPerm)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllAccPermResponse{AccPerm: accPerms, Pagination: pageRes}, nil
}

func (k Keeper) AccPerm(c context.Context, req *types.QueryGetAccPermRequest) (*types.QueryGetAccPermResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetAccPerm(
	    ctx,
	    req.Cat,
        )
	if !found {
	    return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetAccPermResponse{AccPerm: val}, nil
}