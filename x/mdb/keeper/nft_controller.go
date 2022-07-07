package keeper

import (
	"regexp"

	"github.com/LimeChain/mantrachain/x/mdb/types"
	"github.com/LimeChain/mantrachain/x/mdb/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type NftControllerFunc func(controller *NftController) error

type NftController struct {
	actions         []NftControllerFunc
	validators      []NftControllerFunc
	metadata        []*types.MsgNftMetadata
	collectionIndex []byte
	store           msgServer
	conf            *types.Params
	ctx             sdk.Context
}

func NewNftController(ctx sdk.Context, collectionIndex []byte) *NftController {
	return &NftController{
		collectionIndex: collectionIndex,
		ctx:             ctx,
	}
}

func (c *NftController) WithMetadata(metadata []*types.MsgNftMetadata) *NftController {
	c.metadata = metadata
	return c
}

func (c *NftController) WithId(id string) *NftController {
	if c.metadata == nil {
		c.metadata = make([]*types.MsgNftMetadata, 1)
	}

	c.metadata = append(c.metadata, &types.MsgNftMetadata{
		Id: id,
	})
	return c
}

func (c *NftController) WithIds(ids []string) *NftController {
	if c.metadata == nil {
		c.metadata = make([]*types.MsgNftMetadata, 0)
	}

	for _, id := range ids {
		c.metadata = append(c.metadata, &types.MsgNftMetadata{
			Id: id,
		})
	}
	return c
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

func (c *NftController) FilterCannotTransfer(operator sdk.AccAddress) *NftController {
	c.actions = append(c.actions, func(controller *NftController) error {
		return controller.filterCannotTransfer(operator)
	})
	return c
}

func (c *NftController) ValidMetadata() *NftController {
	c.validators = append(c.validators, func(controller *NftController) error {
		return controller.validMetadataMaxCount()
	}, func(controller *NftController) error {
		return controller.validMetadataId()
	}, func(controller *NftController) error {
		return controller.validMetadataTitle()
	}, func(controller *NftController) error {
		return controller.validMetadataUrl()
	}, func(controller *NftController) error {
		return controller.validMetadataDescription()
	}, func(controller *NftController) error {
		return controller.validMetadataImages()
	}, func(controller *NftController) error {
		return controller.validMetadataLinks()
	}, func(controller *NftController) error {
		return controller.validMetadataAttributes()
	})
	return c
}

func (c *NftController) filterNotExist() error {
	byIndex := make(map[string]*types.MsgNftMetadata, 0)
	filtered := []*types.MsgNftMetadata{}
	indexes := [][]byte{}

	for _, nftMetadata := range c.metadata {
		if nftMetadata.Id == "" {
			continue
		}

		index := types.GetNftIndex(c.collectionIndex, nftMetadata.Id)
		indexes = append(indexes, index)
		byIndex[string(index)] = nftMetadata
	}

	indexes = c.store.FilterNotExists(c.ctx, c.collectionIndex, indexes)
	for _, index := range indexes {
		nftMetadata := byIndex[string(index)]
		filtered = append(filtered, nftMetadata)
	}

	c.metadata = filtered

	return nil
}

func (c *NftController) filterNotOwn(owner sdk.AccAddress) error {
	byId := make(map[string]*types.MsgNftMetadata, 0)
	filtered := []*types.MsgNftMetadata{}
	ids := make([]string, 0)

	for _, nftMetadata := range c.metadata {
		id := nftMetadata.Id
		ids = append(ids, id)
		byId[id] = nftMetadata
	}

	ids = c.store.nftKeeper.FilterNotOwn(c.ctx, string(c.collectionIndex), ids)
	for _, id := range ids {
		nftMetadata := byId[id]
		filtered = append(filtered, nftMetadata)
	}

	c.metadata = filtered

	return nil
}

func (c *NftController) filterCannotTransfer(operator sdk.AccAddress) error {
	filtered := []*types.MsgNftMetadata{}

	for _, nftMetadata := range c.metadata {
		id := nftMetadata.Id
		index := types.GetNftIndex(c.collectionIndex, nftMetadata.Id)
		owner := c.store.nftKeeper.GetOwner(c.ctx, string(c.collectionIndex), id)
		if owner.Equals(operator) || c.store.IsApproved(c.ctx, c.collectionIndex, index, owner, operator) {
			filtered = append(filtered, nftMetadata)
		}
	}

	c.metadata = filtered

	return nil
}

func (c *NftController) validMetadataMaxCount() error {
	if int32(len(c.metadata)) > c.conf.ValidNftMetadataMaxCount {
		return sdkerrors.Wrapf(types.ErrInvalidNftsCount, "nfts count %d invalid, max %d", len(c.metadata), c.conf.ValidNftMetadataMaxCount)
	}

	return nil
}

func (c *NftController) validMetadataId() error {
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

func (c *NftController) validMetadataTitle() error {
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

func (c *NftController) validMetadataUrl() error {
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

func (c *NftController) validMetadataDescription() error {
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

func (c *NftController) validMetadataImages() error {
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

func (c *NftController) validMetadataLinks() error {
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

func (c *NftController) validMetadataAttributes() error {
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

func (c *NftController) getMetadata() []*types.MsgNftMetadata {
	return c.metadata
}

func (c *NftController) getNftsIds() []string {
	var nftsIds []string

	for _, nftMetadata := range c.metadata {
		index := types.GetNftIndex(c.collectionIndex, nftMetadata.Id)
		nftsIds = append(nftsIds, string(index))
	}

	return nftsIds
}

func (c *NftController) getIndexes() [][]byte {
	var indexes [][]byte

	for _, nftMetadata := range c.metadata {
		index := types.GetNftIndex(c.collectionIndex, nftMetadata.Id)
		indexes = append(indexes, index)
	}

	return indexes
}
