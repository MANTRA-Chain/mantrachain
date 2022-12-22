package keeper

import (
	"strings"

	"github.com/LimeChain/mantrachain/x/bridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type BridgeControllerFunc func(controller *BridgeController) error

type BridgeController struct {
	validators    []BridgeControllerFunc
	metadata      *types.MsgBridgeMetadata
	bridge        *types.Bridge
	store         msgServer
	conf          *types.Params
	ctx           sdk.Context
	bridgeCreator sdk.AccAddress
}

func NewBridgeController(ctx sdk.Context, bridgeCreator sdk.AccAddress) *BridgeController {
	return &BridgeController{
		ctx:           ctx,
		bridgeCreator: bridgeCreator,
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
	bridge, isFound := c.store.GetBridge(c.ctx, c.bridgeCreator, c.getIndex())
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
	if strings.TrimSpace(c.metadata.Cw20ContractAddress) == "" {
		if strings.TrimSpace(c.metadata.Cw20Name) == "" {
			return sdkerrors.Wrapf(types.ErrInvalidCw20Name, "cw20 name is invalid")
		}

		if int32(len(c.metadata.Cw20Name)) < c.conf.ValidBridgeCw20ContractNameMinLength {
			return sdkerrors.Wrapf(types.ErrInvalidCw20Name, "cw20 name length %d invalid, min %d", len(c.metadata.Cw20Name), c.conf.ValidBridgeCw20ContractNameMinLength)
		}

		if int32(len(c.metadata.Cw20Name)) > c.conf.ValidBridgeCw20ContractNameMaxLength {
			return sdkerrors.Wrapf(types.ErrInvalidCw20Name, "cw20 name length %d invalid, max %d", len(c.metadata.Cw20Name), c.conf.ValidBridgeCw20ContractNameMaxLength)
		}

		if strings.TrimSpace(c.metadata.Cw20Symbol) == "" {
			return sdkerrors.Wrapf(types.ErrInvalidCw20Symbol, "cw20 symbol is invalid")
		}

		if int32(len(c.metadata.Cw20Symbol)) < c.conf.ValidBridgeCw20ContractSymbolMinLength {
			return sdkerrors.Wrapf(types.ErrInvalidCw20Symbol, "cw20 symbol length %d invalid, min %d", len(c.metadata.Cw20Symbol), c.conf.ValidBridgeCw20ContractSymbolMinLength)
		}

		if int32(len(c.metadata.Cw20Symbol)) > c.conf.ValidBridgeCw20ContractSymbolMaxLength {
			return sdkerrors.Wrapf(types.ErrInvalidCw20Symbol, "cw20 symbol length %d invalid, max %d", len(c.metadata.Cw20Symbol), c.conf.ValidBridgeCw20ContractSymbolMaxLength)
		}

		if c.metadata.Cw20InitialBalances == nil || len(c.metadata.Cw20InitialBalances) == 0 {
			return sdkerrors.Wrapf(types.ErrInvalidCw20InitialBalances, "cw20 initial balances is invalid")
		}

		for i, initialBalance := range c.metadata.Cw20InitialBalances {
			if _, err := sdk.AccAddressFromBech32(initialBalance.Address); err != nil {
				return sdkerrors.Wrapf(types.ErrInvalidCw20InitialBalancesAddress, "cw20 initial balances address %s is invalid, index %d", initialBalance.Address, i)
			}

			if strings.TrimSpace(initialBalance.Amount) == "" {
				return sdkerrors.Wrapf(types.ErrInvalidCw20InitialBalancesAmount, "cw20 initial balances amount %s is invalid, index %d", initialBalance.Amount, i)
			}
		}

		if c.metadata.Cw20Mint == nil {
			return sdkerrors.Wrapf(types.ErrInvalidCw20Mint, "cw20 mint is invalid")
		}

		minter, err := sdk.AccAddressFromBech32(c.metadata.Cw20Mint.Minter)

		if err != nil {
			return sdkerrors.Wrapf(types.ErrInvalidCw20MintMinter, "cw20 mint minter address %s is invalid", c.metadata.Cw20Mint.Minter)
		}

		bridgeAccount, _ := sdk.AccAddressFromBech32(c.metadata.BridgeAccount)

		if !minter.Equals(bridgeAccount) {
			return sdkerrors.Wrapf(types.ErrBridgeAccountMismatch, "bridge account %s does not match the minter %s", bridgeAccount.String(), minter.String())
		}
	} else {
		if _, err := sdk.AccAddressFromBech32(c.metadata.Cw20ContractAddress); err != nil {
			return sdkerrors.Wrapf(types.ErrInvalidCw20ContractAddress, "cw20 contract address %s is invalid", c.metadata.Cw20ContractAddress)
		}
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
	return types.GetBridgeIndex(c.bridgeCreator, c.getId())
}
