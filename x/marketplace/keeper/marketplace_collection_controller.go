package keeper

import (
	"github.com/LimeChain/mantrachain/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type MarketplaceCollectionControllerFunc func(controller *MarketplaceCollectionController) error

type MarketplaceCollectionController struct {
	validators            []MarketplaceCollectionControllerFunc
	collection            *types.MsgMarketplaceCollection
	marketplaceCollection *types.MarketplaceCollection
	marketplaceIndex      []byte
	collectionIndex       []byte
	marketplaceId         string
	collectionId          string
	store                 msgServer
	conf                  *types.Params
	ctx                   sdk.Context
}

func NewMarketplaceCollectionController(ctx sdk.Context, marketplaceIndex []byte, collectionIndex []byte, marketplaceId string, collectionId string) *MarketplaceCollectionController {
	return &MarketplaceCollectionController{
		marketplaceIndex: marketplaceIndex,
		collectionIndex:  collectionIndex,
		marketplaceId:    marketplaceId,
		collectionId:     collectionId,
		ctx:              ctx,
	}
}

func (c *MarketplaceCollectionController) WithCollection(collection *types.MsgMarketplaceCollection) *MarketplaceCollectionController {
	c.collection = collection
	return c
}

func (c *MarketplaceCollectionController) WithStore(store msgServer) *MarketplaceCollectionController {
	c.store = store
	return c
}

func (c *MarketplaceCollectionController) WithConfiguration(cfg types.Params) *MarketplaceCollectionController {
	c.conf = &cfg
	return c
}

func (c *MarketplaceCollectionController) Validate() error {
	for _, check := range c.validators {
		if err := check(c); err != nil {
			return err
		}
	}
	return nil
}

func (c *MarketplaceCollectionController) MustNotExist() *MarketplaceCollectionController {
	c.validators = append(c.validators, func(controller *MarketplaceCollectionController) error {
		return controller.mustNotExist()
	})
	return c
}

func (c *MarketplaceCollectionController) MustExist() *MarketplaceCollectionController {
	c.validators = append(c.validators, func(controller *MarketplaceCollectionController) error {
		return controller.mustExist()
	})
	return c
}

func (c *MarketplaceCollectionController) mustExist() error {
	return c.requireMarketplaceCollection()
}

func (c *MarketplaceCollectionController) mustNotExist() error {
	err := c.requireMarketplaceCollection()
	if err == nil {
		return sdkerrors.Wrapf(types.ErrMarketplaceCollectionAlreadyExists, "already exists: marketplace id %s, collection id %s", c.marketplaceId, c.collectionId)
	}
	return nil
}

func (c *MarketplaceCollectionController) requireMarketplaceCollection() error {
	if c.marketplaceCollection != nil {
		return nil
	}
	marketplaceCollection, isFound := c.store.GetMarketplaceCollection(c.ctx, c.marketplaceIndex, c.getIndex())
	if !isFound {
		return sdkerrors.Wrapf(types.ErrMarketplaceCollectionDoesNotExist, "not found: marketplace id %s, collection id %s", c.marketplaceId, c.collectionId)
	}
	c.marketplaceCollection = &marketplaceCollection
	return nil
}

func (c *MarketplaceCollectionController) ValidCollection() *MarketplaceCollectionController {
	c.validators = append(c.validators, func(controller *MarketplaceCollectionController) error {
		return controller.validInitiallyNftMinPrice()
	}, func(controller *MarketplaceCollectionController) error {
		return controller.validNftsEarningsOnSale()
	}, func(controller *MarketplaceCollectionController) error {
		return controller.validNftsEarningsOnYieldReward()
	}, func(controller *MarketplaceCollectionController) error {
		return controller.validInitiallyNftsVaultLockPercentage()
	}, func(controller *MarketplaceCollectionController) error {
		return controller.validNftsEarningAndLockPercentage()
	})
	return c
}

func (c *MarketplaceCollectionController) validNftsEarningAndLockPercentage() error {
	earnAndLockPercentage := sdk.NewInt(0)

	if !c.collection.InitiallyNftsVaultLockPercentage.IsNil() {
		earnAndLockPercentage = earnAndLockPercentage.Add(*c.collection.InitiallyNftsVaultLockPercentage)
	}

	for _, earning := range c.collection.NftsEarningsOnSale {
		if !earning.Percentage.IsNil() {
			earnAndLockPercentage = earnAndLockPercentage.Add(*earning.Percentage)
		}
	}

	if earnAndLockPercentage.GT(sdk.NewInt(100)) {
		return sdkerrors.Wrapf(types.ErrInvalidInitiallyEarnAndLockSummaryPercentage, "the summary of initially vault lock percentage and nfts earnings on sale percentages should not be greater than 100%%, value %s%%", earnAndLockPercentage)
	}

	return nil
}

func (c *MarketplaceCollectionController) validInitiallyNftsVaultLockPercentage() error {
	if c.collection.InitiallyNftsVaultLockPercentage.IsNegative() || c.collection.InitiallyNftsVaultLockPercentage.GT(sdk.NewInt(100)) {
		return sdkerrors.Wrapf(types.ErrInvalidInitiallyNftsVaultLockPercentage, "initially nfts vault lock percentage %s is invalid", c.collection.InitiallyNftsVaultLockPercentage)
	}

	return nil
}

func (c *MarketplaceCollectionController) validInitiallyNftMinPrice() error {
	parsed, err := sdk.ParseCoinNormalized(c.collection.InitiallyNftCollectionOwnerNftsMinPrice)
	if err != nil || parsed.IsNegative() {
		return sdkerrors.Wrapf(types.ErrInvalidInitiallyNftMinPrice, "initially nft min price %s is invalid", c.collection.InitiallyNftCollectionOwnerNftsMinPrice)
	}

	return nil
}

func (c *MarketplaceCollectionController) validNftsEarningsOnSale() error {
	if uint32(len(c.collection.NftsEarningsOnSale)) > c.conf.ValidNftsEarningsOnSaleMaxCount {
		return sdkerrors.Wrapf(types.ErrInvalidNftsEarningsOnSaleMaxCount, "nfts earnings on sale count %d invalid, max %d", len(c.collection.NftsEarningsOnSale), c.conf.ValidNftsEarningsOnSaleMaxCount)
	}

	for i, earning := range c.collection.NftsEarningsOnSale {
		if _, err := sdk.AccAddressFromBech32(earning.Address); err != nil {
			return sdkerrors.Wrapf(types.ErrInvalidNftsEarningsOnSaleAddress, "nfts earnings on sale address %s is invalid, earning index %d", earning.Address, i)
		}
		if earning.Percentage.IsNegative() || earning.Percentage.GT(sdk.NewInt(100)) {
			return sdkerrors.Wrapf(types.ErrInvalidNftsEarningsOnSalePercentage, "nfts earnings on sale percentage %s is invalid, earning index %d", earning.Percentage, i)
		}
		if types.ValidateMarketplaceEarningType(types.MarketplaceEarningType(earning.Type)) != nil {
			return sdkerrors.Wrapf(types.ErrInvalidMarketplaceNftsEarningsOnSaleEarningType, earning.Type)
		}
	}

	return nil
}

func (c *MarketplaceCollectionController) validNftsEarningsOnYieldReward() error {
	if uint32(len(c.collection.NftsEarningsOnYieldReward)) > c.conf.ValidNftsEarningsOnYieldRewardMaxCount {
		return sdkerrors.Wrapf(types.ErrInvalidNftsEarningsOnYieldRewardMaxCount, "nfts earnings on yield reward count %d invalid, max %d", len(c.collection.NftsEarningsOnYieldReward), c.conf.ValidNftsEarningsOnYieldRewardMaxCount)
	}

	for i, earning := range c.collection.NftsEarningsOnYieldReward {
		if _, err := sdk.AccAddressFromBech32(earning.Address); err != nil {
			return sdkerrors.Wrapf(types.ErrInvalidNftsEarningsOnYieldRewardAddress, "nfts earnings on yield reward address %s is invalid, earning index %d", earning.Address, i)
		}
		if earning.Percentage.IsNegative() || earning.Percentage.GT(sdk.NewInt(100)) {
			return sdkerrors.Wrapf(types.ErrInvalidNftsEarningsOnYieldRewardPercentage, "nfts earnings on yield reward percentage %s is invalid, earning index %d", earning.Percentage, i)
		}
		if types.ValidateMarketplaceEarningType(types.MarketplaceEarningType(earning.Type)) != nil {
			return sdkerrors.Wrapf(types.ErrInvalidMarketplaceNftsEarningsOnYieldRewardEarningType, earning.Type)
		}
	}

	return nil
}

func (c *MarketplaceCollectionController) getCollection() *types.MsgMarketplaceCollection {
	return c.collection
}

func (c *MarketplaceCollectionController) getMarketplaceCollection() *types.MarketplaceCollection {
	return c.marketplaceCollection
}

func (c *MarketplaceCollectionController) getIndex() []byte {
	return c.collectionIndex
}
