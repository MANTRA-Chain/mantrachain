package keeper

import (
	"regexp"

	"github.com/LimeChain/mantrachain/x/mdb/types"
	"github.com/LimeChain/mantrachain/x/mdb/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NftControllerFunc is the function signature for nft validation functions
type NftControllerFunc func(controller *NftController) error

type NftController struct {
	actions    []NftControllerFunc
	validators []NftControllerFunc
	metadata   []*types.MsgNftMetadata
	collIndex  []byte
	store      msgServer
	conf       *types.Params
	ctx        sdk.Context
}

func NewNftController(ctx sdk.Context, collIndex []byte, metadata []*types.MsgNftMetadata) *NftController {
	return &NftController{
		metadata:  metadata,
		collIndex: collIndex,
		ctx:       ctx,
	}
}

func (c *NftController) WithStore(store msgServer) *NftController {
	c.store = store
	return c
}

func (c *NftController) WithConfiguration(cfg types.Params) *NftController {
	c.conf = &cfg
	return c
}

func (c *NftController) Validate() error {
	for _, check := range c.validators {
		if err := check(c); err != nil {
			return err
		}
	}
	return nil
}

func (c *NftController) Execute() error {
	for _, check := range c.actions {
		if err := check(c); err != nil {
			return err
		}
	}
	return nil
}

func (c *NftController) FilterNotExist() *NftController {
	c.actions = append(c.actions, func(controller *NftController) error {
		return controller.filterNotExist()
	})
	return c
}

func (c *NftController) FilterNotOwn(owner sdk.AccAddress) *NftController {
	c.actions = append(c.actions, func(controller *NftController) error {
		return controller.filterNotOwn(owner)
	})
	return c
}

func (c *NftController) ValidNftMetadata() *NftController {
	c.validators = append(c.validators, func(controller *NftController) error {
		return controller.validNftMetadataCount()
	}, func(controller *NftController) error {
		return controller.validNftMetadataId()
	}, func(controller *NftController) error {
		return controller.validNftMetadataTitle()
	}, func(controller *NftController) error {
		return controller.validNftMetadataUrl()
	}, func(controller *NftController) error {
		return controller.validNftMetadataDescription()
	}, func(controller *NftController) error {
		return controller.validNftMetadataImages()
	}, func(controller *NftController) error {
		return controller.validNftMetadataLinks()
	}, func(controller *NftController) error {
		return controller.validNftMetadataAttributes()
	})
	return c
}

func (c *NftController) filterNotExist() error {
	filtered := []*types.MsgNftMetadata{}

	for _, nft := range c.metadata {
		if nft.Id == "" {
			continue
		}

		index := types.GetNftIndex(c.collIndex, nft.Id)
		if c.store.HasNft(c.ctx, c.collIndex, index) {
			filtered = append(filtered, nft)
		}
	}
	c.metadata = filtered

	return nil
}

func (c *NftController) filterNotOwn(owner sdk.AccAddress) error {
	filtered := []*types.MsgNftMetadata{}

	for _, nft := range c.metadata {
		if nft.Id == "" {
			continue
		}

		index := types.GetNftIndex(c.collIndex, nft.Id)
		nftOwner := c.store.nftKeeper.GetOwner(c.ctx, string(c.collIndex), string(index))

		if owner.Equals(nftOwner) {
			filtered = append(filtered, nft)
		}
	}
	c.metadata = filtered

	return nil
}

func (c *NftController) validNftMetadataCount() error {
	if len(c.metadata) == 0 {
		return sdkerrors.Wrapf(types.ErrInvalidNftsCount, "nfts count %d invalid, min 1", len(c.metadata))
	}
	if int32(len(c.metadata)) > c.conf.ValidNftMetadataMaxCount {
		return sdkerrors.Wrapf(types.ErrInvalidNftsCount, "nfts count %d invalid, max %d", len(c.metadata), c.conf.ValidNftMetadataMaxCount)
	}

	return nil
}

func (c *NftController) validNftMetadataId() error {
	for i, nft := range c.metadata {
		if nft.Id == "" {
			return sdkerrors.Wrapf(types.ErrInvalidNftId, "id: %s, index %d", nft.Id, i)
		}

		validator := regexp.MustCompile(c.conf.ValidNftId)

		if !validator.MatchString(nft.Id) {
			return sdkerrors.Wrapf(types.ErrInvalidNftId, "id: %s, index %d", nft.Id, i)
		}
	}

	return nil
}

func (c *NftController) validNftMetadataTitle() error {
	for i, nft := range c.metadata {
		if nft.Title == "" {
			continue
		}

		if int32(len(nft.Title)) > c.conf.ValidNftMetadataTitleMaxLength {
			return sdkerrors.Wrapf(types.ErrInvalidNftTitle, "title length %d invalid, max 100, index %d", nft.Title, i, c.conf.ValidNftMetadataTitleMaxLength)
		}
	}

	return nil
}

func (c *NftController) validNftMetadataUrl() error {
	for i, nft := range c.metadata {
		if nft.Url == "" {
			continue
		}

		if !utils.IsUrl(nft.Url) {
			return sdkerrors.Wrapf(types.ErrInvalidNftUrl, "%s invalid url, index %d", nft.Title, i)
		}
	}

	return nil
}

func (c *NftController) validNftMetadataDescription() error {
	for i, nft := range c.metadata {
		if nft.Description == "" {
			continue
		}

		if int32(len(nft.Description)) > c.conf.ValidNftMetadataDescriptionMaxLength {
			return sdkerrors.Wrapf(types.ErrInvalidNftDescription, "description too long, max %d symbols, index %d", c.conf.ValidNftMetadataDescriptionMaxLength, i)
		}
	}

	return nil
}

func (c *NftController) validNftMetadataImages() error {
	for i, nft := range c.metadata {
		if len(nft.Images) == 0 {
			continue
		}

		if int32(len(nft.Images)) == c.conf.ValidNftMetadataImagesMaxCount {
			return sdkerrors.Wrapf(types.ErrInvalidNftImagesCount, "images count %d invalid, max %d, nft index %d", len(nft.Images), c.conf.ValidNftMetadataImagesMaxCount, i)
		}

		for k, image := range nft.Images {
			if image.Type == "" || int32(len(image.Type)) > c.conf.ValidNftMetadataImagesTypeMaxLength {
				return sdkerrors.Wrapf(types.ErrInvalidNftImage, "image index %d type empty or too long, max %d, nft index %d", k, c.conf.ValidNftMetadataImagesTypeMaxLength, i)
			}
			if !utils.IsUrl(image.Url) {
				return sdkerrors.Wrapf(types.ErrInvalidNftImage, "image index %d invalid url, nft index %d", k, i)
			}
		}
	}

	return nil
}

func (c *NftController) validNftMetadataLinks() error {
	for i, nft := range c.metadata {
		if len(nft.Links) == 0 {
			continue
		}

		if int32(len(nft.Links)) == c.conf.ValidNftMetadataLinksMaxCount {
			return sdkerrors.Wrapf(types.ErrInvalidNftLinksCount, "links count %d invalid, max %d, nft index %d", len(nft.Links), c.conf.ValidNftMetadataLinksMaxCount, i)
		}

		for k, link := range nft.Links {
			if link.Type == "" || int32(len(link.Type)) > c.conf.ValidNftMetadataLinksTypeMaxLength {
				return sdkerrors.Wrapf(types.ErrInvalidNftLink, "link index %d type empty or too long, max %d, nft index %d", k, c.conf.ValidNftMetadataLinksTypeMaxLength, i)
			}
			if !utils.IsUrl(link.Url) {
				return sdkerrors.Wrapf(types.ErrInvalidNftLink, "link index %d invalid url, nft index %d", k, i)
			}
		}
	}

	return nil
}

func (c *NftController) validNftMetadataAttributes() error {
	for i, nft := range c.metadata {
		if len(nft.Attributes) == 0 {
			continue
		}

		if int32(len(nft.Attributes)) == c.conf.ValidNftMetadataAttributesMaxCount {
			return sdkerrors.Wrapf(types.ErrInvalidNftAttributesCount, "attributes count %d invalid, max %d, nft index %d", len(nft.Attributes), c.conf.ValidNftMetadataAttributesMaxCount, i)
		}

		for k, attrubute := range nft.Attributes {
			if attrubute.Type == "" || int32(len(attrubute.Type)) > c.conf.ValidNftMetadataAttributesTypeMaxLength {
				return sdkerrors.Wrapf(types.ErrInvalidNftAttribute, "attrubute index %d type empty or too long, max %d, nft index %d", k, c.conf.ValidNftMetadataAttributesTypeMaxLength, i)
			}
			if int32(len(attrubute.Value)) > c.conf.ValidNftMetadataAttributesValueMaxLength || int32(len(attrubute.SubValue)) > c.conf.ValidNftMetadataAttributesSubValueMaxLength {
				return sdkerrors.Wrapf(types.ErrInvalidNftAttribute, "attrubute index %d value/subvalue too long, max %d/%d symbols, nft index %d", k, c.conf.ValidNftMetadataAttributesValueMaxLength, c.conf.ValidNftMetadataAttributesSubValueMaxLength, i)
			}
		}
	}

	return nil
}

func (c *NftController) getFiltered() []*types.MsgNftMetadata {
	return c.metadata
}
