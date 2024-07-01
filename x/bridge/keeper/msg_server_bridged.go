package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/MANTRA-Finance/mantrachain/x/bridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func (k msgServer) CreateMultiBridged(goCtx context.Context, msg *types.MsgCreateMultiBridged) (*types.MsgCreateMultiBridgedResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator := msg.Input.Address

	if err := k.guardKeeper.CheckHasAuthz(ctx, creator, "module:"+types.ModuleName+":CreateMultiBridged"); err != nil {
		return nil, errors.Wrap(err, "unauthorized")
	}

	for i, ethTxHash := range msg.EthTxHashes {
		if k.HasBridged(ctx, ethTxHash) {
			return nil, errors.Wrapf(errorstypes.ErrInvalidRequest, "ethTxHash already set, ethTxHash %s", ethTxHash)
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

	outputs := make([]banktypes.Output, len(msg.Outputs))
	for i, out := range msg.Outputs {
		outAddress := sdk.MustAccAddressFromBech32(out.Address)
		outputs[i] = banktypes.NewOutput(outAddress, out.Coins)
	}

	err := k.bankKeeper.InputOutputCoins(ctx, banktypes.Input{
		Address: msg.Input.Address,
		Coins:   msg.Input.Coins,
	}, outputs)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateMultiBridgedResponse{}, nil
}
