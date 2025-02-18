package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/MANTRA-Chain/mantrachain/v2/x/sanction/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) AddBlacklistAccount(ctx context.Context, msg *types.MsgAddBlacklistAccount) (*types.MsgAddBlacklistAccountResponse, error) {
	if err := msg.Validate(); err != nil {
		return nil, err
	}

	authority := k.GetAuthority()
	if authority != msg.Authority {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid authority; expected %s, got %s", authority, msg.Authority)
	}

	hasAccount, err := k.BlacklistAccounts.Has(ctx, msg.BlacklistAccount)
	if err != nil {
		return nil, err
	}
	if hasAccount {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "account %s has already been blacklisted", msg.BlacklistAccount)
	}

	if err := k.BlacklistAccounts.Set(ctx, msg.BlacklistAccount); err != nil {
		return nil, err
	}

	return &types.MsgAddBlacklistAccountResponse{}, nil
}

func (k msgServer) RemoveBlacklistAccount(ctx context.Context, msg *types.MsgRemoveBlacklistAccount) (*types.MsgRemoveBlacklistAccountResponse, error) {
	if err := msg.Validate(); err != nil {
		return nil, err
	}

	authority := k.GetAuthority()
	if authority != msg.Authority {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid authority; expected %s, got %s", authority, msg.Authority)
	}

	hasAccount, err := k.BlacklistAccounts.Has(ctx, msg.BlacklistAccount)
	if err != nil {
		return nil, err
	}
	if !hasAccount {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "blacklist account %s is not blacklisted", msg.BlacklistAccount)
	}

	if err := k.BlacklistAccounts.Remove(ctx, msg.BlacklistAccount); err != nil {
		return nil, err
	}

	return &types.MsgRemoveBlacklistAccountResponse{}, nil
}
