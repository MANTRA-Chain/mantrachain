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

func (c *RewardsController) getNativeStaked() (staked []*types.NftStakeListItem, err error) {
	if len(c.nftStake.Staked) > 0 {
		for _, v := range c.nftStake.Staked {
			if v.Chain == c.ctx.ChainID() &&
				v.Validator == c.conf.StakingValidatorAddress &&
				v.Denom == c.conf.StakingValidatorDenom {
				staked = append(staked, v)
			}
		}
	}

	return
}

func (c *RewardsController) getLastWithdrawnEpochNative() int64 {
	lastWithdrawnEpoch := types.UndefinedBlockHeight

	if len(c.nftStake.Balances) > 0 {
		for _, v := range c.nftStake.Balances {
			if v.Chain == c.ctx.ChainID() &&
				v.Validator == c.conf.StakingValidatorAddress &&
				v.Denom == c.conf.StakingValidatorDenom {
				lastWithdrawnEpoch = v.LastWithdrawnEpoch

				return lastWithdrawnEpoch
			}
		}
	}

	return lastWithdrawnEpoch
}

func (c *RewardsController) getMinEpochRewardsStartBH(
	staked []*types.NftStakeListItem,
	lastWithdrawnEpoch int64,
) (minEpochRewardsStart int64) {
	minEpochRewardsStart = types.UndefinedBlockHeight

	if len(staked) > 0 {
		for _, v := range staked {
			if minEpochRewardsStart == types.UndefinedBlockHeight {
				if lastWithdrawnEpoch == types.UndefinedBlockHeight {
					minEpochRewardsStart = v.StakedEpoch
				} else {
					minEpochRewardsStart = lastWithdrawnEpoch
				}
			} else {
				if lastWithdrawnEpoch == types.UndefinedBlockHeight {
					if minEpochRewardsStart < v.StakedEpoch {
						minEpochRewardsStart = v.StakedEpoch
					}
				} else {
					if minEpochRewardsStart < lastWithdrawnEpoch {
						minEpochRewardsStart = lastWithdrawnEpoch
					}
				}
			}
		}
	}

	return minEpochRewardsStart
}

func (c *RewardsController) calcNftBalance(epochs []*types.Epoch, staked []*types.NftStakeListItem, stakingValidatorDenom string) (balance sdk.DecCoin) {
	balance = sdk.NewDecCoin(stakingValidatorDenom, sdk.Int(sdk.NewDec(0)))
	for _, epoch := range epochs {
		for _, stake := range staked {
			if stake.StakedEpoch < epoch.BlockStart {
				rewardPerShare := epoch.Rewards.Amount.ToDec().Quo(epoch.Staked)
				balance.Amount = balance.Amount.Add(rewardPerShare.Mul(*stake.Amount))
			}
		}
	}

	return
}

func (c *RewardsController) getBalanceCoinNative() (balance sdk.DecCoin) {
	balance = sdk.NewDecCoin(c.conf.StakingValidatorDenom, sdk.NewInt(0))

	if len(c.nftStake.Balances) > 0 {
		for _, v := range c.nftStake.Balances {
			if v.Chain == c.ctx.ChainID() &&
				v.Validator == c.conf.StakingValidatorAddress &&
				v.Denom == c.conf.StakingValidatorDenom {
				balance = sdk.NewDecCoinFromDec(v.Denom, *v.Amount)

				return balance
			}
		}
	}

	return balance
}

func (c *RewardsController) setBalanceNative(balanceCoin sdk.DecCoin, lastEpochWithdrawn int64, withdrawnAt int64) {
	var balance *types.NftStakeBalance

	if len(c.nftStake.Balances) > 0 {
		for _, v := range c.nftStake.Balances {
			if v.Chain == c.ctx.ChainID() &&
				v.Validator == c.conf.StakingValidatorAddress &&
				v.Denom == c.conf.StakingValidatorDenom {
				balance = v
			}
		}
	}

	if balance != nil {
		balance.Amount = &balanceCoin.Amount
		balance.Denom = balanceCoin.Denom
		balance.LastWithdrawnEpoch = lastEpochWithdrawn
		balance.LastWithdrawnAt = withdrawnAt
	} else {
		c.nftStake.Balances = append(c.nftStake.Balances, &types.NftStakeBalance{
			Chain:              c.ctx.ChainID(),
			Validator:          c.conf.StakingValidatorAddress,
			Denom:              c.conf.StakingValidatorDenom,
			Amount:             &balanceCoin.Amount,
			LastWithdrawnEpoch: lastEpochWithdrawn,
			LastWithdrawnAt:    withdrawnAt,
		})
	}
}
