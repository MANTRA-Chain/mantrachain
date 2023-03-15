package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/LimeChain/mantrachain/x/guard/types"
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

	exists := k.HasLocked(
		ctx,
		msg.Index,
		kind,
	)

	if exists && !msg.Locked {
		k.RemoveLocked(ctx, msg.Index, kind)
	} else if !exists && msg.Locked {
		k.SetLocked(ctx, msg.Index, kind)
	}

	// TODO: add event

	return &types.MsgUpdateLockedResponse{}, nil
}
