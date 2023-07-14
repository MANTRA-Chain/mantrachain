package keeper

import (
	"context"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"mantrachain/x/guard/types"
)

func (k msgServer) UpdateRequiredPrivileges(goCtx context.Context, msg *types.MsgUpdateRequiredPrivileges) (*types.MsgUpdateRequiredPrivilegesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.CheckIsAdmin(ctx, msg.GetCreator()); err != nil {
		return nil, sdkerrors.Wrap(err, "unauthorized")
	}

	if len(msg.Index) == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid index")
	}

	kind, err := types.ParseRequiredPrivilegesKind(msg.Kind)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid kind")
	}

	isFound := k.HasRequiredPrivileges(ctx, msg.Index, kind)

	reqPr := types.PrivilegesFromBytes(msg.Privileges)
	updated := false

	if isFound && reqPr.Empty() {
		k.RemoveRequiredPrivileges(ctx, msg.Index, kind)
		updated = true
	} else if !reqPr.Empty() {
		k.SetRequiredPrivileges(ctx, msg.Index, kind, msg.Privileges)
		updated = true
	}

	if updated {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgUpdateRequiredPrivileges),
				sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
				sdk.NewAttribute(types.AttributeKeyIndex, string(msg.Index)),
				sdk.NewAttribute(types.AttributeKeyKind, kind.String()),
			),
		)
	}

	return &types.MsgUpdateRequiredPrivilegesResponse{}, nil
}

func (k msgServer) UpdateRequiredPrivilegesBatch(goCtx context.Context, msg *types.MsgUpdateRequiredPrivilegesBatch) (*types.MsgUpdateRequiredPrivilegesBatchResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.CheckIsAdmin(ctx, msg.GetCreator()); err != nil {
		return nil, sdkerrors.Wrap(err, "unauthorized")
	}

	kind, err := types.ParseRequiredPrivilegesKind(msg.Kind)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid kind")
	}

	indexes := []string{}

	for i, index := range msg.RequiredPrivileges.Indexes {
		if len(index) == 0 {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid index")
		}

		isFound := k.HasRequiredPrivileges(ctx, index, kind)
		reqPr := types.PrivilegesFromBytes(msg.RequiredPrivileges.Privileges[i])

		if isFound && reqPr.Empty() {
			k.RemoveRequiredPrivileges(ctx, index, kind)
			indexes = append(indexes, string(index))
		} else if !reqPr.Empty() {
			k.SetRequiredPrivileges(ctx, index, kind, msg.RequiredPrivileges.Privileges[i])
			indexes = append(indexes, string(index))
		}
	}

	if len(indexes) > 0 {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgUpdateRequiredPrivilegesBatch),
				sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
				sdk.NewAttribute(types.AttributeKeyIndexes, strings.Join(indexes, ",")),
				sdk.NewAttribute(types.AttributeKeyKind, kind.String()),
			),
		)
	}

	return &types.MsgUpdateRequiredPrivilegesBatchResponse{}, nil
}

func (k msgServer) UpdateRequiredPrivilegesGroupedBatch(goCtx context.Context, msg *types.MsgUpdateRequiredPrivilegesGroupedBatch) (*types.MsgUpdateRequiredPrivilegesGroupedBatchResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.CheckIsAdmin(ctx, msg.GetCreator()); err != nil {
		return nil, sdkerrors.Wrap(err, "unauthorized")
	}

	kind, err := types.ParseRequiredPrivilegesKind(msg.Kind)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid kind")
	}

	indexes := []string{}

	for i := range msg.RequiredPrivilegesGrouped.Indexes {
		reqPr := types.PrivilegesFromBytes(msg.RequiredPrivilegesGrouped.Privileges[i])

		for _, index := range msg.RequiredPrivilegesGrouped.Indexes[i].Indexes {
			if len(index) == 0 {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid index")
			}

			isFound := k.HasRequiredPrivileges(ctx, index, kind)

			if isFound && reqPr.Empty() {
				k.RemoveRequiredPrivileges(ctx, index, kind)
				indexes = append(indexes, string(index))
			} else if !reqPr.Empty() {
				k.SetRequiredPrivileges(ctx, index, kind, msg.RequiredPrivilegesGrouped.Privileges[i])
				indexes = append(indexes, string(index))
			}
		}
	}

	if len(indexes) > 0 {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgUpdateRequiredPrivilegesGroupedBatch),
				sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
				sdk.NewAttribute(types.AttributeKeyIndexes, strings.Join(indexes, ",")),
				sdk.NewAttribute(types.AttributeKeyKind, kind.String()),
			),
		)
	}

	return &types.MsgUpdateRequiredPrivilegesGroupedBatchResponse{}, nil
}
