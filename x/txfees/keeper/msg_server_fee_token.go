package keeper

import (
	"context"

	"github.com/AumegaChain/aumega/x/txfees/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateFeeToken(goCtx context.Context, msg *types.MsgCreateFeeToken) (*types.MsgCreateFeeTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.gk.CheckIsAdmin(ctx, msg.Creator); err != nil {
		return nil, sdkerrors.Wrap(err, "unauthorized")
	}

	// Check if the value already exists
	_, isFound := k.GetFeeToken(
		ctx,
		msg.Denom,
	)
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var feeToken = types.FeeToken{
		Denom:  msg.Denom,
		PairId: msg.PairId,
	}

	k.SetFeeToken(
		ctx,
		feeToken,
	)
	return &types.MsgCreateFeeTokenResponse{}, nil
}

func (k msgServer) UpdateFeeToken(goCtx context.Context, msg *types.MsgUpdateFeeToken) (*types.MsgUpdateFeeTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.gk.CheckIsAdmin(ctx, msg.Creator); err != nil {
		return nil, sdkerrors.Wrap(err, "unauthorized")
	}

	// Check if the value exists
	_, isFound := k.GetFeeToken(
		ctx,
		msg.Denom,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	var feeToken = types.FeeToken{
		Denom:  msg.Denom,
		PairId: msg.PairId,
	}

	k.SetFeeToken(ctx, feeToken)

	return &types.MsgUpdateFeeTokenResponse{}, nil
}

func (k msgServer) DeleteFeeToken(goCtx context.Context, msg *types.MsgDeleteFeeToken) (*types.MsgDeleteFeeTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.gk.CheckIsAdmin(ctx, msg.Creator); err != nil {
		return nil, sdkerrors.Wrap(err, "unauthorized")
	}

	// Check if the value exists
	_, isFound := k.GetFeeToken(
		ctx,
		msg.Denom,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	k.RemoveFeeToken(
		ctx,
		msg.Denom,
	)

	return &types.MsgDeleteFeeTokenResponse{}, nil
}
