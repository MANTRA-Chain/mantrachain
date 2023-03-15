package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/LimeChain/mantrachain/x/guard/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) UpdateRequiredPrivileges(goCtx context.Context, msg *types.MsgUpdateRequiredPrivileges) (*types.MsgUpdateRequiredPrivilegesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if len(msg.Index) == 0 {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid index")
	}

	kind, err := types.ParseRequiredPrivilegesKind(msg.Kind)
	if err != nil {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid kind")
	}

	isFound := k.HasRequiredPrivileges(
		ctx,
		msg.Index,
		kind,
	)

	reqPr := types.PrivilegesFromBytes(msg.Privileges)

	if isFound && reqPr.Empty() {
		k.RemoveRequiredPrivileges(ctx, msg.Index, kind)
	} else if !reqPr.Empty() {
		k.SetRequiredPrivileges(ctx, msg.Index, kind, msg.Privileges)
	}

	// TODO: add event

	return &types.MsgUpdateRequiredPrivilegesResponse{}, nil
}

func (k msgServer) UpdateRequiredPrivilegesBatch(goCtx context.Context, msg *types.MsgUpdateRequiredPrivilegesBatch) (*types.MsgUpdateRequiredPrivilegesBatchResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.RequiredPrivilegesList == nil || len(msg.RequiredPrivilegesList.Indexes) == 0 {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "indexes and/or privileges are empty")
	}

	if msg.RequiredPrivilegesList.Privileges == nil || len(msg.RequiredPrivilegesList.Indexes) != len(msg.RequiredPrivilegesList.Privileges) {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "indexes and privileges length are not equal")
	}

	kind, err := types.ParseRequiredPrivilegesKind(msg.Kind)
	if err != nil {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid kind")
	}

	for i, index := range msg.RequiredPrivilegesList.Indexes {
		if len(index) == 0 {
			return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid index")
		}

		isFound := k.HasRequiredPrivileges(
			ctx,
			index,
			kind,
		)

		reqPr := types.PrivilegesFromBytes(msg.RequiredPrivilegesList.Privileges[i])

		if isFound && reqPr.Empty() {
			k.RemoveRequiredPrivileges(ctx, index, kind)
		} else if !reqPr.Empty() {
			k.SetRequiredPrivileges(ctx, index, kind, msg.RequiredPrivilegesList.Privileges[i])
		}
	}

	// TODO: add event

	return &types.MsgUpdateRequiredPrivilegesBatchResponse{}, nil
}

func (k msgServer) UpdateRequiredPrivilegesGroupedBatch(goCtx context.Context, msg *types.MsgUpdateRequiredPrivilegesGroupedBatch) (*types.MsgUpdateRequiredPrivilegesGroupedBatchResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.RequiredPrivilegesListGrouped == nil ||
		msg.RequiredPrivilegesListGrouped.Indexes == nil ||
		len(msg.RequiredPrivilegesListGrouped.Indexes) == 0 {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "grouped accounts and/or privileges are empty")
	}

	if msg.RequiredPrivilegesListGrouped.Privileges == nil ||
		len(msg.RequiredPrivilegesListGrouped.Indexes) != len(msg.RequiredPrivilegesListGrouped.Privileges) {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "grouped accounts and privileges length is not equal")
	}

	kind, err := types.ParseRequiredPrivilegesKind(msg.Kind)
	if err != nil {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid kind")
	}

	for i := range msg.RequiredPrivilegesListGrouped.Indexes {
		reqPr := types.PrivilegesFromBytes(msg.RequiredPrivilegesListGrouped.Privileges[i])

		for _, index := range msg.RequiredPrivilegesListGrouped.Indexes[i].Indexes {
			if len(index) == 0 {
				return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid index")
			}

			isFound := k.HasRequiredPrivileges(
				ctx,
				index,
				kind,
			)

			if isFound && reqPr.Empty() {
				k.RemoveRequiredPrivileges(ctx, index, kind)
			} else if !reqPr.Empty() {
				k.SetRequiredPrivileges(ctx, index, kind, msg.RequiredPrivilegesListGrouped.Privileges[i])
			}
		}
	}

	// TODO: add event

	return &types.MsgUpdateRequiredPrivilegesGroupedBatchResponse{}, nil
}
