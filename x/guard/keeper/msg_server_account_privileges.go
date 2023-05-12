package keeper

import (
	"context"
	"strings"

	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) UpdateAccountPrivileges(goCtx context.Context, msg *types.MsgUpdateAccountPrivileges) (*types.MsgUpdateAccountPrivilegesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.CheckIsAdmin(ctx, msg.GetCreator()); err != nil {
		return nil, sdkerrors.Wrap(err, "unauthorized")
	}

	conf := k.GetParams(ctx)

	account, err := sdk.AccAddressFromBech32(msg.Account)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid account address")
	}

	isFound := k.HasAccountPrivileges(ctx, account)

	accPr := types.PrivilegesFromBytes(msg.Privileges)
	areDefaultAccountPrivileges := accPr.Equal(conf.DefaultPrivileges)
	updated := false

	if isFound && (accPr.Empty() || areDefaultAccountPrivileges) {
		k.RemoveAccountPrivileges(ctx, account)
		updated = true
	} else if !accPr.Empty() && !areDefaultAccountPrivileges {
		k.SetAccountPrivileges(ctx, account, msg.Privileges)
		updated = true
	}

	if updated {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgUpdateAccountPrivileges),
				sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
				sdk.NewAttribute(types.AttributeKeyAccount, msg.Account),
			),
		)
	}

	return &types.MsgUpdateAccountPrivilegesResponse{}, nil
}

func (k msgServer) UpdateAccountPrivilegesBatch(goCtx context.Context, msg *types.MsgUpdateAccountPrivilegesBatch) (*types.MsgUpdateAccountPrivilegesBatchResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.CheckIsAdmin(ctx, msg.GetCreator()); err != nil {
		return nil, sdkerrors.Wrap(err, "unauthorized")
	}

	conf := k.GetParams(ctx)
	accounts := []string{}

	for i, acc := range msg.AccountsPrivileges.Accounts {
		account, err := sdk.AccAddressFromBech32(acc)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid account address")
		}

		isFound := k.HasAccountPrivileges(ctx, account)

		accPr := types.PrivilegesFromBytes(msg.AccountsPrivileges.Privileges[i])
		areDefaultAccountPrivileges := accPr.Equal(conf.DefaultPrivileges)

		if isFound && (accPr.Empty() || areDefaultAccountPrivileges) {
			k.RemoveAccountPrivileges(ctx, account)
			accounts = append(accounts, acc)
		} else if !accPr.Empty() && !areDefaultAccountPrivileges {
			k.SetAccountPrivileges(ctx, account, msg.AccountsPrivileges.Privileges[i])
			accounts = append(accounts, acc)
		}
	}

	if len(accounts) > 0 {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgUpdateAccountPrivilegesBatch),
				sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
				sdk.NewAttribute(types.AttributeKeyAccounts, strings.Join(accounts, ",")),
			),
		)
	}

	return &types.MsgUpdateAccountPrivilegesBatchResponse{}, nil
}

func (k msgServer) UpdateAccountPrivilegesGroupedBatch(goCtx context.Context, msg *types.MsgUpdateAccountPrivilegesGroupedBatch) (*types.MsgUpdateAccountPrivilegesGroupedBatchResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.CheckIsAdmin(ctx, msg.GetCreator()); err != nil {
		return nil, sdkerrors.Wrap(err, "unauthorized")
	}

	conf := k.GetParams(ctx)
	accounts := []string{}

	for i := range msg.AccountsPrivilegesGrouped.Accounts {
		accPr := types.PrivilegesFromBytes(msg.AccountsPrivilegesGrouped.Privileges[i])
		areDefaultAccountPrivileges := accPr.Equal(conf.DefaultPrivileges)

		for _, acc := range msg.AccountsPrivilegesGrouped.Accounts[i].Accounts {
			account, err := sdk.AccAddressFromBech32(acc)
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid account address")
			}

			isFound := k.HasAccountPrivileges(ctx, account)

			if isFound && (accPr.Empty() || areDefaultAccountPrivileges) {
				k.RemoveAccountPrivileges(ctx, account)
				accounts = append(accounts, acc)
			} else if !accPr.Empty() && !areDefaultAccountPrivileges {
				k.SetAccountPrivileges(ctx, account, msg.AccountsPrivilegesGrouped.Privileges[i])
				accounts = append(accounts, acc)
			}
		}
	}

	if len(accounts) > 0 {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgUpdateAccountPrivilegesGroupedBatch),
				sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
				sdk.NewAttribute(types.AttributeKeyAccounts, strings.Join(accounts, ",")),
			),
		)
	}

	return &types.MsgUpdateAccountPrivilegesGroupedBatchResponse{}, nil
}
