package keeper

import (
	"context"

	"github.com/LimeChain/mantrachain/x/guard/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) UpdateGuardTransfer(goCtx context.Context, msg *types.MsgUpdateGuardTransfer) (*types.MsgUpdateGuardTransferResponse, error) {
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
	valFound, isFound := k.GetGuardTransfer(ctx)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "not set")
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != valFound.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	var guardTransfer = types.GuardTransfer{
		Creator: msg.Creator,
		Enabled: msg.Enabled,
	}

	k.SetGuardTransfer(ctx, guardTransfer)

	return &types.MsgUpdateGuardTransferResponse{}, nil
}
