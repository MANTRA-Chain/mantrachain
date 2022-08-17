package keeper

import (
	"github.com/LimeChain/mantrachain/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type RewardsControllerFunc func(controller *RewardsController) error

type RewardsController struct {
	validators         []RewardsControllerFunc
	nftStake           *types.NftStake
	keeper             Keeper
	conf               *types.Params
	ctx                sdk.Context
	marketplaceCreator sdk.AccAddress
	collectionCreator  sdk.AccAddress
	marketplaceId      string
	collectionId       string
	marketplaceIndex   []byte
	collectionIndex    []byte
	id                 string
	index              []byte
}

func NewRewardsController(
	ctx sdk.Context,
	marketplaceCreator sdk.AccAddress,
	marketplaceId string,
	collectionCreator sdk.AccAddress,
	collectionId string,
) *RewardsController {
	marketplaceIndex := NewMarketplaceResolver().GetMarketplaceIndex(marketplaceCreator, marketplaceId)
	collectionIndex := NewTokenResolver().GetCollectionIndex(collectionCreator, collectionId)

	return &RewardsController{
		ctx:                ctx,
		marketplaceCreator: marketplaceCreator,
		marketplaceId:      marketplaceId,
		collectionCreator:  collectionCreator,
		collectionId:       collectionId,
		marketplaceIndex:   marketplaceIndex,
		collectionIndex:    collectionIndex,
	}
}

func (c *RewardsController) WithNftId(id string) *RewardsController {
	c.id = id
	index := NewTokenResolver().GetNftIndex(c.collectionCreator, c.collectionId, c.getId())
	c.index = index
	return c
}

func (c *RewardsController) WithKeeper(keeper Keeper) *RewardsController {
	c.keeper = keeper
	return c
}

func (c *RewardsController) WithNftStake(nftStake types.NftStake) *RewardsController {
	c.nftStake = &nftStake
	return c
}

func (c *RewardsController) WithConfiguration(cfg types.Params) *RewardsController {
	c.conf = &cfg
	return c
}

func (c *RewardsController) Validate() error {
	for _, check := range c.validators {
		if err := check(c); err != nil {
			return err
		}
	}
	return nil
}

func (c *RewardsController) NftStakeMustExist() *RewardsController {
	c.validators = append(c.validators, func(controller *RewardsController) error {
		return controller.nftStakeMustExist()
	})
	return c
}

func (c *RewardsController) nftStakeMustExist() error {
	return c.requireNftStake()
}

func (c *RewardsController) requireNftStake() error {
	if c.nftStake != nil {
		return nil
	}

	nftStake, isFound := c.keeper.GetNftStake(
		c.ctx,
		c.marketplaceIndex,
		c.collectionIndex,
		c.getIndex(),
	)

	if !isFound {
		return sdkerrors.Wrapf(types.ErrNftStakeDoesNotExist, "not found: %s", c.getId())
	}
	c.nftStake = &nftStake
	return nil
}

func (c *RewardsController) getNftStake() *types.NftStake {
	return c.nftStake
}

func (c *RewardsController) getId() string {
	return c.id
}

func (c *RewardsController) getCollectionIndex() []byte {
	return c.collectionIndex
}

func (c *RewardsController) getNftIndex() []byte {
	return c.index
}

func (c *RewardsController) getIndex() []byte {
	return NewTokenResolver().GetNftIndex(c.collectionCreator, c.collectionId, c.getId())
}

func (c *RewardsController) getNativeStaked() (staked []*types.Stake, err error) {
	if len(c.nftStake.Staked) > 0 {
		for _, v := range c.nftStake.Staked {
			if v.Chain == c.ctx.ChainID() && v.Validator == c.conf.StakingValidatorAddress {
				parsed, err := sdk.ParseCoinNormalized(v.Amount)

				if err != nil {
					return nil, err
				}

				if parsed.Denom == c.conf.StakingValidatorDenom {
					staked = append(staked, &v)
				}
			}
		}
	}

	return
}

func (c *RewardsController) getMinEpochRewardsStartBH(staked []*types.Stake) (minEpochRewardsStart int64) {
	minEpochRewardsStart = types.UndefinedBlockHeight

	if len(staked) > 0 {
		for _, v := range staked {
			if minEpochRewardsStart == types.UndefinedBlockHeight {
				if v.LastEpochWithdrawn == types.UndefinedBlockHeight {
					minEpochRewardsStart = v.Epoch
				} else {
					minEpochRewardsStart = v.LastEpochWithdrawn
				}
			} else {
				if v.LastEpochWithdrawn == types.UndefinedBlockHeight {
					if minEpochRewardsStart < v.Epoch {
						minEpochRewardsStart = v.Epoch
					}
				} else {
					if minEpochRewardsStart < v.LastEpochWithdrawn {
						minEpochRewardsStart = v.LastEpochWithdrawn
					}
				}
			}
		}
	}

	return minEpochRewardsStart
}

func (c *RewardsController) calcNftBalance(epochs []*types.Epoch, staked []*types.Stake, stakingValidatorDenom string) (balance sdk.DecCoin) {
	balance = sdk.NewDecCoin(stakingValidatorDenom, sdk.Int(sdk.NewDec(0)))
	for _, epoch := range epochs {
		for _, stake := range staked {
			if stake.Epoch < epoch.BlockStart {
				rewardPerShare := epoch.Rewards.Amount.ToDec().Quo(epoch.Staked)
				if s, err := sdk.ParseCoinNormalized(stake.Amount); err == nil {
					balance.Amount = balance.Amount.Add(rewardPerShare.Mul(sdk.NewDec(s.Amount.Int64())))
				}
			}
		}
	}

	return
}

func (c *RewardsController) getPrevBalanceNative() (balance sdk.DecCoin) {
	balance = sdk.NewDecCoin(c.conf.StakingValidatorDenom, sdk.NewInt(0))
	if len(c.nftStake.Balances) > 0 {
		for _, v := range c.nftStake.Balances {
			if v.Denom == c.conf.StakingValidatorDenom {
				balance.Amount = balance.Amount.Add(v.Amount)
			}
		}
	}

	return
}

func (c *RewardsController) setNextBalanceNative(balance sdk.DecCoin) {
	balances := sdk.DecCoins{balance}
	if len(c.nftStake.Balances) > 0 {
		for _, v := range c.nftStake.Balances {
			if v.Denom != c.conf.StakingValidatorDenom {
				balances = append(balances, v)
			}
		}
	}

	c.nftStake.Balances = balances
	c.keeper.SetNftStake(c.ctx, *c.getNftStake())
}
