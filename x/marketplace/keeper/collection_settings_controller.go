package keeper

import (
	"github.com/LimeChain/mantrachain/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type CollectionSettingsControllerFunc func(controller *CollectionSettingsController) error

type CollectionSettingsController struct {
	validators         []CollectionSettingsControllerFunc
	settings           *types.MsgCollectionSettings
	collectionSettings *types.CollectionSettings
	marketplaceIndex   []byte
	collectionIndex    []byte
	marketplaceId      string
	collectionId       string
	store              msgServer
	conf               *types.Params
	ctx                sdk.Context
}

func NewCollectionSettingsController(ctx sdk.Context, marketplaceIndex []byte, collectionIndex []byte, marketplaceId string, collectionId string) *CollectionSettingsController {
	return &CollectionSettingsController{
		marketplaceIndex: marketplaceIndex,
		collectionIndex:  collectionIndex,
		marketplaceId:    marketplaceId,
		collectionId:     collectionId,
		ctx:              ctx,
	}
}

func (c *CollectionSettingsController) WithSettings(settings *types.MsgCollectionSettings) *CollectionSettingsController {
	c.settings = settings
	return c
}

func (c *CollectionSettingsController) WithStore(store msgServer) *CollectionSettingsController {
	c.store = store
	return c
}

func (c *CollectionSettingsController) WithConfiguration(cfg types.Params) *CollectionSettingsController {
	c.conf = &cfg
	return c
}

func (c *CollectionSettingsController) Validate() error {
	for _, check := range c.validators {
		if err := check(c); err != nil {
			return err
		}
	}
	return nil
}

func (c *CollectionSettingsController) MustNotExist() *CollectionSettingsController {
	c.validators = append(c.validators, func(controller *CollectionSettingsController) error {
		return controller.mustNotExist()
	})
	return c
}

func (c *CollectionSettingsController) MustExist() *CollectionSettingsController {
	c.validators = append(c.validators, func(controller *CollectionSettingsController) error {
		return controller.mustExist()
	})
	return c
}

func (c *CollectionSettingsController) mustExist() error {
	return c.requireCollectionSettings()
}

func (c *CollectionSettingsController) mustNotExist() error {
	err := c.requireCollectionSettings()
	if err == nil {
		return sdkerrors.Wrapf(types.ErrCollectionSettingsAlreadyExists, "already exists: marketplace id %s, collection id %s", c.marketplaceId, c.collectionId)
	}
	return nil
}

func (c *CollectionSettingsController) requireCollectionSettings() error {
	if c.collectionSettings != nil {
		return nil
	}
	collectionSettings, isFound := c.store.GetCollectionSettings(c.ctx, c.marketplaceIndex, c.getIndex())
	if !isFound {
		return sdkerrors.Wrapf(types.ErrCollectionSettingsDoesNotExist, "not found: marketplace id %s, collection id %s", c.marketplaceId, c.collectionId)
	}
	c.collectionSettings = &collectionSettings
	return nil
}

func (c *CollectionSettingsController) ValidSettings() *CollectionSettingsController {
	c.validators = append(c.validators, func(controller *CollectionSettingsController) error {
		return controller.validInitiallyNftMinPrice()
	}, func(controller *CollectionSettingsController) error {
		return controller.validNftsEarningsOnSale()
	}, func(controller *CollectionSettingsController) error {
		return controller.validNftsEarningsOnYieldReward()
	}, func(controller *CollectionSettingsController) error {
		return controller.validInitiallyNftsVaultLockPercentage()
	}, func(controller *CollectionSettingsController) error {
		return controller.validNftsEarningAndLockPercentage()
	})
	return c
}

func (c *CollectionSettingsController) validNftsEarningAndLockPercentage() error {
	earnAndLockPercentage := sdk.NewInt(0)

	if !c.settings.InitiallyNftsVaultLockPercentage.IsNil() {
		earnAndLockPercentage = earnAndLockPercentage.Add(*c.settings.InitiallyNftsVaultLockPercentage)
	}

	for _, earning := range c.settings.NftsEarningsOnSale {
		if !earning.Percentage.IsNil() {
			earnAndLockPercentage = earnAndLockPercentage.Add(*earning.Percentage)
		}
	}

	if earnAndLockPercentage.GT(sdk.NewInt(100)) {
		return sdkerrors.Wrapf(types.ErrInvalidInitiallyEarnAndLockSummaryPercentage, "the summary of initially vault lock percentage and nfts earnings on sale percentages should not be greater than 100%%, value %s%%", earnAndLockPercentage)
	}

	return nil
}

func (c *CollectionSettingsController) validInitiallyNftsVaultLockPercentage() error {
	if c.settings.InitiallyNftsVaultLockPercentage.IsNegative() || c.settings.InitiallyNftsVaultLockPercentage.GT(sdk.NewInt(100)) {
		return sdkerrors.Wrapf(types.ErrInvalidInitiallyNftsVaultLockPercentage, "initially nfts vault lock percentage %s is invalid", c.settings.InitiallyNftsVaultLockPercentage)
	}

	return nil
}

func (c *CollectionSettingsController) validInitiallyNftMinPrice() error {
	parsed, err := sdk.ParseCoinNormalized(c.settings.InitiallyCollectionOwnerNftsMinPrice)
	if err != nil || parsed.IsNegative() {
		return sdkerrors.Wrapf(types.ErrInvalidInitiallyNftMinPrice, "initially nft min price %s is invalid", c.settings.InitiallyCollectionOwnerNftsMinPrice)
	}

	return nil
}

func (c *CollectionSettingsController) validNftsEarningsOnSale() error {
	if uint32(len(c.settings.NftsEarningsOnSale)) > c.conf.ValidNftsEarningsOnSaleMaxCount {
		return sdkerrors.Wrapf(types.ErrInvalidNftsEarningsOnSaleMaxCount, "nfts earnings on sale count %d invalid, max %d", len(c.settings.NftsEarningsOnSale), c.conf.ValidNftsEarningsOnSaleMaxCount)
	}

	for i, earning := range c.settings.NftsEarningsOnSale {
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

func (c *CollectionSettingsController) validNftsEarningsOnYieldReward() error {
	if uint32(len(c.settings.NftsEarningsOnYieldReward)) > c.conf.ValidNftsEarningsOnYieldRewardMaxCount {
		return sdkerrors.Wrapf(types.ErrInvalidNftsEarningsOnYieldRewardMaxCount, "nfts earnings on yield reward count %d invalid, max %d", len(c.settings.NftsEarningsOnYieldReward), c.conf.ValidNftsEarningsOnYieldRewardMaxCount)
	}

	for i, earning := range c.settings.NftsEarningsOnYieldReward {
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

func (c *CollectionSettingsController) getSettings() *types.MsgCollectionSettings {
	return c.settings
}

func (c *CollectionSettingsController) getCollectionSettings() *types.CollectionSettings {
	return c.collectionSettings
}

func (c *CollectionSettingsController) getIndex() []byte {
	return c.collectionIndex
}
