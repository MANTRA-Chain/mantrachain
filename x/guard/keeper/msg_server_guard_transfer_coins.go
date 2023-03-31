package keeper

import (
	"context"
	"strconv"

	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateGuardTransferCoins(goCtx context.Context, msg *types.MsgUpdateGuardTransferCoins) (*types.MsgUpdateGuardTransferCoinsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	exists := k.HasGuardTransferCoins(ctx)
	updated := false

	if exists && !msg.Enabled {
		k.RemoveGuardTransferCoins(ctx)
		updated = true
	} else if !exists && msg.Enabled {
		k.SetGuardTransferCoins(ctx)
		updated = true
	}

	if updated {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgUpdateGuardTransferCoins),
				sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
				sdk.NewAttribute(types.AttributeKeyGuardTransferCoins, strconv.FormatBool(msg.Enabled)),
			),
		)
	}

	return &types.MsgUpdateGuardTransferCoinsResponse{}, nil
}