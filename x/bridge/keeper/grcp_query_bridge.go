package keeper

import (
	"context"

	"github.com/LimeChain/mantrachain/x/bridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Bridge(c context.Context, req *types.QueryGetBridgeRequest) (*types.QueryGetBridgeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	creator, err := sdk.AccAddressFromBech32(req.Creator)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	conf := k.GetParams(ctx)
	err = types.ValidateBridgeId(conf.ValidBridgeId, req.Id)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	index := types.GetBridgeIndex(creator, req.Id)

	bridge, found := k.GetBridge(
		ctx,
		sdk.AccAddress(creator),
		index,
	)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetBridgeResponse{
		Id:                  bridge.Id,
		BridgeAccount:       bridge.BridgeAccount,
		Cw20ContractAddress: bridge.Cw20ContractAddress,
		Creator:             bridge.Creator.String(),
		Owner:               bridge.Owner.String(),
	}, nil
}
