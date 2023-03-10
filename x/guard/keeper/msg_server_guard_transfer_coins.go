package keeper

import (
	"context"

	"github.com/LimeChain/mantrachain/x/guard/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateGuardTransferCoins(goCtx context.Context, msg *types.MsgUpdateGuardTransferCoins) (*types.MsgUpdateGuardTransferCoinsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	exists := k.HasGuardTransferCoins(ctx)

	if exists && !msg.Enabled {
		k.RemoveGuardTransferCoins(ctx)
	} else if !exists && msg.Enabled {
		k.SetGuardTransferCoins(ctx)
	}

	// TODO: add event

	return &types.MsgUpdateGuardTransferCoinsResponse{}, nil
}
