package keeper

import (
	"regexp"
	"strconv"

	"github.com/LimeChain/mantrachain/x/mdb/types"
	"github.com/LimeChain/mantrachain/x/mdb/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NftCollectionControllerFunc is the function signature for nftCollection validation functions
type NftCollectionControllerFunc func(controller *NftCollectionController) error

type NftCollectionController struct {
	actions       []NftCollectionControllerFunc
	validators    []NftCollectionControllerFunc
	metadata      *types.MsgCreateNftCollectionMetadata
	nftCollection *types.NftCollection
	store         msgServer
	conf          *types.Params
	ctx           sdk.Context
	creator       sdk.AccAddress
}

func NewNftCollectionController(ctx sdk.Context, metadata *types.MsgCreateNftCollectionMetadata, creator sdk.AccAddress) *NftCollectionController {
	return &NftCollectionController{
		metadata: metadata,
		ctx:      ctx,
		creator:  creator,
	}
}

func (c *NftCollectionController) WithStore(store msgServer) *NftCollectionController {
	c.store = store
	return c
}

func (c *NftCollectionController) WithNftCollection(coll types.NftCollection) *NftCollectionController {
	c.nftCollection = &coll
	return c
}

func (c *NftCollectionController) WithConfiguration(cfg types.Params) *NftCollectionController {
	c.conf = &cfg
	return c
}

func (c *NftCollectionController) CreateDefaultIfNotExists() *NftCollectionController {
	if c.metadata.Id == "" {
		c.actions = append(c.actions, func(controller *NftCollectionController) error {
			return controller.CreateDefault()
		})
	} else {
		c.actions = append(c.actions, func(controller *NftCollectionController) error {
			return c.MustExist().Validate()
		})
	}
	return c
}

func (c *NftCollectionController) CreateDefault() error {
	c.metadata.Id = "default"
	c.metadata.Opened = true
	collIndex := c.getIndex()

	err := c.requireNftCollection()

	if err != nil {
		return err
	}

	if c.nftCollection == nil {
		return nil
	}

	nftExecutor := NewNftExecutor(c.ctx, c.store.nftKeeper)
	_, err = nftExecutor.SetDefaultClass(collIndex)
	if err != nil {
		return err
	}

	newNftCollection := types.NftCollection{
		Index:   collIndex,
		Opened:  c.metadata.Opened,
		Creator: c.creator,
		Owner:   c.creator,
	}

	c.store.SetNftCollection(c.ctx, newNftCollection)

	c.nftCollection = &newNftCollection

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

func (c *NftCollectionController) ValidNftCollectionMetadata() *NftCollectionController {
	c.validators = append(c.validators, func(controller *NftCollectionController) error {
		return controller.validNftCollectionMetadataId()
	}, func(controller *NftCollectionController) error {
		return controller.validNftCollectionMetadataName()
	}, func(controller *NftCollectionController) error {
		return controller.validNftCollectionMetadataCategory()
	}, func(controller *NftCollectionController) error {
		return controller.validNftCollectionMetadataUrl()
	}, func(controller *NftCollectionController) error {
		return controller.validNftCollectionMetadataDescription()
	}, func(controller *NftCollectionController) error {
		return controller.validNftCollectionMetadataDisplayTheme()
	}, func(controller *NftCollectionController) error {
		return controller.validNftCollectionMetadataSymbol()
	}, func(controller *NftCollectionController) error {
		return controller.validNftCollectionMetadataCreatorEarnings()
	}, func(controller *NftCollectionController) error {
		return controller.validNftCollectionMetadataImages()
	}, func(controller *NftCollectionController) error {
		return controller.validNftCollectionMetadataLinks()
	}, func(controller *NftCollectionController) error {
		return controller.validNftCollectionMetadataOpened()
	})
	return c
}

func (c *NftCollectionController) HasOwner(owner sdk.AccAddress) *NftCollectionController {
	c.validators = append(c.validators, func(controller *NftCollectionController) error {
		return controller.hasOwner(owner)
	})
	return c
}

func (c *NftCollectionController) hasOwner(owner sdk.AccAddress) error {
	// assert nftCollection exists
	if err := c.requireNftCollection(); err != nil {
		panic("validation check is not allowed on a non existing nftCollection")
	}
	if owner.Equals(c.nftCollection.Owner) {
		return nil
	}
	// if it has expired return error
	return sdkerrors.Wrapf(types.ErrUnauthorized, "unauthorized")
}

func (c *NftCollectionController) mustExist() error {
	return c.requireNftCollection()
}

func (c *NftCollectionController) mustNotBeDefault() error {
	if c.metadata.Id == "default" {
		return sdkerrors.Wrap(types.ErrInvalidNftCollectionId, c.metadata.Id)
	}
	return nil
}

func (c *NftCollectionController) requireNftCollection() error {
	if c.nftCollection != nil {
		return nil
	}
	creator := sdk.AccAddress(c.creator)
	nftCollection, isFound := c.store.GetNftCollection(c.ctx, creator, c.getIndex())
	if !isFound {
		return sdkerrors.Wrapf(types.ErrNftCollectionDoesNotExist, "not found: %s", c.getIndex())
	}
	c.nftCollection = &nftCollection
	return nil
}

func (c *NftCollectionController) mustNotExist() error {
	err := c.requireNftCollection()
	if err == nil {
		return sdkerrors.Wrapf(types.ErrNftCollectionAlreadyExists, c.metadata.Name)
	}
	return nil
}

func (c *NftCollectionController) validNftCollectionMetadataCategory() error {
	if c.metadata.Category == "" {
		return nil
	}
	if types.ValidateNftCollectionCategory(types.NftCollectionCategory(c.metadata.Category)) != nil {
		return sdkerrors.Wrapf(types.ErrInvalidNftCollectionCategory, c.metadata.Category)
	}
	return nil
}

func (c *NftCollectionController) validNftCollectionMetadataDisplayTheme() error {
	if c.metadata.DisplayTheme == "" {
		return nil
	}
	if types.ValidateNftCollectionDisplayTheme(types.NftCollectionDisplayTheme(c.metadata.DisplayTheme)) != nil {
		return sdkerrors.Wrapf(types.ErrInvalidNftCollectionDisplayTheme, c.metadata.DisplayTheme)
	}
	return nil
}

func (c *NftCollectionController) validNftCollectionMetadataCreatorEarnings() error {
	if c.metadata.CreatorEarnings == "" {
		return nil
	}

	if _, err := strconv.ParseFloat(c.metadata.CreatorEarnings, 64); err != nil {
		return sdkerrors.Wrapf(types.ErrInvalidNftCollectionCreatorEarnings, "%s creator_ernings not a number", c.metadata.CreatorEarnings)
	}

	return nil
}

func (c *NftCollectionController) validNftCollectionMetadataSymbol() error {
	if c.metadata.Symbol == "" {
		return nil
	}

	if len(c.metadata.Symbol) < 2 {
		return sdkerrors.Wrapf(types.ErrInvalidNftCollectionSymbol, "%s symbol too short, min 2 letters", c.metadata.Symbol)
	}

	if len(c.metadata.Symbol) > 5 {
		return sdkerrors.Wrapf(types.ErrInvalidNftCollectionSymbol, "%s symbol too long, max 5 letters", c.metadata.Symbol)
	}

	return nil
}

func (c *NftCollectionController) validNftCollectionMetadataUrl() error {
	if c.metadata.Url == "" {
		return nil
	}

	if !utils.IsUrl(c.metadata.Url) {
		return sdkerrors.Wrapf(types.ErrInvalidNftCollectionUrl, "%s invalid url", c.metadata.Url)
	}

	return nil
}

func (c *NftCollectionController) validNftCollectionMetadataDescription() error {
	if c.metadata.Description == "" {
		return nil
	}

	if len(c.metadata.Description) > 1000 {
		return sdkerrors.Wrapf(types.ErrInvalidNftCollectionDescription, "%s description too long, max 1000 symbols", c.metadata.Description)
	}

	return nil
}

func (c *NftCollectionController) validNftCollectionMetadataId() error {
	validator := regexp.MustCompile(c.conf.ValidNftCollectionId)

	if !validator.MatchString(c.metadata.Id) {
		return sdkerrors.Wrap(types.ErrInvalidNftCollectionId, c.metadata.Id)
	}

	return nil
}

func (c *NftCollectionController) validNftCollectionMetadataName() error {
	if len(c.metadata.Name) == 0 {
		return nil
	}

	if len(c.metadata.Name) > 100 {
		return sdkerrors.Wrapf(types.ErrInvalidNftCollectionName, "name length %d invalid, max 100", len(c.metadata.Name))
	}

	return nil
}

func (c *NftCollectionController) validNftCollectionMetadataImages() error {
	if len(c.metadata.Images) == 0 {
		return nil
	}
	if len(c.metadata.Images) > 10 {
		return sdkerrors.Wrapf(types.ErrInvalidNftCollectionImage, "images length %d invalid, max 10", len(c.metadata.Images))
	}
	for i, image := range c.metadata.Images {
		if image.Type == "" || len(image.Type) > 100 {
			return sdkerrors.Wrapf(types.ErrInvalidNftCollectionImage, "image index %d invalid type", i)
		}
		if !utils.IsUrl(image.Url) {
			return sdkerrors.Wrapf(types.ErrInvalidNftCollectionImage, "image index %d invalid url", i)
		}
	}
	return nil
}

func (c *NftCollectionController) validNftCollectionMetadataLinks() error {
	if len(c.metadata.Links) == 0 {
		return nil
	}
	if len(c.metadata.Links) > 10 {
		return sdkerrors.Wrapf(types.ErrInvalidNftCollectionLink, "links length %d invalid, max 10", len(c.metadata.Links))
	}
	for i, link := range c.metadata.Links {
		if link.Type == "" || len(link.Type) > 100 {
			return sdkerrors.Wrapf(types.ErrInvalidNftCollectionLink, "link index %d invalid type", i)
		}
		if !utils.IsUrl(link.Url) {
			return sdkerrors.Wrapf(types.ErrInvalidNftCollectionLink, "link index %d invalid url", i)
		}
	}
	return nil
}

func (c *NftCollectionController) validNftCollectionMetadataOpened() error {
	if c.metadata.Id != "default" {
		return nil
	}
	if !c.metadata.Opened {
		return sdkerrors.Wrapf(types.ErrInvalidNftCollectionOpened, "collection %d can not be opened", len(c.metadata.Id))
	}
	return nil
}

func (c *NftCollectionController) getId() string {
	return c.metadata.Id
}

func (c *NftCollectionController) getIndex() []byte {
	id := c.getId()
	return types.GetNftCollectionIndex(c.creator, id)
}
