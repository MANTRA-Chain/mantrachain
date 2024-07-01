package keeper

import (
	"context"

	"github.com/MANTRA-Finance/mantrachain/x/guard/types"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
)

func (k msgServer) UpdateAuthzGenericGrantRevokeBatch(goCtx context.Context, msg *types.MsgUpdateAuthzGenericGrantRevokeBatch) (*types.MsgUpdateAuthzGenericGrantRevokeBatchResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.CheckIsAdmin(ctx, msg.GetCreator()); err != nil {
		return nil, errors.Wrap(err, "unauthorized")
	}

	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	grantee, err := sdk.AccAddressFromBech32(msg.Grantee)
	if err != nil {
		return nil, err
	}

	for _, msg := range msg.AuthzGrantRevokeMsgsTypes.Msgs {
		if k.router.HandlerByTypeURL(msg.TypeUrl) == nil {
			return nil, errorstypes.ErrInvalidType.Wrapf("%s doesn't exist", msg.TypeUrl)
		}

		if msg.Grant {
			authorization := authztypes.NewGenericAuthorization(msg.TypeUrl)

			err = k.authzKeeper.SaveGrant(ctx, grantee, creator, authorization, nil)
			if err != nil {
				return nil, err
			}
		} else {
			err = k.authzKeeper.DeleteGrant(ctx, grantee, creator, msg.TypeUrl)
			if err != nil {
				return nil, err
			}
		}
	}

	return &types.MsgUpdateAuthzGenericGrantRevokeBatchResponse{}, nil
}
