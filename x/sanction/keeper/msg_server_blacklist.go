package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/MANTRA-Chain/mantrachain/v8/x/sanction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
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

		// revoke all authz grants for the blacklisted account
		grants, err := k.authzKeeper.GranterGrants(ctx, &authz.QueryGranterGrantsRequest{Granter: account})
		if err != nil {
			return nil, err
		}
		granter, err := sdk.AccAddressFromBech32(account)
		if err != nil {
			return nil, err
		}
		for _, grant := range grants.Grants {
			grantee, err := sdk.AccAddressFromBech32(grant.Grantee)
			if err != nil {
				return nil, err
			}

			var authorization authz.Authorization
			err = k.cdc.UnpackAny(grant.Authorization, &authorization)
			if err != nil {
				return nil, err
			}

			if err := k.authzKeeper.DeleteGrant(ctx, grantee, granter, authorization.MsgTypeURL()); err != nil {
				return nil, err
			}
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
