package keeper

import (
	"context"
	"strings"

	"github.com/LimeChain/mantrachain/x/bridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) RegisterBridge(goCtx context.Context, msg *types.MsgRegisterBridge) (*types.MsgRegisterBridgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator, err := sdk.AccAddressFromBech32(msg.Creator)

	if err != nil {
		return nil, err
	}

	bridgeAccount, err := sdk.AccAddressFromBech32(msg.Bridge.BridgeAccount)

	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(msg.Bridge.Id) == "" {
		return nil, sdkerrors.Wrap(types.ErrInvalidBridgeId, "bridge id should not be empty")
	}

	bridgeController := NewBridgeController(ctx, creator).
		WithMetadata(msg.Bridge).
		WithStore(k).
		WithConfiguration(k.GetParams(ctx))

	err = bridgeController.
		MustNotExist().
		ValidMetadata().
		Validate()

	if err != nil {
		return nil, err
	}

	var cw20ContractAddress sdk.AccAddress
	wasmExecutor := NewWasmExecutor(ctx, k.wasmViewKeeper, k.wasmContractKeeper)

	if strings.TrimSpace(msg.Bridge.Cw20ContractAddress) == "" {
		cw20Contract, isFound := k.GetCw20Contract(ctx)
		if !isFound {
			return nil, sdkerrors.Wrapf(types.ErrInvalidCw20ContractAddress, "cw20 contract address is invalid")
		}

		cw20ContractAddress, err = wasmExecutor.Instantiate(
			cw20Contract.CodeId,
			creator,
			bridgeAccount,
			msg.Bridge.Cw20Name,
			msg.Bridge.Cw20Symbol,
			msg.Bridge.Cw20Decimals,
			msg.Bridge.Cw20InitialBalances,
			msg.Bridge.Cw20Mint,
		)

		if err != nil {
			return nil, err
		}
	} else {
		cw20ContractAddress, err = sdk.AccAddressFromBech32(msg.Bridge.Cw20ContractAddress)

		if err != nil {
			return nil, err
		}

		minter, err := wasmExecutor.GetMinter(cw20ContractAddress)

		if err != nil {
			return nil, err
		}

		if !minter.Equals(bridgeAccount) {
			return nil, sdkerrors.Wrapf(types.ErrBridgeAccountMismatch, "bridge account %s does not match the minter %s", bridgeAccount.String(), minter.String())
		}
	}

	bridgeIndex := bridgeController.getIndex()
	bridgeId := bridgeController.getId()

	newBridge := types.Bridge{
		Index:               bridgeIndex,
		Id:                  bridgeId,
		Cw20ContractAddress: cw20ContractAddress.String(),
		BridgeAccount:       bridgeAccount.String(),
		Creator:             creator,
		Owner:               creator,
	}

	k.SetBridge(ctx, newBridge)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgRegisterBridge),
			sdk.NewAttribute(types.AttributeKeyBridgeId, bridgeId),
			sdk.NewAttribute(types.AttributeKeyBridgeAccount, bridgeAccount.String()),
			sdk.NewAttribute(types.AttributeKeySigner, creator.String()),
			sdk.NewAttribute(types.AttributeKeyOwner, creator.String()),
		),
	)

	return &types.MsgRegisterBridgeResponse{
		BridgeId:            bridgeId,
		Cw20ContractAddress: msg.Bridge.Cw20ContractAddress,
		BridgeCreator:       creator.String(),
	}, nil
}
