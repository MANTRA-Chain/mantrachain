package keeper

import (
	"context"

	"github.com/MANTRA-Finance/mantrachain/x/txfees/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) QueryGasEstimation(goCtx context.Context, req *types.QueryGetGasEstimationRequest) (*types.QueryGetGasEstimationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetFeeToken(ctx, req.Denom)

	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	coin, err := sdk.ParseCoinNormalized(req.Amount)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	swapAmount, _, err := k.liquidityKeeper.GetSwapAmount(ctx, val.PairId, coin)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetGasEstimationResponse{Amount: swapAmount}, nil
}
