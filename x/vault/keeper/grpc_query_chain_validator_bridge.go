package keeper

import (
	"context"

	"github.com/LimeChain/mantrachain/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ChainValidatorBridge(c context.Context, req *types.QueryGetChainValidatorBridgeRequest) (*types.QueryGetChainValidatorBridgeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetChainValidatorBridge(
		ctx,
		req.Chain,
		req.Validator,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetChainValidatorBridgeResponse{ChainValidatorBridge: val}, nil
}
