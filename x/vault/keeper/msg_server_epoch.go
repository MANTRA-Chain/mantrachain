package keeper

import (
	"context"
	"strings"

	"github.com/LimeChain/mantrachain/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) StartEpoch(goCtx context.Context, msg *types.MsgStartEpoch) (*types.MsgStartEpochResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator, err := sdk.AccAddressFromBech32(msg.Creator)

	if err != nil {
		return nil, err
	}

	if msg.BlockStart <= 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidBlockStart, "block start should be positive")
	}

	if strings.TrimSpace(msg.Chain) == "" {
		return nil, sdkerrors.Wrap(types.ErrInvalidChain, "chain should not be empty")
	}

	if strings.TrimSpace(msg.Validator) == "" {
		return nil, sdkerrors.Wrap(types.ErrInvalidValidator, "validator should not be empty")
	}

	reward, err := sdk.ParseCoinNormalized(msg.Reward)

	if err != nil {
		return nil, err
	}

	chainValidatorBridge, found := k.GetChainValidatorBridge(ctx, msg.Chain, msg.Validator)

	if !found {
		return nil, sdkerrors.Wrap(types.ErrChainValidatorBridgeNotFound, "chain validator bridge not found")
	}

	be := NewBridgeExecutor(ctx, k.bridgeKeeper)
	bridge, found := be.GetBridge(creator, chainValidatorBridge.BridgeId)

	if !found {
		return nil, sdkerrors.Wrapf(types.ErrBridgeDoesNotExist, "bridge not exists")
	}

	lastEpochBlock, found := k.GetLastEpochBlock(ctx, msg.Chain, msg.Validator)
	lastEpochBlockHeight := int64(0)

	if !found {
		k.InitEpoch(ctx, msg.Chain, msg.Validator, msg.BlockStart)
	} else {
		lastEpochBlockHeight = lastEpochBlock.BlockHeight
		lastEpoch, found := k.GetEpoch(ctx, msg.Chain, msg.Validator, lastEpochBlockHeight)

		if !found {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "last epoch not found %s", lastEpochBlock)
		}

		if lastEpochBlockHeight >= msg.BlockStart {
			return nil, sdkerrors.Wrap(types.ErrInvalidBlockStart, "block start should be greater than the last epoch block")
		}

		cw20ContractAddress, err := sdk.AccAddressFromBech32(bridge.Cw20ContractAddress)

		if err != nil {
			return nil, err
		}

		we := NewWasmExecutor(ctx, k.wasmViewKeeper, k.wasmContractKeeper)
		err = we.Mint(cw20ContractAddress, creator, k.ac.GetModuleAddress(types.ModuleName), reward.Amount.Uint64())

		if err != nil {
			return nil, err
		}

		lastEpoch.Rewards = sdk.NewCoins(reward)
		lastEpoch.BlockEnd = msg.BlockStart
		lastEpoch.NextEpochBlock = msg.BlockStart
		lastEpoch.EndAt = ctx.BlockHeader().Time.Unix()

		k.SetEpoch(ctx, msg.Chain, msg.Validator, lastEpochBlockHeight, lastEpoch)

		newEpoch := types.Epoch{
			PrevEpochBlock: lastEpochBlockHeight,
			NextEpochBlock: types.UndefinedBlockHeight,
			BlockStart:     msg.BlockStart,
			BlockEnd:       types.UndefinedBlockHeight,
			StartAt:        ctx.BlockHeader().Time.Unix(),
			Staked:         chainValidatorBridge.Staked,
		}

		k.SetEpoch(ctx, msg.Chain, msg.Validator, msg.BlockStart, newEpoch)
	}

	k.SetLastEpochBlock(ctx, msg.Chain, msg.Validator, types.LastEpochBlock{
		BlockHeight: msg.BlockStart,
	})

	return &types.MsgStartEpochResponse{
		PrevEpochBlock:      lastEpochBlockHeight,
		NextEpochBlock:      types.UndefinedBlockHeight,
		BlockStart:          msg.BlockStart,
		BlockEnd:            types.UndefinedBlockHeight,
		Chain:               msg.Chain,
		Validator:           msg.Validator,
		PrevEpochRewards:    []*sdk.Coin{&reward},
		Staked:              chainValidatorBridge.Staked.String(),
		Cw20ContractAddress: bridge.Cw20ContractAddress,
	}, nil
}
