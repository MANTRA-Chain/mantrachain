package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/MANTRA-Chain/mantrachain/v8/x/sanction/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) AddBlacklistAccounts(ctx context.Context, msg *types.MsgAddBlacklistAccounts) (*types.MsgAddBlacklistAccountsResponse, error) {
	if err := msg.Validate(); err != nil {
		return nil, err
	}

	authority := k.GetAuthority()
	if authority != msg.Authority {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid authority; expected %s, got %s", authority, msg.Authority)
	}

	for _, account := range msg.BlacklistAccounts {
		hasAccount, err := k.BlacklistAccounts.Has(ctx, account)
		if err != nil {
			return nil, err
		}
		if hasAccount {
			return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "account %s has already been blacklisted", account)
		}

		if err := k.BlacklistAccounts.Set(ctx, account); err != nil {
			return nil, err
		}
	}

	return &types.MsgAddBlacklistAccountsResponse{}, nil
}

func (k msgServer) RemoveBlacklistAccounts(ctx context.Context, msg *types.MsgRemoveBlacklistAccounts) (*types.MsgRemoveBlacklistAccountsResponse, error) {
	if err := msg.Validate(); err != nil {
		return nil, err
	}

	authority := k.GetAuthority()
	if authority != msg.Authority {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid authority; expected %s, got %s", authority, msg.Authority)
	}

	for _, account := range msg.BlacklistAccounts {
		hasAccount, err := k.BlacklistAccounts.Has(ctx, account)
		if err != nil {
			return nil, err
		}
		if !hasAccount {
			return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "blacklist account %s is not blacklisted", account)
		}

		if err := k.BlacklistAccounts.Remove(ctx, account); err != nil {
			return nil, err
		}
	}

	return &types.MsgRemoveBlacklistAccountsResponse{}, nil
}
