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
	validators []NftControllerFunc
	metadata   []*types.MsgMintNftMetadata
	collIndex  []byte
	store      msgServer
	conf       *types.Params
	ctx        sdk.Context
}

func NewNftController(ctx sdk.Context, collIndex []byte, metadata []*types.MsgMintNftMetadata) *NftController {
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

func (c *NftController) FilterNotExist() *NftController {
	c.validators = append(c.validators, func(controller *NftController) error {
		return controller.filterNotExist()
	})
	return c
}

func (c *NftController) ValidNftMetadata() *NftController {
	c.validators = append(c.validators, func(controller *NftController) error {
		return controller.validNftMetadataLength()
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
	filtered := []*types.MsgMintNftMetadata{}

	for _, nft := range c.metadata {
		if nft.Id == "" {
			continue
		}

		index := types.GetNftIndex(c.collIndex, nft.Id)
		if !c.store.HasNft(c.ctx, c.collIndex, index) {
			filtered = append(filtered, nft)
		}
	}
	c.metadata = filtered

	return nil
}

func (c *NftController) validNftMetadataLength() error {
	if len(c.metadata) > 100 {
		return sdkerrors.Wrapf(types.ErrInvalidNftsLength, "nfts length %d invalid, max 100", len(c.metadata))
	}

	return nil
}

func (c *NftController) validNftMetadataId() error {
	for i, nft := range c.metadata {
		if nft.Id == "" {
			return sdkerrors.Wrapf(types.ErrInvalidNftId, "id: %s, index: %d", nft.Id, i)
		}

		validator := regexp.MustCompile(c.conf.ValidNftId)

		if !validator.MatchString(nft.Id) {
			return sdkerrors.Wrapf(types.ErrInvalidNftId, "id: %s, index: %d", nft.Id, i)
		}
	}

	return nil
}

func (c *NftController) validNftMetadataTitle() error {
	for i, nft := range c.metadata {
		if nft.Title == "" {
			continue
		}

		if len(nft.Title) > 100 {
			return sdkerrors.Wrapf(types.ErrInvalidNftTitle, "title length %d invalid, max 100, index: %d", nft.Title, i)
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
			return sdkerrors.Wrapf(types.ErrInvalidNftUrl, "%s invalid url, index: %d", nft.Title, i)
		}
	}

	return nil
}

func (c *NftController) validNftMetadataDescription() error {
	for i, nft := range c.metadata {
		if nft.Description == "" {
			continue
		}

		if len(nft.Description) > 1000 {
			return sdkerrors.Wrapf(types.ErrInvalidNftDescription, "%s description too long, max 1000 symbols, index: %d", nft.Description, i)
		}
	}

	return nil
}

func (c *NftController) validNftMetadataImages() error {
	for i, nft := range c.metadata {
		if len(nft.Images) == 0 {
			continue
		}

		for k, image := range nft.Images {
			if image.Type == "" || len(image.Type) > 100 {
				return sdkerrors.Wrapf(types.ErrInvalidNftCollectionImage, "image index %d invalid type, nft index: %d", k, i)
			}
			if !utils.IsUrl(image.Url) {
				return sdkerrors.Wrapf(types.ErrInvalidNftCollectionImage, "image index %d invalid url, nft index: %d", k, i)
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

		for k, link := range nft.Links {
			if link.Type == "" || len(link.Type) > 100 {
				return sdkerrors.Wrapf(types.ErrInvalidNftCollectionLink, "link index %d invalid type, nft index: %d", k, i)
			}
			if !utils.IsUrl(link.Url) {
				return sdkerrors.Wrapf(types.ErrInvalidNftCollectionLink, "link index %d invalid url, nft index: %d", k, i)
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

		for k, attrubute := range nft.Attributes {
			if attrubute.Type == "" || len(attrubute.Type) > 100 {
				return sdkerrors.Wrapf(types.ErrInvalidNftAttribute, "attrubute index %d invalid type, nft index: %d", k, i)
			}
			if len(attrubute.Value) > 100 || len(attrubute.SubValue) > 100 {
				return sdkerrors.Wrapf(types.ErrInvalidNftAttribute, "attrubute index %d value/subvalue too long, max 1000 symbols, nft index: %d", k, i)
			}
		}
	}

	return nil
}

func (c *NftController) getFiltered() []*types.MsgMintNftMetadata {
	return c.metadata
}
