package keeper

import (
	"strings"

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
	}, func(controller *MarketplaceController) error {
		return controller.validMarketplaceMetadataOptions()
	}, func(controller *MarketplaceController) error {
		return controller.validMarketplaceMetadataImages()
	}, func(controller *MarketplaceController) error {
		return controller.validMarketplaceMetadataLinks()
	}, func(controller *MarketplaceController) error {
		return controller.validMarketplaceMetadataAttributes()
	})
	return c
}

func (c *MarketplaceController) validMarketplaceMetadataAttributes() error {
	if len(c.metadata.Attributes) == 0 {
		return nil
	}

	if int32(len(c.metadata.Attributes)) == c.conf.ValidMarketplaceMetadataAttributesMaxCount {
		return sdkerrors.Wrapf(types.ErrInvalidMarketplaceAttributesCount, "attributes count %d invalid, max %d", len(c.metadata.Attributes), c.conf.ValidMarketplaceMetadataAttributesMaxCount)
	}

	for k, attrubute := range c.metadata.Attributes {
		if attrubute.Type == "" || int32(len(attrubute.Type)) > c.conf.ValidMarketplaceMetadataAttributesTypeMaxLength {
			return sdkerrors.Wrapf(types.ErrInvalidMarketplaceAttribute, "attrubute index %d type empty or too long, max %d", k, c.conf.ValidMarketplaceMetadataAttributesTypeMaxLength)
		}
		if int32(len(attrubute.Value)) > c.conf.ValidMarketplaceMetadataAttributesValueMaxLength || int32(len(attrubute.SubValue)) > c.conf.ValidMarketplaceMetadataAttributesSubValueMaxLength {
			return sdkerrors.Wrapf(types.ErrInvalidMarketplaceAttribute, "attrubute index %d value/subvalue too long, max %d/%d symbols", k, c.conf.ValidMarketplaceMetadataAttributesValueMaxLength, c.conf.ValidMarketplaceMetadataAttributesSubValueMaxLength)
		}
	}

	return nil
}

func (c *MarketplaceController) validMarketplaceMetadataImages() error {
	if len(c.metadata.Images) == 0 {
		return nil
	}
	if int32(len(c.metadata.Images)) > c.conf.ValidMarketplaceMetadataImagesMaxCount {
		return sdkerrors.Wrapf(types.ErrInvalidMarketplaceImagesCount, "images count %d invalid, max %d", len(c.metadata.Images), c.conf.ValidMarketplaceMetadataImagesMaxCount)
	}
	for i, image := range c.metadata.Images {
		if image.Type == "" || int32(len(image.Type)) > c.conf.ValidMarketplaceMetadataImagesTypeMaxLength {
			return sdkerrors.Wrapf(types.ErrInvalidMarketplaceImage, "image type empty or length invalid, index %d , max %d", i, c.conf.ValidMarketplaceMetadataImagesTypeMaxLength)
		}
		if !utils.IsUrl(image.Url) {
			return sdkerrors.Wrapf(types.ErrInvalidMarketplaceImage, "image index %d invalid url", i)
		}
	}
	return nil
}

func (c *MarketplaceController) validMarketplaceMetadataLinks() error {
	if len(c.metadata.Links) == 0 {
		return nil
	}
	if int32(len(c.metadata.Links)) > c.conf.ValidMarketplaceMetadataLinksMaxCount {
		return sdkerrors.Wrapf(types.ErrInvalidMarketplaceLinksCount, "links count %d invalid, max %d", len(c.metadata.Links), c.conf.ValidMarketplaceMetadataLinksMaxCount)
	}
	for i, link := range c.metadata.Links {
		if link.Type == "" || int32(len(link.Type)) > c.conf.ValidMarketplaceMetadataLinksTypeMaxLength {
			return sdkerrors.Wrapf(types.ErrInvalidMarketplaceLink, "link type empty or length invalid, index %d , max %d", i, c.conf.ValidMarketplaceMetadataLinksTypeMaxLength)
		}
		if !utils.IsUrl(link.Url) {
			return sdkerrors.Wrapf(types.ErrInvalidMarketplaceLink, "link index %d invalid url", i)
		}
	}
	return nil
}

func (c *MarketplaceController) validMarketplaceMetadataOptions() error {
	if len(c.metadata.Options) == 0 {
		return nil
	}
	if int32(len(c.metadata.Options)) > c.conf.ValidMarketplaceMetadataOptionsMaxCount {
		return sdkerrors.Wrapf(types.ErrInvalidMarketplaceOptionsCount, "options count %d invalid, max %d", len(c.metadata.Options), c.conf.ValidMarketplaceMetadataOptionsMaxCount)
	}
	for i, option := range c.metadata.Options {
		if option.Type == "" || int32(len(option.Type)) > c.conf.ValidMarketplaceMetadataOptionsTypeMaxLength {
			return sdkerrors.Wrapf(types.ErrInvalidMarketplaceOption, "option type empty or length invalid, index %d , max %d", i, c.conf.ValidMarketplaceMetadataOptionsTypeMaxLength)
		}
		if int32(len(option.Value)) > c.conf.ValidMarketplaceMetadataOptionsValueMaxLength || int32(len(option.SubValue)) > c.conf.ValidMarketplaceMetadataOptionsSubValueMaxLength {
			return sdkerrors.Wrapf(types.ErrInvalidMarketplaceOption, "option index %d value/subvalue too long, max %d/%d symbols", i, c.conf.ValidMarketplaceMetadataOptionsValueMaxLength, c.conf.ValidMarketplaceMetadataOptionsSubValueMaxLength)
		}
	}
	return nil
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

	if int32(len(c.metadata.Description)) > c.conf.ValidMarketplaceMetadataDescriptionMaxLength {
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
	if strings.TrimSpace(c.metadata.Name) == "" {
		return nil
	}

	if int32(len(c.metadata.Name)) > c.conf.ValidMarketplaceMetadataNameMaxLength {
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
