package keeper

import (
	"strings"

	"cosmossdk.io/errors"
	"github.com/MANTRA-Finance/mantrachain/x/token/types"
	"github.com/MANTRA-Finance/mantrachain/x/token/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type NftCollectionControllerFunc func(controller *NftCollectionController) error

type NftCollectionController struct {
	actions       []NftCollectionControllerFunc
	validators    []NftCollectionControllerFunc
	metadata      *types.MsgCreateNftCollectionMetadata
	strict        bool
	nftCollection *types.NftCollection
	store         msgServer
	conf          *types.Params
	ctx           sdk.Context
	creator       sdk.AccAddress
}

func NewNftCollectionController(ctx sdk.Context, creator sdk.AccAddress, strict bool) *NftCollectionController {
	return &NftCollectionController{
		ctx:     ctx,
		creator: creator,
		strict:  strict,
	}
}

func (c *NftCollectionController) WithMetadata(metadata *types.MsgCreateNftCollectionMetadata) *NftCollectionController {
	c.metadata = metadata
	return c
}

func (c *NftCollectionController) WithStore(store msgServer) *NftCollectionController {
	c.store = store
	return c
}

func (c *NftCollectionController) WithId(id string) *NftCollectionController {
	if c.metadata == nil {
		c.metadata = &types.MsgCreateNftCollectionMetadata{}
	}
	c.metadata.Id = id
	return c
}

func (c *NftCollectionController) WithNftCollection(collection types.NftCollection) *NftCollectionController {
	c.nftCollection = &collection
	return c
}

func (c *NftCollectionController) WithConfiguration(cfg types.Params) *NftCollectionController {
	c.conf = &cfg
	return c
}

func (c *NftCollectionController) CreateDefaultIfNotExists() *NftCollectionController {
	if !c.strict && strings.TrimSpace(c.metadata.Id) == "" {
		c.actions = append(c.actions, func(controller *NftCollectionController) error {
			return controller.CreateDefaultNftCollection()
		})
	}
	return c
}

func (c *NftCollectionController) CreateDefaultNftCollection() error {
	c.metadata.Id = ""
	c.metadata.Category = string(types.GeneralNftCollectionCat)
	collectionIndex := c.getIndex()

	c.requireNftCollection()

	if c.nftCollection == nil {
		nftExecutor := NewNftExecutor(c.ctx, c.store.nftKeeper)
		err := nftExecutor.SetDefaultClass(collectionIndex)
		if err != nil {
			return err
		}

		newNftCollection := types.NftCollection{
			Index:    collectionIndex,
			Id:       c.getId(),
			Category: c.metadata.Category,
			Creator:  c.creator,
		}

		c.store.SetNftCollection(c.ctx, newNftCollection)
		c.store.SetOpenedNftsCollection(c.ctx, collectionIndex)
		c.store.SetNftCollectionOwner(c.ctx, collectionIndex, c.creator)

		c.nftCollection = &newNftCollection
	}

	return nil
}

func (c *NftCollectionController) Execute() error {
	for _, check := range c.actions {
		if err := check(c); err != nil {
			return err
		}
	}
	return nil
}

func (c *NftCollectionController) Validate() error {
	for _, check := range c.validators {
		if err := check(c); err != nil {
			return err
		}
	}
	return nil
}

func (c *NftCollectionController) MustNotExist() *NftCollectionController {
	c.validators = append(c.validators, func(controller *NftCollectionController) error {
		return controller.mustNotExist()
	})
	return c
}

func (c *NftCollectionController) MustExist() *NftCollectionController {
	c.validators = append(c.validators, func(controller *NftCollectionController) error {
		return controller.mustExist()
	})
	return c
}

func (c *NftCollectionController) MustNotBeDefault() *NftCollectionController {
	c.validators = append(c.validators, func(controller *NftCollectionController) error {
		return controller.mustNotBeDefault()
	})
	return c
}

func (c *NftCollectionController) OpenedOrOwner(owner sdk.AccAddress) *NftCollectionController {
	c.validators = append(c.validators, func(controller *NftCollectionController) error {
		return controller.openedOrOwner(owner)
	})
	return c
}

func (c *NftCollectionController) ValidMetadata() *NftCollectionController {
	c.validators = append(c.validators, func(controller *NftCollectionController) error {
		return controller.validCollectionMetadataId()
	}, func(controller *NftCollectionController) error {
		return controller.validCollectionMetadataName()
	}, func(controller *NftCollectionController) error {
		return controller.validCollectionMetadataCategory()
	}, func(controller *NftCollectionController) error {
		return controller.validCollectionMetadataUrl()
	}, func(controller *NftCollectionController) error {
		return controller.validCollectionMetadataDescription()
	}, func(controller *NftCollectionController) error {
		return controller.validCollectionMetadataSymbol()
	}, func(controller *NftCollectionController) error {
		return controller.validCollectionMetadataOptions()
	}, func(controller *NftCollectionController) error {
		return controller.validCollectionMetadataImages()
	}, func(controller *NftCollectionController) error {
		return controller.validCollectionMetadataLinks()
	}, func(controller *NftCollectionController) error {
		return controller.validCollectionMetadataOpened()
	}, func(controller *NftCollectionController) error {
		return controller.validCollectionMetadataSoulBondedNfts()
	}, func(controller *NftCollectionController) error {
		return controller.validCollectionMetadataRestricted()
	})
	return c
}

func (c *NftCollectionController) mustExist() error {
	return c.requireNftCollection()
}

func (c *NftCollectionController) mustNotBeDefault() error {
	if c.isDefault() {
		return errors.Wrap(types.ErrInvalidNftCollectionId, c.getId())
	}
	return nil
}

func (c *NftCollectionController) openedOrOwner(owner sdk.AccAddress) error {
	if err := c.requireNftCollection(); err != nil {
		return errors.Wrap(err, "validation check is not allowed on a non existing nftCollection")
	}

	if c.store.HasOpenedNftsCollection(
		c.ctx,
		c.getIndex(),
	) {
		return nil
	}

	collOwner, found := c.store.GetNftCollectionOwner(
		c.ctx,
		c.getIndex(),
	)

	if !found {
		return errors.Wrapf(types.ErrUnauthorized, "unauthorized")
	}

	if owner.Equals(sdk.AccAddress(collOwner)) {
		return nil
	}

	return errors.Wrapf(types.ErrUnauthorized, "unauthorized")
}

func (c *NftCollectionController) requireNftCollection() error {
	if c.nftCollection != nil {
		return nil
	}
	nftCollection, isFound := c.store.GetNftCollection(c.ctx, c.creator, c.getIndex())
	if !isFound {
		return errors.Wrapf(types.ErrNftCollectionDoesNotExist, "not found: %s", c.getId())
	}
	c.nftCollection = &nftCollection
	return nil
}

func (c *NftCollectionController) mustNotExist() error {
	err := c.requireNftCollection()
	if err == nil {
		return errors.Wrapf(types.ErrNftCollectionAlreadyExists, c.metadata.Name)
	}
	return nil
}

func (c *NftCollectionController) validCollectionMetadataCategory() error {
	if c.metadata.Category == "" {
		return nil
	}
	if _, err := types.ParseNftCollectionCategory(c.metadata.Category); err != nil {
		return errors.Wrapf(types.ErrInvalidNftCollectionCategory, c.metadata.Category)
	}
	return nil
}

func (c *NftCollectionController) validCollectionMetadataSymbol() error {
	if c.metadata.Symbol == "" {
		return nil
	}

	if int32(len(c.metadata.Symbol)) < c.conf.ValidNftCollectionMetadataSymbolMinLength {
		return errors.Wrapf(types.ErrInvalidNftCollectionSymbol, "%s symbol too short, min %d letters", c.metadata.Symbol, c.conf.ValidNftCollectionMetadataSymbolMinLength)
	}

	if int32(len(c.metadata.Symbol)) > c.conf.ValidNftCollectionMetadataSymbolMaxLength {
		return errors.Wrapf(types.ErrInvalidNftCollectionSymbol, "%s symbol too long, max %d letters", c.metadata.Symbol, c.conf.ValidNftCollectionMetadataSymbolMaxLength)
	}

	return nil
}

func (c *NftCollectionController) validCollectionMetadataUrl() error {
	if c.metadata.Url == "" {
		return nil
	}

	if !utils.IsUrl(c.metadata.Url) {
		return errors.Wrapf(types.ErrInvalidNftCollectionUrl, "%s invalid url", c.metadata.Url)
	}

	return nil
}

func (c *NftCollectionController) validCollectionMetadataDescription() error {
	if c.metadata.Description == "" {
		return nil
	}

	if int32(len(c.metadata.Description)) > c.conf.ValidNftCollectionMetadataDescriptionMaxLength {
		return errors.Wrapf(types.ErrInvalidNftCollectionDescription, "description too long, max %d symbols", c.conf.ValidNftCollectionMetadataDescriptionMaxLength)
	}

	return nil
}

func (c *NftCollectionController) validCollectionMetadataId() error {
	return types.ValidateNftCollectionId(c.conf.ValidNftCollectionId, c.getId())
}

func (c *NftCollectionController) validCollectionMetadataName() error {
	if len(c.metadata.Name) == 0 {
		return nil
	}

	if int32(len(c.metadata.Name)) > c.conf.ValidNftCollectionMetadataNameMaxLength {
		return errors.Wrapf(types.ErrInvalidNftCollectionName, "name length %d invalid, max %d", len(c.metadata.Name), c.conf.ValidNftCollectionMetadataNameMaxLength)
	}

	return nil
}

func (c *NftCollectionController) validCollectionMetadataImages() error {
	if len(c.metadata.Images) == 0 {
		return nil
	}
	if int32(len(c.metadata.Images)) > c.conf.ValidNftCollectionMetadataImagesMaxCount {
		return errors.Wrapf(types.ErrInvalidNftCollectionImagesCount, "images count %d invalid, max %d", len(c.metadata.Images), c.conf.ValidNftCollectionMetadataImagesMaxCount)
	}
	for i, image := range c.metadata.Images {
		if image.Type == "" || int32(len(image.Type)) > c.conf.ValidNftCollectionMetadataImagesTypeMaxLength {
			return errors.Wrapf(types.ErrInvalidNftCollectionImage, "image type empty or length invalid, index %d , max %d", i, c.conf.ValidNftCollectionMetadataImagesTypeMaxLength)
		}
		if !utils.IsUrl(image.Url) {
			return errors.Wrapf(types.ErrInvalidNftCollectionImage, "image index %d invalid url", i)
		}
	}
	return nil
}

func (c *NftCollectionController) validCollectionMetadataLinks() error {
	if len(c.metadata.Links) == 0 {
		return nil
	}
	if int32(len(c.metadata.Links)) > c.conf.ValidNftCollectionMetadataLinksMaxCount {
		return errors.Wrapf(types.ErrInvalidNftCollectionLinksCount, "links count %d invalid, max %d", len(c.metadata.Links), c.conf.ValidNftCollectionMetadataLinksMaxCount)
	}
	for i, link := range c.metadata.Links {
		if link.Type == "" || int32(len(link.Type)) > c.conf.ValidNftCollectionMetadataLinksTypeMaxLength {
			return errors.Wrapf(types.ErrInvalidNftCollectionLink, "link type empty or length invalid, index %d , max %d", i, c.conf.ValidNftCollectionMetadataLinksTypeMaxLength)
		}
		if !utils.IsUrl(link.Url) {
			return errors.Wrapf(types.ErrInvalidNftCollectionLink, "link index %d invalid url", i)
		}
	}
	return nil
}

func (c *NftCollectionController) validCollectionMetadataOptions() error {
	if len(c.metadata.Options) == 0 {
		return nil
	}
	if int32(len(c.metadata.Options)) > c.conf.ValidNftCollectionMetadataOptionsMaxCount {
		return errors.Wrapf(types.ErrInvalidNftCollectionOptionsCount, "options count %d invalid, max %d", len(c.metadata.Options), c.conf.ValidNftCollectionMetadataOptionsMaxCount)
	}
	for i, option := range c.metadata.Options {
		if option.Type == "" || int32(len(option.Type)) > c.conf.ValidNftCollectionMetadataOptionsTypeMaxLength {
			return errors.Wrapf(types.ErrInvalidNftCollectionOption, "option type empty or length invalid, index %d , max %d", i, c.conf.ValidNftCollectionMetadataOptionsTypeMaxLength)
		}
		if int32(len(option.Value)) > c.conf.ValidNftCollectionMetadataOptionsValueMaxLength || int32(len(option.SubValue)) > c.conf.ValidNftCollectionMetadataOptionsSubValueMaxLength {
			return errors.Wrapf(types.ErrInvalidNftCollectionOption, "option index %d value/subvalue too long, max %d/%d symbols", i, c.conf.ValidNftCollectionMetadataOptionsValueMaxLength, c.conf.ValidNftCollectionMetadataOptionsSubValueMaxLength)
		}
	}
	return nil
}

func (c *NftCollectionController) validCollectionMetadataOpened() error {
	if !c.isDefault() {
		return nil
	}
	if !c.metadata.Opened {
		return errors.Wrapf(types.ErrInvalidNftCollectionOpened, "collection %d must be opened", len(c.getId()))
	}
	return nil
}

func (c *NftCollectionController) validCollectionMetadataSoulBondedNfts() error {
	if !c.isDefault() {
		return nil
	}
	if c.metadata.SoulBondedNfts {
		return errors.Wrapf(types.ErrInvalidNftCollectionSoulBondedNfts, "collection %d cannot be with soul bonded nfts", len(c.getId()))
	}
	return nil
}

func (c *NftCollectionController) validCollectionMetadataRestricted() error {
	if !c.metadata.RestrictedNfts {
		return nil
	}
	if c.metadata.Opened {
		return errors.Wrapf(types.ErrInvalidNftCollectionRestricted, "opened collection %d cannot be restricted", len(c.getId()))
	}
	if c.isDefault() {
		return errors.Wrapf(types.ErrInvalidNftCollectionRestricted, "collection %d cannot be restricted", len(c.getId()))
	}
	return nil
}

func (c *NftCollectionController) isDefault() bool {
	return c.getId() == c.conf.NftCollectionDefaultId
}

func (c *NftCollectionController) getId() string {
	if strings.TrimSpace(c.metadata.Id) == "" && !c.strict {
		return c.conf.NftCollectionDefaultId
	}
	return c.metadata.Id
}

func (c *NftCollectionController) getIndex() []byte {
	id := c.getId()
	return types.GetNftCollectionIndex(c.creator, id)
}
