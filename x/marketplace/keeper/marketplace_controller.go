package keeper

import (
	"github.com/LimeChain/mantrachain/x/marketplace/types"
	"github.com/LimeChain/mantrachain/x/marketplace/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type MarketplaceControllerFunc func(controller *MarketplaceController) error

type MarketplaceController struct {
	validators  []MarketplaceControllerFunc
	metadata    *types.MsgMarketplaceMetadata
	marketplace *types.Marketplace
	store       msgServer
	conf        *types.Params
	ctx         sdk.Context
	creator     sdk.AccAddress
}

func NewMarketplaceController(ctx sdk.Context, creator sdk.AccAddress) *MarketplaceController {
	return &MarketplaceController{
		ctx:     ctx,
		creator: creator,
	}
}

func (c *MarketplaceController) WithMetadata(metadata *types.MsgMarketplaceMetadata) *MarketplaceController {
	c.metadata = metadata
	return c
}

func (c *MarketplaceController) WithStore(store msgServer) *MarketplaceController {
	c.store = store
	return c
}

func (c *MarketplaceController) WithId(id string) *MarketplaceController {
	if c.metadata == nil {
		c.metadata = &types.MsgMarketplaceMetadata{}
	}
	c.metadata.Id = id
	return c
}

func (c *MarketplaceController) WithMarketplace(collection types.Marketplace) *MarketplaceController {
	c.marketplace = &collection
	return c
}

func (c *MarketplaceController) WithConfiguration(cfg types.Params) *MarketplaceController {
	c.conf = &cfg
	return c
}

func (c *MarketplaceController) Validate() error {
	for _, check := range c.validators {
		if err := check(c); err != nil {
			return err
		}
	}
	return nil
}

func (c *MarketplaceController) IsOpenedOrHasOwner(owner sdk.AccAddress) *MarketplaceController {
	c.validators = append(c.validators, func(controller *MarketplaceController) error {
		return controller.isOpenedOrHasOwner(owner)
	})
	return c
}

func (c *MarketplaceController) MustNotExist() *MarketplaceController {
	c.validators = append(c.validators, func(controller *MarketplaceController) error {
		return controller.mustNotExist()
	})
	return c
}

func (c *MarketplaceController) MustExist() *MarketplaceController {
	c.validators = append(c.validators, func(controller *MarketplaceController) error {
		return controller.mustExist()
	})
	return c
}

func (c *MarketplaceController) ValidMetadata() *MarketplaceController {
	// TODO: Validate options, attrubute, images and links
	c.validators = append(c.validators, func(controller *MarketplaceController) error {
		return controller.marketplaceMetadataNotNil()
	}, func(controller *MarketplaceController) error {
		return controller.validMarketplaceMetadataId()
	}, func(controller *MarketplaceController) error {
		return controller.validMarketplaceMetadataName()
	}, func(controller *MarketplaceController) error {
		return controller.validMarketplaceMetadataUrl()
	}, func(controller *MarketplaceController) error {
		return controller.validMarketplaceMetadataDescription()
	})
	return c
}

func (c *MarketplaceController) isOpenedOrHasOwner(owner sdk.AccAddress) error {
	if err := c.requireMarketplace(); err != nil {
		panic("validation check is not allowed on a non existing marketplace")
	}

	if c.marketplace.Opened {
		return nil
	}

	if owner.Equals(c.marketplace.Owner) {
		return nil
	}
	return sdkerrors.Wrapf(types.ErrUnauthorized, "unauthorized")
}

func (c *MarketplaceController) HasOwner(owner sdk.AccAddress) *MarketplaceController {
	c.validators = append(c.validators, func(controller *MarketplaceController) error {
		return controller.hasOwner(owner)
	})
	return c
}

func (c *MarketplaceController) hasOwner(owner sdk.AccAddress) error {
	if err := c.requireMarketplace(); err != nil {
		panic("validation check is not allowed on a non existing marketplace")
	}
	if owner.Equals(c.marketplace.Owner) {
		return nil
	}
	return sdkerrors.Wrapf(types.ErrUnauthorized, "unauthorized")
}

func (c *MarketplaceController) mustExist() error {
	return c.requireMarketplace()
}

func (c *MarketplaceController) requireMarketplace() error {
	if c.marketplace != nil {
		return nil
	}
	marketplace, isFound := c.store.GetMarketplace(c.ctx, c.creator, c.getIndex())
	if !isFound {
		return sdkerrors.Wrapf(types.ErrMarketplaceDoesNotExist, "not found: %s", c.getId())
	}
	c.marketplace = &marketplace
	return nil
}

func (c *MarketplaceController) mustNotExist() error {
	err := c.requireMarketplace()
	if err == nil {
		return sdkerrors.Wrapf(types.ErrMarketplaceAlreadyExists, c.metadata.Name)
	}
	return nil
}

func (c *MarketplaceController) validMarketplaceMetadataUrl() error {
	if c.metadata.Url == "" {
		return nil
	}

	if !utils.IsUrl(c.metadata.Url) {
		return sdkerrors.Wrapf(types.ErrInvalidMarketplaceUrl, "%s invalid url", c.metadata.Url)
	}

	return nil
}

func (c *MarketplaceController) validMarketplaceMetadataDescription() error {
	if c.metadata.Description == "" {
		return nil
	}

	if uint32(len(c.metadata.Description)) > c.conf.ValidMarketplaceMetadataDescriptionMaxLength {
		return sdkerrors.Wrapf(types.ErrInvalidMarketplaceDescription, "description too long, max %d symbols", c.conf.ValidMarketplaceMetadataDescriptionMaxLength)
	}

	return nil
}

func (c *MarketplaceController) marketplaceMetadataNotNil() error {
	if c.metadata == nil {
		return sdkerrors.Wrapf(types.ErrInvalidMarketplaceMetadata, "marketplace metadata is invalid")
	}

	return nil
}

func (c *MarketplaceController) validMarketplaceMetadataId() error {
	return types.ValidateMarketplaceId(c.conf.ValidMarketplaceId, c.metadata.Id)
}

func (c *MarketplaceController) validMarketplaceMetadataName() error {
	if len(c.metadata.Name) == 0 {
		return nil
	}

	if uint32(len(c.metadata.Name)) > c.conf.ValidMarketplaceMetadataNameMaxLength {
		return sdkerrors.Wrapf(types.ErrInvalidMarketplaceName, "name length %d invalid, max %d", len(c.metadata.Name), c.conf.ValidMarketplaceMetadataNameMaxLength)
	}

	return nil
}

func (c *MarketplaceController) getMarketplace() *types.Marketplace {
	return c.marketplace
}

func (c *MarketplaceController) getId() string {
	return c.metadata.Id
}

func (c *MarketplaceController) getIndex() []byte {
	return types.GetMarketplaceIndex(c.creator, c.getId())
}
