package keeper

import (
	"context"
	"strconv"
	"strings"

	"github.com/LimeChain/mantrachain/x/bridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k msgServer) Mint(goCtx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator, err := sdk.AccAddressFromBech32(msg.Creator)

	if err != nil {
		return nil, err
	}

	if msg.Mint == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidMint, "mint cannot be empty")
	}

	if strings.TrimSpace(msg.Mint.TxHash) == "" {
		return nil, sdkerrors.Wrap(types.ErrInvalidTxHash, "tx hash should not be empty")
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

	receiver, err := sdk.AccAddressFromBech32(msg.Mint.Receiver)

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

	// TODO: export as constant
	txHashType := "deposit_in"
	txHashIndex := types.GetTxHashIndex(msg.Mint.TxHash)
	txHash, found := k.GetTxHash(ctx, bridge.Index, txHashIndex)

	if found && txHash.Processed {
		return nil, sdkerrors.Wrapf(types.ErrTxAlreadyProcessed, "tx %s already processed", msg.Mint.TxHash)
	}

	cw20ContractAddress, err := sdk.AccAddressFromBech32(bridge.Cw20ContractAddress)

	if err != nil {
		return nil, err
	}

	wasmExecutor := NewWasmExecutor(ctx, k.wasmViewKeeper, k.wasmContractKeeper)
	err = wasmExecutor.Mint(cw20ContractAddress, creator, receiver, msg.Mint.Amount)

	if err != nil {
		return nil, err
	}

	k.SetTxHash(ctx, types.TxHash{
		Index:       txHashIndex,
		BridgeIndex: bridge.Index,
		TxHash:      msg.Mint.TxHash,
		Type:        txHashType,
		Processed:   true,
	})

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgMint),
			sdk.NewAttribute(types.AttributeKeyBridgeCreator, bridgeCreator.String()),
			sdk.NewAttribute(types.AttributeKeyBridgeId, msg.BridgeId),
			sdk.NewAttribute(types.AttributeKeyCw20ContractAddress, bridge.Cw20ContractAddress),
			sdk.NewAttribute(types.AttributeKeyAmount, strconv.FormatUint(msg.Mint.Amount, 10)),
			sdk.NewAttribute(types.AttributeKeyTxHash, msg.Mint.TxHash),
			sdk.NewAttribute(types.AttributeKeyType, txHashType),
			sdk.NewAttribute(types.AttributeKeySigner, creator.String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, receiver.String()),
		),
	)

	return &types.MsgMintResponse{
		Creator:             creator.String(),
		Receiver:            receiver.String(),
		BridgeCreator:       bridgeCreator.String(),
		BridgeId:            msg.BridgeId,
		Cw20ContractAddress: bridge.Cw20ContractAddress,
		Amount:              msg.Mint.Amount,
		TxHash:              msg.Mint.TxHash,
		Type:                txHashType,
	}, nil
}
