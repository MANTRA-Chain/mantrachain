package keeper

import (
	"github.com/LimeChain/mantrachain/x/bridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type BridgeControllerFunc func(controller *BridgeController) error

type BridgeController struct {
	validators []BridgeControllerFunc
	metadata   *types.MsgBridgeMetadata
	bridge     *types.Bridge
	store      msgServer
	conf       *types.Params
	ctx        sdk.Context
	creator    sdk.AccAddress
}

func NewBridgeController(ctx sdk.Context, creator sdk.AccAddress) *BridgeController {
	return &BridgeController{
		ctx:     ctx,
		creator: creator,
	}
}

func (c *BridgeController) WithMetadata(metadata *types.MsgBridgeMetadata) *BridgeController {
	c.metadata = metadata
	return c
}

func (c *BridgeController) WithStore(store msgServer) *BridgeController {
	c.store = store
	return c
}

func (c *BridgeController) WithId(id string) *BridgeController {
	if c.metadata == nil {
		c.metadata = &types.MsgBridgeMetadata{}
	}
	c.metadata.Id = id
	return c
}

func (c *BridgeController) WithBridge(bridge types.Bridge) *BridgeController {
	c.bridge = &bridge
	return c
}

func (c *BridgeController) WithConfiguration(cfg types.Params) *BridgeController {
	c.conf = &cfg
	return c
}

func (c *BridgeController) Validate() error {
	for _, check := range c.validators {
		if err := check(c); err != nil {
			return err
		}
	}
	return nil
}

func (c *BridgeController) MustNotExist() *BridgeController {
	c.validators = append(c.validators, func(controller *BridgeController) error {
		return controller.mustNotExist()
	})
	return c
}

func (c *BridgeController) MustExist() *BridgeController {
	c.validators = append(c.validators, func(controller *BridgeController) error {
		return controller.mustExist()
	})
	return c
}

func (c *BridgeController) ValidMetadata() *BridgeController {
	// TODO: Validate options, attrubute, images and links
	c.validators = append(c.validators, func(controller *BridgeController) error {
		return controller.bridgeMetadataNotNil()
	}, func(controller *BridgeController) error {
		return controller.validBridgeMetadataId()
	}, func(controller *BridgeController) error {
		return controller.validBridgeMetadataBridgeAccount()
	}, func(controller *BridgeController) error {
		return controller.validBridgeMetadataCw20ContractInitParams()
	})
	return c
}

func (c *BridgeController) HasOwner(owner sdk.AccAddress) *BridgeController {
	c.validators = append(c.validators, func(controller *BridgeController) error {
		return controller.hasOwner(owner)
	})
	return c
}

func (c *BridgeController) hasOwner(owner sdk.AccAddress) error {
	if err := c.requireBridge(); err != nil {
		panic("validation check is not allowed on a non existing bridge")
	}
	if owner.Equals(c.bridge.Owner) {
		return nil
	}
	return sdkerrors.Wrapf(types.ErrUnauthorized, "unauthorized")
}

func (c *BridgeController) mustExist() error {
	return c.requireBridge()
}

func (c *BridgeController) requireBridge() error {
	if c.bridge != nil {
		return nil
	}
	bridge, isFound := c.store.GetBridge(c.ctx, c.creator, c.getIndex())
	if !isFound {
		return sdkerrors.Wrapf(types.ErrBridgeDoesNotExist, "not found: %s", c.getId())
	}
	c.bridge = &bridge
	return nil
}

func (c *BridgeController) mustNotExist() error {
	err := c.requireBridge()
	if err == nil {
		return sdkerrors.Wrapf(types.ErrBridgeAlreadyExists, c.metadata.Id)
	}
	return nil
}

func (c *BridgeController) bridgeMetadataNotNil() error {
	// TODO: move to types -> validate basic
	if c.metadata == nil {
		return sdkerrors.Wrapf(types.ErrInvalidBridgeMetadata, "bridge metadata is invalid")
	}

	return nil
}

func (c *BridgeController) validBridgeMetadataId() error {
	return types.ValidateBridgeId(c.conf.ValidBridgeId, c.metadata.Id)
}

func (c *BridgeController) validBridgeMetadataBridgeAccount() error {
	_, err := sdk.AccAddressFromBech32(c.metadata.BridgeAccount)

	if err != nil {
		return err
	}

	return nil
}

func (c *BridgeController) validBridgeMetadataCw20ContractInitParams() error {
	// TODO: use strings.TrimSpace ... == "" for strings
	if len(c.metadata.Cw20ContractAddress) == 0 {
		if len(c.metadata.Cw20Name) == 0 {
			return sdkerrors.Wrapf(types.ErrInvalidCw20Name, "cw20 name is invalid")
		}

		if len(c.metadata.Cw20Symbol) == 0 {
			return sdkerrors.Wrapf(types.ErrInvalidCw20Symbol, "cw20 symbol is invalid")
		}

		if c.metadata.Cw20InitialBalances == nil || len(c.metadata.Cw20InitialBalances) == 0 {
			return sdkerrors.Wrapf(types.ErrInvalidCw20InitialBalances, "cw20 initial balances is invalid")
		}

		// TODO: check c.metadata.Cw20InitialBalances items are valid

		if c.metadata.Cw20Mint == nil {
			return sdkerrors.Wrapf(types.ErrInvalidCw20Mint, "cw20 mint is invalid")
		}

		// TODO: check c.metadata.Cw20Mint object is valid and minter address equals to bridge account
	}

	return nil
}

func (c *BridgeController) getBridge() *types.Bridge {
	return c.bridge
}

func (c *BridgeController) getId() string {
	return c.metadata.Id
}

func (c *BridgeController) getIndex() []byte {
	return types.GetBridgeIndex(c.creator, c.getId())
}
