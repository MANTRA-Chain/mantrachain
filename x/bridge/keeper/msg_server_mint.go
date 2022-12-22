package keeper

import (
	"context"
	"strings"

	"github.com/LimeChain/mantrachain/x/bridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k msgServer) Mint(goCtx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)

	creator, err := sdk.AccAddressFromBech32(msg.Creator)

	if err != nil {
		return nil, err
	}

	if msg.Mint == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidMint, "mint cannot be empty")
	}

	if msg.Mint.MintList == nil || int32(len(msg.Mint.MintList)) == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidMint, "mint mint list cannot be empty")
	}

	bridgeCreator, err := sdk.AccAddressFromBech32(msg.BridgeCreator)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid bridge creator")
	}

	bridgeController := NewBridgeController(ctx, bridgeCreator).
		WithId(msg.BridgeId).
		WithStore(k)

	err = bridgeController.
		MustExist().
		Validate()

	if err != nil {
		return nil, err
	}

	bridge := bridgeController.getBridge()

	bridgeAccount, err := sdk.AccAddressFromBech32(bridge.BridgeAccount)

	if err != nil {
		return nil, err
	}

	if !creator.Equals(bridgeAccount) {
		return nil, sdkerrors.Wrapf(types.ErrBridgeAccountMismatch, "bridge account %s does not match the creator %s", bridgeAccount.String(), creator.String())
	}

	cw20ContractAddress, err := sdk.AccAddressFromBech32(bridge.Cw20ContractAddress)

	if err != nil {
		return nil, err
	}

	var receivers []string

	wasmExecutor := NewWasmExecutor(ctx, k.wasmViewKeeper, k.wasmContractKeeper)

	if int32(len(msg.Mint.MintList)) > params.ValidMintMintListMetadataMintListMaxCount {
		return nil, sdkerrors.Wrapf(types.ErrInvalidMint, "mint mint list  metadata mint list count invalid %d, max %d", len(msg.Mint.MintList), params.ValidMintMintListMetadataMintListMaxCount)
	}

	for i, mint := range msg.Mint.MintList {
		if strings.TrimSpace(mint.TxHash) == "" {
			return nil, sdkerrors.Wrapf(types.ErrInvalidTxHash, "tx hash should not be empty,d", i)
		}

		_, err = sdk.AccAddressFromBech32(mint.Receiver)

		if err != nil {
			return nil, err
		}

		receiver, err := sdk.AccAddressFromBech32(mint.Receiver)

		if err != nil {
			return nil, err
		}

		txHashType := types.DepositIn
		txHashIndex := types.GetTxHashIndex(mint.TxHash)
		txHash, found := k.GetTxHash(ctx, bridge.Index, txHashIndex)

		if found && txHash.Processed {
			return nil, sdkerrors.Wrapf(types.ErrTxAlreadyProcessed, "tx %s already processed, index %d", mint.TxHash, i)
		}

		err = wasmExecutor.Mint(cw20ContractAddress, creator, receiver, mint.Amount)

		if err != nil {
			return nil, err
		}

		k.SetTxHash(ctx, types.TxHash{
			Index:               txHashIndex,
			BridgeIndex:         bridge.Index,
			TxHash:              mint.TxHash,
			Type:                string(txHashType),
			Processed:           true,
			Receiver:            receiver.String(),
			Amount:              mint.Amount,
			Cw20ContractAddress: bridge.Cw20ContractAddress,
		})

		receivers = append(receivers, receiver.String())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgMint),
			sdk.NewAttribute(types.AttributeKeyBridgeCreator, bridgeCreator.String()),
			sdk.NewAttribute(types.AttributeKeyBridgeId, msg.BridgeId),
			sdk.NewAttribute(types.AttributeKeyCw20ContractAddress, bridge.Cw20ContractAddress),
			sdk.NewAttribute(types.AttributeKeySigner, creator.String()),
			sdk.NewAttribute(types.AttributeKeyReceivers, strings.Join(receivers[:], ",")),
		),
	)

	return &types.MsgMintResponse{
		Creator:             creator.String(),
		Receivers:           receivers,
		BridgeCreator:       bridgeCreator.String(),
		BridgeId:            msg.BridgeId,
		Cw20ContractAddress: bridge.Cw20ContractAddress,
	}, nil
}
