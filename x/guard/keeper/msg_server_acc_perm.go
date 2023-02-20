package keeper

import (
	"context"

	"github.com/LimeChain/mantrachain/x/guard/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateAccPerm(goCtx context.Context, msg *types.MsgCreateAccPerm) (*types.MsgCreateAccPermResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	conf := k.GetParams(ctx)

	creator, err := sdk.AccAddressFromBech32(msg.Creator)

	if err != nil {
		return nil, err
	}

	adminAccount, err := sdk.AccAddressFromBech32(conf.AdminAccount)

	if err != nil {
		return nil, err
	}

	if !creator.Equals(adminAccount) {
		return nil, sdkerrors.Wrapf(types.ErrAdminAccountParamMismatch, "admin account param %s does not match the creator %s", adminAccount.String(), creator.String())
	}

	// Check if the value already exists
	_, isFound := k.GetAccPerm(
		ctx,
		msg.Id,
	)
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var accPerm = types.AccPerm{
		Creator:    msg.Creator,
		Id:         msg.Id,
		Priviliges: msg.Priviliges,
	}

	// TODO: add event
	k.SetAccPerm(
		ctx,
		accPerm,
	)
	return &types.MsgCreateAccPermResponse{}, nil
}

func (k msgServer) UpdateAccPerm(goCtx context.Context, msg *types.MsgUpdateAccPerm) (*types.MsgUpdateAccPermResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	conf := k.GetParams(ctx)

	creator, err := sdk.AccAddressFromBech32(msg.Creator)

	if err != nil {
		return nil, err
	}

	adminAccount, err := sdk.AccAddressFromBech32(conf.AdminAccount)

	if err != nil {
		return nil, err
	}

	if !creator.Equals(adminAccount) {
		return nil, sdkerrors.Wrapf(types.ErrAdminAccountParamMismatch, "admin account param %s does not match the creator %s", adminAccount.String(), creator.String())
	}

	// Check if the value exists
	valFound, isFound := k.GetAccPerm(
		ctx,
		msg.Id,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != valFound.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	var accPerm = types.AccPerm{
		Creator:    msg.Creator,
		Id:         msg.Id,
		Priviliges: msg.Priviliges,
	}

	// TODO: add event
	k.SetAccPerm(ctx, accPerm)

	return &types.MsgUpdateAccPermResponse{}, nil
}

func (k msgServer) DeleteAccPerm(goCtx context.Context, msg *types.MsgDeleteAccPerm) (*types.MsgDeleteAccPermResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	conf := k.GetParams(ctx)

	creator, err := sdk.AccAddressFromBech32(msg.Creator)

	if err != nil {
		return nil, err
	}

	adminAccount, err := sdk.AccAddressFromBech32(conf.AdminAccount)

	if err != nil {
		return nil, err
	}

	if !creator.Equals(adminAccount) {
		return nil, sdkerrors.Wrapf(types.ErrAdminAccountParamMismatch, "admin account param %s does not match the creator %s", adminAccount.String(), creator.String())
	}

	// Check if the value exists
	valFound, isFound := k.GetAccPerm(
		ctx,
		msg.Id,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != valFound.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	// TODO: add event
	k.RemoveAccPerm(
		ctx,
		msg.Id,
	)

	return &types.MsgDeleteAccPermResponse{}, nil
}
