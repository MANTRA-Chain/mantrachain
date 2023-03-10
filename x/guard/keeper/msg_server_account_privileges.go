package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/LimeChain/mantrachain/x/guard/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) UpdateAccountPrivileges(goCtx context.Context, msg *types.MsgUpdateAccountPrivileges) (*types.MsgUpdateAccountPrivilegesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	conf := k.GetParams(ctx)

	account, err := sdk.AccAddressFromBech32(msg.Account)
	if err != nil {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid account address")
	}

	isFound := k.HasAccountPrivileges(
		ctx,
		account,
	)

	accPr := types.AccPrivilegesFromBytes(msg.Privileges)
	areDefaultAccountPrivileges := accPr.Equal(conf.DefaultAccountPrivileges)

	if isFound && (accPr.Empty() || areDefaultAccountPrivileges) {
		k.RemoveAccountPrivileges(ctx, account)
	} else if !accPr.Empty() && !areDefaultAccountPrivileges {
		k.SetAccountPrivileges(ctx, account, msg.Privileges)
	}

	// TODO: add event

	return &types.MsgUpdateAccountPrivilegesResponse{}, nil
}

func (k msgServer) UpdateAccountPrivilegesBatch(goCtx context.Context, msg *types.MsgUpdateAccountPrivilegesBatch) (*types.MsgUpdateAccountPrivilegesBatchResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	conf := k.GetParams(ctx)

	if msg.AccountsPrivileges == nil || len(msg.AccountsPrivileges.Accounts) == 0 {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "accounts and/or privileges are empty")
	}

	if msg.AccountsPrivileges.Privileges == nil || len(msg.AccountsPrivileges.Accounts) != len(msg.AccountsPrivileges.Privileges) {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "accounts and privileges length are not equal")
	}

	for i, acc := range msg.AccountsPrivileges.Accounts {
		account, err := sdk.AccAddressFromBech32(acc)
		if err != nil {
			return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid account address")
		}

		isFound := k.HasAccountPrivileges(
			ctx,
			account,
		)

		accPr := types.AccPrivilegesFromBytes(msg.AccountsPrivileges.Privileges[i])
		areDefaultAccountPrivileges := accPr.Equal(conf.DefaultAccountPrivileges)

		if isFound && (accPr.Empty() || areDefaultAccountPrivileges) {
			k.RemoveAccountPrivileges(ctx, account)
		} else if !accPr.Empty() && !areDefaultAccountPrivileges {
			k.SetAccountPrivileges(ctx, account, msg.AccountsPrivileges.Privileges[i])
		}
	}

	// TODO: add event

	return &types.MsgUpdateAccountPrivilegesBatchResponse{}, nil
}

func (k msgServer) UpdateAccountPrivilegesGroupedBatch(goCtx context.Context, msg *types.MsgUpdateAccountPrivilegesGroupedBatch) (*types.MsgUpdateAccountPrivilegesGroupedBatchResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	conf := k.GetParams(ctx)

	if msg.AccountsPrivilegesGrouped == nil || msg.AccountsPrivilegesGrouped.Accounts == nil || len(msg.AccountsPrivilegesGrouped.Accounts) == 0 {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "grouped accounts and/or privileges are empty")
	}

	if msg.AccountsPrivilegesGrouped.Privileges == nil || len(msg.AccountsPrivilegesGrouped.Accounts) != len(msg.AccountsPrivilegesGrouped.Privileges) {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "grouped accounts and privileges length is not equal")
	}

	for i := range msg.AccountsPrivilegesGrouped.Accounts {
		accPr := types.AccPrivilegesFromBytes(msg.AccountsPrivilegesGrouped.Privileges[i])
		areDefaultAccountPrivileges := accPr.Equal(conf.DefaultAccountPrivileges)

		for _, acc := range msg.AccountsPrivilegesGrouped.Accounts[i].Accounts {
			account, err := sdk.AccAddressFromBech32(acc)
			if err != nil {
				return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid account address")
			}

			isFound := k.HasAccountPrivileges(
				ctx,
				account,
			)

			if isFound && (accPr.Empty() || areDefaultAccountPrivileges) {
				k.RemoveAccountPrivileges(ctx, account)
			} else if !accPr.Empty() && !areDefaultAccountPrivileges {
				k.SetAccountPrivileges(ctx, account, msg.AccountsPrivilegesGrouped.Privileges[i])
			}
		}
	}

	// TODO: add event

	return &types.MsgUpdateAccountPrivilegesGroupedBatchResponse{}, nil
}
