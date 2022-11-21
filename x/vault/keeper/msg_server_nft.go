package keeper

import (
	"context"
	"strings"

	"github.com/LimeChain/mantrachain/x/vault/types"
	"github.com/LimeChain/mantrachain/x/vault/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// TODO: Add min threshold for withdraw yield rewards
func (k msgServer) WithdrawNftRewards(goCtx context.Context, msg *types.MsgWithdrawNftRewards) (*types.MsgWithdrawNftRewardsResponse, error) {
	var stakingChain = ""
	var stakingValidator = ""
	var isNativeReward = false

	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)

	if strings.TrimSpace(msg.StakingChain) != "" {
		stakingChain = msg.StakingChain
		stakingValidator = msg.StakingValidator
	} else {
		isNativeReward = true
		stakingChain = ctx.ChainID()
		stakingValidator = params.StakingValidatorAddress
	}

	if strings.TrimSpace(stakingValidator) == "" {
		return nil, sdkerrors.Wrap(types.ErrUnavailable, "staking validator address param not set")
	}

	creator, err := sdk.AccAddressFromBech32(msg.Creator)

	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(msg.Receiver) == "" {
		msg.Receiver = msg.Creator
	}

	receiver, err := sdk.AccAddressFromBech32(msg.Receiver)

	if err != nil {
		return nil, err
	}

	marketplaceCreator, err := sdk.AccAddressFromBech32(msg.MarketplaceCreator)

	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(msg.MarketplaceId) == "" {
		return nil, sdkerrors.Wrap(types.ErrInvalidMarketplaceId, "marketplace id should not be empty")
	}

	collectionCreator, err := sdk.AccAddressFromBech32(msg.CollectionCreator)

	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(msg.CollectionId) == "" {
		return nil, sdkerrors.Wrap(types.ErrInvalidCollectionId, "marketplace id should not be empty")
	}

	if strings.TrimSpace(msg.NftId) == "" {
		return nil, sdkerrors.Wrap(types.ErrInvalidNftId, "nft id should not be empty")
	}

	rewardsController := NewRewardsController(ctx, marketplaceCreator, msg.MarketplaceId, collectionCreator, msg.CollectionId).
		WithNftId(msg.NftId).
		WithKeeper(k.Keeper)

	err = rewardsController.
		NftStakeMustExist().
		Validate()

	if err != nil {
		return nil, err
	}

	collectionIndex := rewardsController.getCollectionIndex()
	index := rewardsController.getNftIndex()

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)
	owner := nftExecutor.GetNftOwner(string(collectionIndex), string(index))

	if owner == nil || owner.Empty() {
		return nil, sdkerrors.Wrap(types.ErrUnavailable, "nft owner not found")
	}

	if !owner.Equals(creator) {
		return nil, sdkerrors.Wrap(types.ErrUnauthorized, "nft is not owned by the rewards recipient")
	}

	var startAt int64
	var endAt int64
	var lastEpochWithdrawn int64 = rewardsController.getLastWithdrawnEpoch(stakingChain, stakingValidator)
	var epochs []*types.Epoch
	var rewards []*sdk.DecCoin = nil
	var intRewards []*sdk.Coin = nil
	var remainBalances []*sdk.DecCoin = nil
	var sent = make(map[string]sdk.Int)
	var balances = make(map[string]sdk.Dec)
	var cw20ContractAddress sdk.AccAddress

	staked, err := rewardsController.getStaked(stakingChain, stakingValidator)

	if err != nil {
		return nil, err
	}

	minEpochRewardsStartBH := rewardsController.getMinEpochRewardsStartBH(staked, lastEpochWithdrawn)

	if minEpochRewardsStartBH != types.UndefinedBlockHeight {
		epochs = k.GetNextRewardsEpochsFromPrevEpochId(
			ctx,
			stakingChain,
			stakingValidator,
			minEpochRewardsStartBH,
		)
	}

	if len(epochs) > 0 {
		startAt = epochs[0].StartAt
		endAt = epochs[len(epochs)-1].EndAt
		lastEpochWithdrawn = epochs[len(epochs)-1].BlockStart

		rewards = rewardsController.calcNftBalances(epochs, staked)
	}

	prevBalances := rewardsController.getBalancesCoin(stakingChain, stakingValidator)
	rewards = append(rewards, prevBalances...)

	if len(rewards) > 0 {
		rewards = utils.SumCoins(rewards)

		if !isNativeReward {
			chainValidatorBridge, found := k.GetChainValidatorBridge(ctx, stakingChain, stakingValidator)

			if !found {
				return nil, sdkerrors.Wrap(types.ErrChainValidatorBridgeNotFound, "chain validator bridge not found")
			}

			bridgeAccount, err := sdk.AccAddressFromBech32(chainValidatorBridge.BridgeAccount)

			if err != nil {
				return nil, err
			}

			be := NewBridgeExecutor(ctx, k.bridgeKeeper)
			bridge, found := be.GetBridge(bridgeAccount, chainValidatorBridge.BridgeId)

			if !found {
				return nil, sdkerrors.Wrapf(types.ErrBridgeDoesNotExist, "bridge not exists")
			}

			cw20ContractAddress, err = sdk.AccAddressFromBech32(bridge.Cw20ContractAddress)

			if err != nil {
				return nil, err
			}
		}

		for _, reward := range rewards {
			intReward, remainder := reward.TruncateDecimal()

			if isNativeReward {
				err = k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, sdk.NewCoins(intReward))
			} else {
				we := NewWasmExecutor(ctx, k.wasmViewKeeper, k.wasmContractKeeper)
				err = we.IncreaseAllowance(cw20ContractAddress, k.ac.GetModuleAddress(types.ModuleName), receiver, intReward.Amount.Uint64())

				if err != nil {
					return nil, err
				}

				err = we.TransferFrom(cw20ContractAddress, creator, k.ac.GetModuleAddress(types.ModuleName), receiver, intReward.Amount.Uint64())
			}

			if err != nil {
				return nil, err
			}

			if sent[intReward.Denom].IsNil() {
				sent[intReward.Denom] = sdk.NewInt(0)
			}
			sent[intReward.Denom] = sent[intReward.Denom].Add(intReward.Amount)

			if balances[remainder.Denom].IsNil() {
				balances[remainder.Denom] = sdk.NewDec(0)
			}
			balances[remainder.Denom] = balances[remainder.Denom].Add(remainder.Amount)
		}

		for denom, amount := range sent {
			intRewards = append(intRewards, &sdk.Coin{
				Denom:  denom,
				Amount: amount,
			})
		}

		for denom, amount := range balances {
			remainBalances = append(remainBalances, &sdk.DecCoin{
				Denom:  denom,
				Amount: amount,
			})
		}

		rewardsController.setBalances(stakingChain, stakingValidator, remainBalances, lastEpochWithdrawn, ctx.BlockHeader().Time.Unix())
		k.SetNftStake(ctx, *rewardsController.getNftStake())
	}

	return &types.MsgWithdrawNftRewardsResponse{
		MarketplaceCreator: marketplaceCreator.String(),
		MarketplaceId:      msg.MarketplaceId,
		CollectionCreator:  collectionCreator.String(),
		CollectionId:       msg.CollectionId,
		NftId:              msg.NftId,
		Owner:              owner.String(),
		Receiver:           receiver.String(),
		Balances:           remainBalances,
		Rewards:            intRewards,
		StartAt:            startAt,
		EndAt:              endAt,
	}, nil
}

func (k msgServer) SetStaked(goCtx context.Context, msg *types.MsgSetStaked) (*types.MsgSetStakedResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if strings.TrimSpace(msg.StakingChain) == "" {
		return nil, sdkerrors.Wrap(types.ErrUnavailable, "staking chain should not be empty")
	}

	if strings.TrimSpace(msg.StakingValidator) == "" {
		return nil, sdkerrors.Wrap(types.ErrUnavailable, "staking validator address should not be empty")
	}

	creator, err := sdk.AccAddressFromBech32(msg.Creator)

	if err != nil {
		return nil, err
	}

	marketplaceCreator, err := sdk.AccAddressFromBech32(msg.MarketplaceCreator)

	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(msg.MarketplaceId) == "" {
		return nil, sdkerrors.Wrap(types.ErrInvalidMarketplaceId, "marketplace id should not be empty")
	}

	collectionCreator, err := sdk.AccAddressFromBech32(msg.CollectionCreator)

	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(msg.CollectionId) == "" {
		return nil, sdkerrors.Wrap(types.ErrInvalidCollectionId, "marketplace id should not be empty")
	}

	if strings.TrimSpace(msg.NftId) == "" {
		return nil, sdkerrors.Wrap(types.ErrInvalidNftId, "nft id should not be empty")
	}

	if msg.BlockHeight <= 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidCollectionId, "block height should be positive")
	}

	rewardsController := NewRewardsController(ctx, marketplaceCreator, msg.MarketplaceId, collectionCreator, msg.CollectionId).
		WithNftId(msg.NftId).
		WithKeeper(k.Keeper)

	err = rewardsController.
		NftStakeMustExist().
		Validate()

	if err != nil {
		return nil, err
	}

	valFound, isFound := k.GetChainValidatorBridge(
		ctx,
		msg.StakingChain,
		msg.StakingValidator,
	)
	if !isFound {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "missing bridge %s %s", msg.StakingChain, msg.StakingValidator)
	}

	be := NewBridgeExecutor(ctx, k.bridgeKeeper)
	bridge, found := be.GetBridge(creator, valFound.BridgeId)

	if !found {
		return nil, sdkerrors.Wrapf(types.ErrBridgeDoesNotExist, "bridge not exists")
	}

	bridgeAccount, err := sdk.AccAddressFromBech32(bridge.BridgeAccount)

	if err != nil {
		return nil, err
	}

	if !bridgeAccount.Equals(creator) {
		return nil, sdkerrors.Wrap(types.ErrUnauthorized, "not authorized to set staked")
	}

	nftStake := rewardsController.getNftStake()

	lastEpochBlock, found := k.GetLastEpochBlock(ctx, msg.StakingChain, msg.StakingValidator)

	if !found {
		return nil, sdkerrors.Wrap(types.ErrLastEpochBlockNotFound, "last epoch block not found")
	}

	// TODO: Validate if msg.Shares is valid decimal

	if nftStake.Staked[msg.StakedIndex].Chain == msg.StakingChain &&
		nftStake.Staked[msg.StakedIndex].Validator == msg.StakingValidator &&
		// Do not update if StakedAt is already set
		nftStake.Staked[msg.StakedIndex].StakedAt == 0 {
		nftStake.Staked[msg.StakedIndex].StakedAt = ctx.BlockHeader().Time.Unix()
		nftStake.Staked[msg.StakedIndex].StakedEpoch = lastEpochBlock.BlockHeight
		nftStake.Staked[msg.StakedIndex].BlockHeight = msg.BlockHeight
		nftStake.Staked[msg.StakedIndex].Shares = msg.Shares
	} else {
		return nil, sdkerrors.Wrap(types.ErrNftStakeStakedNotFound, "nft stake staked not found")
	}

	k.SetNftStake(ctx, *rewardsController.getNftStake())

	return &types.MsgSetStakedResponse{
		MarketplaceCreator: marketplaceCreator.String(),
		MarketplaceId:      msg.MarketplaceId,
		CollectionCreator:  collectionCreator.String(),
		CollectionId:       msg.CollectionId,
		NftId:              msg.NftId,
		StakingChain:       msg.StakingChain,
		StakingValidator:   msg.StakingValidator,
		BlockHeight:        msg.BlockHeight,
		StakedIndex:        msg.StakedIndex,
		Shares:             msg.Shares,
	}, nil
}
