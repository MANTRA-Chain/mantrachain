package keeper

import (
	"context"
	"strings"

	"github.com/LimeChain/mantrachain/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) WithdrawNftReward(goCtx context.Context, msg *types.MsgWithdrawNftReward) (*types.MsgWithdrawNftRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)

	if params.StakingValidatorAddress == "" {
		return nil, sdkerrors.Wrap(types.ErrUnavailable, "staking validator address param not set")
	}

	if params.StakingValidatorDenom == "" {
		return nil, sdkerrors.Wrap(types.ErrUnavailable, "staking validator denom param not set")
	}

	if params.MinRewardWithdrawAmount == 0 {
		return nil, sdkerrors.Wrap(types.ErrUnavailable, "min reward withdraw amount param not set")
	}

	if ctx.ChainID() == "" {
		return nil, sdkerrors.Wrap(types.ErrUnavailable, "chain id not set yet")
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
		WithKeeper(k.Keeper).
		WithConfiguration(params)

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
	var lastEpochWithdrawn int64 = rewardsController.getLastWithdrawnEpochNative()
	var epochs []*types.Epoch
	var rewards sdk.DecCoin = sdk.NewDecCoin(params.StakingValidatorDenom, sdk.Int(sdk.NewDec(0)))
	var sent sdk.Coin = sdk.NewCoin(params.StakingValidatorDenom, sdk.Int(sdk.NewDec(0)))
	var balance sdk.DecCoin = sdk.NewDecCoin(params.StakingValidatorDenom, sdk.Int(sdk.NewDec(0)))

	staked, err := rewardsController.getNativeStaked()

	if err != nil {
		return nil, err
	}

	minEpochRewardsStartBH := rewardsController.getMinEpochRewardsStartBH(staked, lastEpochWithdrawn)

	if minEpochRewardsStartBH != types.UndefinedBlockHeight {
		epochs = k.GetNextRewardsEpochsFromPrevEpochId(
			ctx,
			ctx.ChainID(),
			params.StakingValidatorAddress,
			params.StakingValidatorDenom,
			minEpochRewardsStartBH,
		)
	}

	if len(epochs) > 0 {
		startAt = epochs[0].StartAt
		endAt = epochs[len(epochs)-1].EndAt
		lastEpochWithdrawn = epochs[len(epochs)-1].BlockStart

		rewards = rewardsController.calcNftBalance(epochs, staked, params.StakingValidatorDenom)
	}

	prevBalance := rewardsController.getBalanceCoinNative()
	rewards.Amount = rewards.Amount.Add(prevBalance.Amount)

	if rewards.Amount.GTE(sdk.NewDec(params.MinRewardWithdrawAmount)) {
		coins, remainder := rewards.TruncateDecimal()
		err = k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, sdk.NewCoins(coins))

		if err != nil {
			return nil, err
		}

		sent = coins
		balance = remainder

		rewardsController.setBalanceNative(balance, lastEpochWithdrawn, ctx.BlockHeader().Time.Unix())

		k.SetNftStake(ctx, *rewardsController.getNftStake())
	}

	return &types.MsgWithdrawNftRewardResponse{
		MarketplaceCreator: marketplaceCreator.String(),
		MarketplaceId:      msg.MarketplaceId,
		CollectionCreator:  collectionCreator.String(),
		CollectionId:       msg.CollectionId,
		NftId:              msg.NftId,
		Owner:              owner.String(),
		Receiver:           receiver.String(),
		Balance:            &balance,
		Reward:             &sent,
		StartAt:            startAt,
		EndAt:              endAt,
	}, nil
}
