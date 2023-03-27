package keeper

import (
	"context"
	"strconv"

	"cosmossdk.io/errors"
	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) UpdateLocked(goCtx context.Context, msg *types.MsgUpdateLocked) (*types.MsgUpdateLockedResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if len(msg.Index) == 0 {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid index")
	}

	kind, err := types.ParseLockedKind(msg.Kind)
	if err != nil {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid kind")
	}

	exists := k.HasLocked(ctx, msg.Index, kind)
	updated := false

	if exists && !msg.Locked {
		k.RemoveLocked(ctx, msg.Index, kind)
		updated = true
	} else if !exists && msg.Locked {
		k.SetLocked(ctx, msg.Index, kind)
		updated = true
	}

	if updated {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgUpdateLocked),
				sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
				sdk.NewAttribute(types.AttributeKeyLocked, strconv.FormatBool(msg.Locked)),
				sdk.NewAttribute(types.AttributeKeyIndex, string(msg.Index)),
				sdk.NewAttribute(types.AttributeKeyKind, kind.String()),
			),
		)
	}

	return &types.MsgUpdateLockedResponse{}, nil
}
