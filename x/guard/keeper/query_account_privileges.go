package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) QueryAccountPrivilegesAll(goCtx context.Context, req *types.QueryAllAccountPrivilegesRequest) (*types.QueryAllAccountPrivilegesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var accounts []string
	var privileges [][]byte
	ctx := sdk.UnwrapSDKContext(goCtx)

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.AccountPrivilegesStoreKey())

	pageRes, err := query.Paginate(store, req.Pagination, func(key []byte, value []byte) error {
		accounts = append(accounts, sdk.AccAddress(key).String())
		privileges = append(privileges, value)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllAccountPrivilegesResponse{Accounts: accounts, Privileges: privileges, Pagination: pageRes}, nil
}

func (k Keeper) QueryAccountPrivileges(c context.Context, req *types.QueryGetAccountPrivilegesRequest) (*types.QueryGetAccountPrivilegesResponse, error) {
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
