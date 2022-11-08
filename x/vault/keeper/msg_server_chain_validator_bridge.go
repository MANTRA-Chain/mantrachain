package keeper

import (
	"context"

	"github.com/LimeChain/mantrachain/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateChainValidatorBridge(goCtx context.Context, msg *types.MsgCreateChainValidatorBridge) (*types.MsgCreateChainValidatorBridgeResponse, error) {
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

	// Check if the value already exists
	_, isFound := k.GetChainValidatorBridge(
		ctx,
		msg.Chain,
		msg.Validator,
	)
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var chainValidatorBridge = types.ChainValidatorBridge{
		Creator:  msg.Creator,
		BridgeId: msg.BridgeId,
	}

	k.SetChainValidatorBridge(
		ctx,
		msg.Chain,
		msg.Validator,
		chainValidatorBridge,
	)
	return &types.MsgCreateChainValidatorBridgeResponse{
		Chain:     msg.Chain,
		Validator: msg.Validator,
		BridgeId:  msg.BridgeId,
		Creator:   msg.Creator,
	}, nil
}

func (k msgServer) UpdateChainValidatorBridge(goCtx context.Context, msg *types.MsgUpdateChainValidatorBridge) (*types.MsgUpdateChainValidatorBridgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetChainValidatorBridge(
		ctx,
		msg.Chain,
		msg.Validator,
	)
	if !isFound {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "missing bridge %s %s", msg.Chain, msg.Validator)
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != valFound.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	var chainValidatorBridge = types.ChainValidatorBridge{
		Creator:  msg.Creator,
		BridgeId: msg.BridgeId,
		Staked:   valFound.Staked,
	}

	k.SetChainValidatorBridge(
		ctx,
		msg.Chain,
		msg.Validator,
		chainValidatorBridge,
	)

	return &types.MsgUpdateChainValidatorBridgeResponse{
		Chain:     msg.Chain,
		Validator: msg.Validator,
		BridgeId:  msg.BridgeId,
		Creator:   msg.Creator,
	}, nil
}

func (k msgServer) DeleteChainValidatorBridge(goCtx context.Context, msg *types.MsgDeleteChainValidatorBridge) (*types.MsgDeleteChainValidatorBridgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetChainValidatorBridge(
		ctx,
		msg.Chain,
		msg.Validator,
	)
	if !isFound {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "missing bridge %s %s", msg.Chain, msg.Validator)
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != valFound.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveChainValidatorBridge(
		ctx,
		msg.Chain,
		msg.Validator,
	)

	return &types.MsgDeleteChainValidatorBridgeResponse{
		Chain:     msg.Chain,
		Validator: msg.Validator,
		BridgeId:  valFound.BridgeId,
		Creator:   msg.Creator,
	}, nil
}
