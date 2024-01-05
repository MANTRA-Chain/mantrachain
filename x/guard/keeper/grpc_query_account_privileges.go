package keeper

import (
	"context"

	"github.com/MANTRA-Finance/aumega/x/guard/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AccountPrivilegesAll(goCtx context.Context, req *types.QueryAllAccountPrivilegesRequest) (*types.QueryAllAccountPrivilegesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var accounts []string
	var privileges [][]byte
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	accountPrivilegesStore := prefix.NewStore(store, types.AccountPrivilegesStoreKey())

	pageRes, err := query.Paginate(accountPrivilegesStore, req.Pagination, func(key []byte, value []byte) error {
		accounts = append(accounts, sdk.AccAddress(key).String())
		privileges = append(privileges, value)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllAccountPrivilegesResponse{Accounts: accounts, Privileges: privileges, Pagination: pageRes}, nil
}

func (k Keeper) AccountPrivileges(c context.Context, req *types.QueryGetAccountPrivilegesRequest) (*types.QueryGetAccountPrivilegesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	conf := k.GetParams(ctx)

	account, err := sdk.AccAddressFromBech32(req.Account)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid account address")
	}

	val, found := k.GetAccountPrivileges(
		ctx,
		account,
		conf.DefaultPrivileges,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetAccountPrivilegesResponse{
		Account:    req.Account,
		Privileges: val,
	}, nil
}
