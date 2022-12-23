package keeper

import (
	"context"
	"strconv"
	"strings"

	"github.com/LimeChain/mantrachain/x/vault/types"
	"github.com/LimeChain/mantrachain/x/vault/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) TransferYieldReward(ctx sdk.Context, we *WasmExecutor, cw20ContractAddress sdk.AccAddress, creator sdk.AccAddress, isNativeReward bool, intReward sdk.Coin, receiver sdk.AccAddress) error {
	var err error

	if isNativeReward {
		err = k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, sdk.NewCoins(intReward))
	} else {
		err = we.IncreaseAllowance(cw20ContractAddress, k.ac.GetModuleAddress(types.ModuleName), receiver, intReward.Amount.String())

		if err != nil {
			return err
		}

		err = we.TransferFrom(cw20ContractAddress, creator, k.ac.GetModuleAddress(types.ModuleName), receiver, intReward.Amount.String())
	}

	if err != nil {
		return err
	}

	return nil
}

// TODO: May need min threshold for withdraw yield rewards
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
		return nil, sdkerrors.Wrap(types.ErrInvalidCollectionId, "collection id should not be empty")
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

			bridgeCreator, err := sdk.AccAddressFromBech32(chainValidatorBridge.BridgeCreator)

			if err != nil {
				return nil, err
			}

			be := NewBridgeExecutor(ctx, k.bridgeKeeper)
			bridge, found := be.GetBridge(bridgeCreator, chainValidatorBridge.BridgeId)

			if !found {
				return nil, sdkerrors.Wrapf(types.ErrBridgeDoesNotExist, "bridge not exists")
			}

			cw20ContractAddress, err = sdk.AccAddressFromBech32(bridge.Cw20ContractAddress)

			if err != nil {
				return nil, err
			}
		}

		we := NewWasmExecutor(ctx, k.wasmViewKeeper, k.wasmContractKeeper)

		for _, v := range rewards {
			nftEarningsOnYieldReward := rewardsController.getNftEarningsOnYieldReward()
			reward := sdk.NewDecCoinFromDec(v.Denom, v.Amount)

			if len(nftEarningsOnYieldReward) > 0 {
				for _, j := range nftEarningsOnYieldReward {
					earningAmount := sdk.NewDecFromInt(*j.Percentage).Mul(reward.Amount).Quo(sdk.NewDec(100))
					earningCoin := sdk.NewDecCoinFromDec(reward.Denom, earningAmount)
					intReward, _ := earningCoin.TruncateDecimal()

					err = k.TransferYieldReward(ctx, we, cw20ContractAddress, creator, isNativeReward, intReward, sdk.AccAddress(j.Address))

					if err != nil {
						return nil, err
					}

					reward = reward.Sub(sdk.NewDecCoin(reward.Denom, intReward.Amount))
				}
			}

			intReward, remainder := reward.TruncateDecimal()

			if intReward.IsPositive() {
				err = k.TransferYieldReward(ctx, we, cw20ContractAddress, creator, isNativeReward, intReward, sdk.AccAddress(receiver))

				if err != nil {
					return nil, err
				}
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

		rewardsController.setBalances(
			stakingChain,
			stakingValidator,
			remainBalances,
			lastEpochWithdrawn,
			ctx.BlockHeader().Time.Unix(),
		)
		rewardsController.setInitiallyRewardWithdrawn(stakingChain, stakingValidator, true)

		k.SetNftStake(ctx, *rewardsController.getNftStake())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgWithdrawNftRewards),
			sdk.NewAttribute(types.AttributeKeyMarketplaceCreator, marketplaceCreator.String()),
			sdk.NewAttribute(types.AttributeKeyMarketplaceId, msg.MarketplaceId),
			sdk.NewAttribute(types.AttributeKeyCollectionCreator, collectionCreator.String()),
			sdk.NewAttribute(types.AttributeKeyCollectionId, msg.CollectionId),
			sdk.NewAttribute(types.AttributeKeyNftId, msg.NftId),
			sdk.NewAttribute(types.AttributeKeyReceiver, receiver.String()),
			sdk.NewAttribute(types.AttributeKeyStartAt, strconv.FormatInt(startAt, 10)),
			sdk.NewAttribute(types.AttributeKeyEndAt, strconv.FormatInt(endAt, 10)),
			sdk.NewAttribute(types.AttributeKeySigner, creator.String()),
			sdk.NewAttribute(types.AttributeKeyOwner, creator.String()),
		),
	)

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

func (k msgServer) UpdateNftStakeStaked(goCtx context.Context, msg *types.MsgUpdateNftStakeStaked) (*types.MsgUpdateNftStakeStakedResponse, error) {
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
		return nil, sdkerrors.Wrap(types.ErrInvalidCollectionId, "collection id should not be empty")
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

	chainValidatorBridge, isFound := k.GetChainValidatorBridge(
		ctx,
		msg.StakingChain,
		msg.StakingValidator,
	)
	if !isFound {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "missing bridge %s %s", msg.StakingChain, msg.StakingValidator)
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
		return nil, sdkerrors.Wrap(types.ErrUnauthorized, "not authorized to set staked")
	}

	nftStake := rewardsController.getNftStake()

	lastEpochBlock, found := k.GetLastEpochBlock(ctx, msg.StakingChain, msg.StakingValidator)

	if !found {
		return nil, sdkerrors.Wrap(types.ErrLastEpochBlockNotFound, "last epoch block not found")
	}

	shares, err := sdk.NewDecFromStr(msg.Shares)

	if err != nil {
		return nil, err
	}

	if nftStake.Staked[msg.StakedIndex] == nil {
		return nil, sdkerrors.Wrap(types.ErrNftStakeStakedNotFound, "nft stake staked not found")
	}

	if nftStake.Staked[msg.StakedIndex].Chain != msg.StakingChain ||
		nftStake.Staked[msg.StakedIndex].Validator != msg.StakingValidator {
		return nil, sdkerrors.Wrap(types.ErrNftStakeStakedChainValidatorNotMatch, "nft stake staked chain validator not match")
	}

	if nftStake.Staked[msg.StakedIndex].StakedAt != 0 {
		return nil, sdkerrors.Wrap(types.ErrNftStakeStakedAlreadyBeingSet, "nft stake staked at not zero")
	}

	nftStake.Staked[msg.StakedIndex].StakedAt = ctx.BlockHeader().Time.Unix()
	nftStake.Staked[msg.StakedIndex].StakedEpoch = lastEpochBlock.BlockHeight
	nftStake.Staked[msg.StakedIndex].BlockHeight = msg.BlockHeight
	nftStake.Staked[msg.StakedIndex].Shares = msg.Shares

	k.SetNftStake(ctx, *rewardsController.getNftStake())

	chainValidatorBridge.Staked = chainValidatorBridge.Staked.Add(shares)

	k.SetChainValidatorBridge(ctx, msg.StakingChain, msg.StakingValidator, chainValidatorBridge)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgUpdateNftStakeStaked),
			sdk.NewAttribute(types.AttributeKeyMarketplaceCreator, marketplaceCreator.String()),
			sdk.NewAttribute(types.AttributeKeyMarketplaceId, msg.MarketplaceId),
			sdk.NewAttribute(types.AttributeKeyCollectionCreator, collectionCreator.String()),
			sdk.NewAttribute(types.AttributeKeyCollectionId, msg.CollectionId),
			sdk.NewAttribute(types.AttributeKeyNftId, msg.NftId),
			sdk.NewAttribute(types.AttributeKeyStakingChain, msg.StakingChain),
			sdk.NewAttribute(types.AttributeKeyStakingValidator, msg.StakingValidator),
			sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(msg.BlockHeight, 10)),
			sdk.NewAttribute(types.AttributeKeyStakedIndex, strconv.FormatInt(msg.StakedIndex, 10)),
			sdk.NewAttribute(types.AttributeKeyShares, msg.Shares),
		),
	)

	return &types.MsgUpdateNftStakeStakedResponse{
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
