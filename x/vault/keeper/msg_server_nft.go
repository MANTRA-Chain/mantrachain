package keeper

import (
	"context"
	"strings"

	"github.com/LimeChain/mantrachain/x/vault/types"
	"github.com/LimeChain/mantrachain/x/vault/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

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
	var minTreshold sdk.Coin = sdk.Coin{}
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

	if params.RewardMinClaim != "" {
		minTreshold, err = sdk.ParseCoinNormalized(params.RewardMinClaim)

		if err != nil {
			return nil, err
		}
	}

	if len(rewards) > 0 {
		rewards = utils.SumCoins(rewards, minTreshold)

		if isNativeReward {
			chainValidatorBridge, found := k.GetChainValidatorBridge(ctx, stakingChain, stakingValidator)

			if !found {
				return nil, sdkerrors.Wrap(types.ErrChainValidatorBridgeNotFound, "chain validator bridge not found")
			}

			be := NewBridgeExecutor(ctx, k.bridgeKeeper)
			bridge, found := be.GetBridge(creator, chainValidatorBridge.BridgeId)

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
				err = we.Mint(cw20ContractAddress, k.ac.GetModuleAddress(types.ModuleName), receiver, intReward.Amount.Uint64())
			}

			if err != nil {
				return nil, err
			}

			sent[intReward.Denom] = sent[intReward.Denom].Add(intReward.Amount)
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
