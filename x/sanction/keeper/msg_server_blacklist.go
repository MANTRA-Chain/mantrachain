package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/x/feegrant"
	"github.com/MANTRA-Chain/mantrachain/v8/x/sanction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
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
		granter, err := sdk.AccAddressFromBech32(account)
		if err != nil {
			return nil, err
		}
		var nextKey []byte
		for {
			resp, err := k.authzKeeper.GranterGrants(ctx, &authz.QueryGranterGrantsRequest{
				Granter: account,
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 1000,
				},
			})
			if err != nil {
				return nil, fmt.Errorf("failed to query authz grants for %s: %w", account, err)
			}
			for _, grant := range resp.Grants {
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

			if resp.Pagination == nil || len(resp.Pagination.NextKey) == 0 {
				break
			}
			nextKey = resp.Pagination.NextKey
		}

		// revoke all feegrant allowances where the blacklisted account is the granter
		if k.feegrantKeeper != nil {
			nextKey = nil
			for {
				resp, err := k.feegrantKeeper.AllowancesByGranter(ctx, &feegrant.QueryAllowancesByGranterRequest{
					Granter: account,
					Pagination: &query.PageRequest{
						Key:   nextKey,
						Limit: 1000,
					},
				})
				if err != nil {
					return nil, fmt.Errorf("failed to query feegrant allowances for %s: %w", account, err)
				}

				for _, allowance := range resp.Allowances {
					if allowance == nil {
						continue
					}
					if _, err := k.feegrantKeeper.RevokeAllowance(ctx, &feegrant.MsgRevokeAllowance{
						Granter: account,
						Grantee: allowance.Grantee,
					}); err != nil {
						return nil, fmt.Errorf("failed to revoke feegrant allowance granter=%s grantee=%s: %w", account, allowance.Grantee, err)
					}
				}

				if resp.Pagination == nil || len(resp.Pagination.NextKey) == 0 {
					break
				}
				nextKey = resp.Pagination.NextKey
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
