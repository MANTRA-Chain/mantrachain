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

func (c *RewardsController) getStaked(chain string, validator string) (staked []*types.NftStakeListItem, err error) {
	if len(c.nftStake.Staked) > 0 {
		for _, v := range c.nftStake.Staked {
			if v.Chain == chain &&
				v.Validator == validator {
				staked = append(staked, v)
			}
		}
	}

	return
}

func (c *RewardsController) getLastWithdrawnEpoch(chain string, validator string) int64 {
	lastWithdrawnEpoch := types.UndefinedBlockHeight

	if len(c.nftStake.Balances) > 0 {
		for _, v := range c.nftStake.Balances {
			if v.Chain == chain &&
				v.Validator == validator {
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

func (c *RewardsController) calcNftBalances(epochs []*types.Epoch, staked []*types.NftStakeListItem) (balances []*sdk.DecCoin) {
	balances = []*sdk.DecCoin{}
	for _, epoch := range epochs {
		for _, stake := range staked {
			if stake.StakedEpoch < epoch.BlockStart {
				for _, reward := range epoch.Rewards {
					var balance *sdk.DecCoin = nil
					for i := range balances {
						if balances[i].Denom == reward.Denom {
							balance = balances[i]
						}
					}
					if balance == nil {
						coin := sdk.NewDecCoinFromDec(reward.Denom, sdk.ZeroDec())
						balance = &coin
						balances = append(balances, balance)
					}
					rewardPerShare := reward.Amount.ToDec().Quo(epoch.Staked)
					balance.Amount = balance.Amount.Add(rewardPerShare.Mul(sdk.MustNewDecFromStr(stake.Shares)))
				}
			}
		}
	}

	return
}

func (c *RewardsController) getBalancesCoin(chain string, validator string) (balances []*sdk.DecCoin) {
	if len(c.nftStake.Balances) > 0 {
		for _, v := range c.nftStake.Balances {
			if v.Chain == chain &&
				v.Validator == validator {
				coin := sdk.DecCoin{
					Denom:  v.Denom,
					Amount: *v.Amount,
				}
				balances = append(balances, &coin)
			}
		}
	}

	return balances
}

func (c *RewardsController) setBalances(chain string, validator string, balancesCoins []*sdk.DecCoin, lastEpochWithdrawn int64, withdrawnAt int64) {
	var filtered []*types.NftStakeBalance = nil
	for _, v := range c.nftStake.Balances {
		if v.Chain != chain &&
			v.Validator != validator {
			filtered = append(filtered, v)
		}
	}

	c.nftStake.Balances = filtered

	for _, v := range balancesCoins {
		c.nftStake.Balances = append(c.nftStake.Balances, &types.NftStakeBalance{
			Chain:              chain,
			Validator:          validator,
			Denom:              v.Denom,
			Amount:             &v.Amount,
			LastWithdrawnEpoch: lastEpochWithdrawn,
			LastWithdrawnAt:    withdrawnAt,
		})
	}
}
