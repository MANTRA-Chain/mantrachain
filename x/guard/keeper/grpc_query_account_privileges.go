package keeper

import (
	"context"

	"github.com/LimeChain/mantrachain/x/guard/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
		conf.DefaultAccountPrivileges,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetAccountPrivilegesResponse{
		Account:    req.Account,
		Privileges: val,
	}, nil
}

func (k Keeper) AccountPrivilegesMany(c context.Context, req *types.QueryGetAccountPrivilegesManyRequest) (*types.QueryGetAccountPrivilegesManyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	conf := k.GetParams(ctx)

	if req.Accounts == nil || len(req.Accounts) == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var accounts []sdk.AccAddress

	for _, account := range req.Accounts {
		acc, err := sdk.AccAddressFromBech32(account)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid account address")
		}
		accounts = append(accounts, acc)
	}

	values := k.GetAccountPrivilegesMany(
		ctx,
		accounts,
		conf.DefaultAccountPrivileges,
	)

	return &types.QueryGetAccountPrivilegesManyResponse{
		Accounts:   req.Accounts,
		Privileges: values,
	}, nil
}
