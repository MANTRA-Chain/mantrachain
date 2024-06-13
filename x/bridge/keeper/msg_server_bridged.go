package keeper

import (
	"context"
	"strings"

	"github.com/MANTRA-Finance/mantrachain/x/bridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func (k msgServer) CreateMultiBridged(goCtx context.Context, msg *types.MsgCreateMultiBridged) (*types.MsgCreateMultiBridgedResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if len(msg.Inputs) == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "no inputs")
	}

	if len(msg.Inputs) != 1 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "multiple senders")
	}

	if len(msg.Outputs) == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "no outputs")
	}

	if len(msg.Outputs) != len(msg.EthTxHashes) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "outputs ethTxHashes length mismatch")
	}

	creator := msg.Inputs[0].Address

	if err := k.guardKeeper.CheckHasAuthz(ctx, creator, "module:"+types.ModuleName+":CreateMultiBridged"); err != nil {
		return nil, sdkerrors.Wrap(err, "unauthorized")
	}

	for i, ethTxHash := range msg.EthTxHashes {
		if strings.TrimSpace(ethTxHash) == "" {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "ethTxHash cannot be empty")
		}

		if k.HasBridged(ctx, ethTxHash) {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "ethTxHash already set", "ethTxHash", ethTxHash)
		}

		if len(msg.Outputs[i].Coins) != 1 {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "multiple output coins not supported")
		}

		coin := sdk.NewCoin(msg.Outputs[i].Coins[0].Denom, msg.Outputs[i].Coins[0].Amount)

		bridged := types.Bridged{
			Index: types.BridgedKey(
				ethTxHash,
			),
			EthTxHash: ethTxHash,
			Receiver:  msg.Outputs[i].Address,
			Amount:    &coin,
		}

		k.SetBridged(ctx, bridged)
	}

	inputs := make([]banktypes.Input, len(msg.Inputs))
	for i, in := range msg.Inputs {
		inAddress := sdk.MustAccAddressFromBech32(in.Address)
		inputs[i] = banktypes.NewInput(inAddress, in.Coins)
	}

	outputs := make([]banktypes.Output, len(msg.Outputs))
	for i, out := range msg.Outputs {
		outAddress := sdk.MustAccAddressFromBech32(out.Address)
		outputs[i] = banktypes.NewOutput(outAddress, out.Coins)
	}

	err := k.bankKeeper.InputOutputCoins(ctx, inputs, outputs)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateMultiBridgedResponse{}, nil
}
