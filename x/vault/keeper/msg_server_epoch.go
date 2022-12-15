package keeper

import (
	"context"
	"strconv"
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

	if strings.TrimSpace(msg.StakingChain) == "" {
		return nil, sdkerrors.Wrap(types.ErrInvalidStakingChain, "staking chain should not be empty")
	}

	if strings.TrimSpace(msg.StakingValidator) == "" {
		return nil, sdkerrors.Wrap(types.ErrInvalidStakingValidator, "staking validator should not be empty")
	}

	if msg.BlockStart <= 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidBlockStart, "block start should be positive")
	}

	chainValidatorBridge, found := k.GetChainValidatorBridge(ctx, msg.StakingChain, msg.StakingValidator)

	if !found {
		return nil, sdkerrors.Wrap(types.ErrChainValidatorBridgeNotFound, "chain validator bridge not found")
	}

	be := NewBridgeExecutor(ctx, k.bridgeKeeper)
	bridge, found := be.GetBridge(creator, chainValidatorBridge.BridgeId)

	if !found {
		return nil, sdkerrors.Wrapf(types.ErrBridgeDoesNotExist, "bridge not exists")
	}

	bridgeAccount, err := sdk.AccAddressFromBech32(bridge.BridgeAccount)

	if err != nil {
		return nil, err
	}

	if !bridgeAccount.Equals(creator) {
		return nil, sdkerrors.Wrap(types.ErrUnauthorized, "not authorized to start epoch")
	}

	reward := sdk.Coin{}

	if msg.Reward != "" {
		reward, err = sdk.ParseCoinNormalized(msg.Reward)

		if err != nil {
			return nil, err
		}
	}

	lastEpochBlock, found := k.GetLastEpochBlock(ctx, msg.StakingChain, msg.StakingValidator)
	lastEpochBlockHeight := lastEpochBlock.BlockHeight

	if found {
		if lastEpochBlockHeight >= msg.BlockStart {
			return nil, sdkerrors.Wrap(types.ErrInvalidBlockStart, "block start should be greater than the last epoch block")
		}

		lastEpoch, found := k.GetEpoch(ctx, msg.StakingChain, msg.StakingValidator, lastEpochBlockHeight)

		if !found {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "last epoch not found %s", lastEpochBlock)
		}

		cw20ContractAddress, err := sdk.AccAddressFromBech32(bridge.Cw20ContractAddress)

		if err != nil {
			return nil, err
		}

		we := NewWasmExecutor(ctx, k.wasmViewKeeper, k.wasmContractKeeper)

		if !reward.IsNil() && !reward.IsZero() {
			err = we.Mint(cw20ContractAddress, creator, k.ac.GetModuleAddress(types.ModuleName), reward.Amount.String())

			if err != nil {
				return nil, err
			}

			lastEpoch.Rewards = sdk.NewCoins(reward)
		}

		lastEpoch.BlockEnd = msg.BlockStart
		lastEpoch.NextEpochBlock = msg.BlockStart
		lastEpoch.EndAt = ctx.BlockHeader().Time.Unix()

		k.SetEpoch(ctx, msg.StakingChain, msg.StakingValidator, lastEpochBlockHeight, lastEpoch)

		newEpoch := types.Epoch{
			Index:          types.GetEpochIndex(msg.StakingValidator, []byte(strconv.FormatInt(msg.BlockStart, 10))),
			PrevEpochBlock: lastEpochBlockHeight,
			NextEpochBlock: types.UndefinedBlockHeight,
			BlockStart:     msg.BlockStart,
			BlockEnd:       types.UndefinedBlockHeight,
			StartAt:        ctx.BlockHeader().Time.Unix(),
			Staked:         chainValidatorBridge.Staked,
		}

		k.SetEpoch(ctx, msg.StakingChain, msg.StakingValidator, msg.BlockStart, newEpoch)

		k.SetLastEpochBlock(ctx, msg.StakingChain, msg.StakingValidator, types.LastEpochBlock{
			BlockHeight:      msg.BlockStart,
			Creator:          msg.Creator,
			StakingChain:     msg.StakingChain,
			StakingValidator: msg.StakingValidator,
		})
	} else {
		if !reward.IsNil() && !reward.IsZero() {
			return nil, sdkerrors.Wrap(types.ErrInitEpochNotFound, "init epoch not found")
		}

		k.InitEpoch(ctx, msg.StakingChain, msg.StakingValidator, msg.BlockStart)
	}

	return &types.MsgStartEpochResponse{
		PrevEpochBlock:      lastEpochBlockHeight,
		NextEpochBlock:      types.UndefinedBlockHeight,
		BlockStart:          msg.BlockStart,
		BlockEnd:            types.UndefinedBlockHeight,
		StakingChain:        msg.StakingChain,
		StakingValidator:    msg.StakingValidator,
		Staked:              chainValidatorBridge.Staked.String(),
		Cw20ContractAddress: bridge.Cw20ContractAddress,
	}, nil
}
