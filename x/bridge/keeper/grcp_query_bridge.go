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

	bridgeCreator, err := sdk.AccAddressFromBech32(req.BridgeCreator)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	conf := k.GetParams(ctx)
	err = types.ValidateBridgeId(conf.ValidBridgeId, req.BridgeId)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	index := types.GetBridgeIndex(bridgeCreator, req.BridgeId)

	bridge, found := k.GetBridge(
		ctx,
		sdk.AccAddress(bridgeCreator),
		index,
	)

	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetBridgeResponse{
		BridgeCreator:       bridge.Creator.String(),
		BridgeId:            bridge.Id,
		BridgeAccount:       bridge.BridgeAccount,
		Cw20ContractAddress: bridge.Cw20ContractAddress,
		Owner:               bridge.Owner.String(),
	}, nil
}
