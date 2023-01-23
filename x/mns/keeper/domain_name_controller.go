package keeper

import (
	"regexp"

	"github.com/LimeChain/mantrachain/x/mns/types"
	"github.com/LimeChain/mantrachain/x/mns/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// DomainNameControllerFunc is the function signature for domain validation functions
type DomainNameControllerFunc func(controller *DomainNameController) error

type DomainNameController struct {
	validators    []DomainNameControllerFunc
	msgDomain     string
	msgDomainName string
	domainName    *types.DomainName
	store         msgServer
	conf          *types.Params
	ctx           sdk.Context
}

func NewDomainNameController(ctx sdk.Context, msgDomain string, msgDomainName string) *DomainNameController {
	return &DomainNameController{
		msgDomain:     msgDomain,
		msgDomainName: msgDomainName,
		ctx:           ctx,
	}
}

func (c *DomainNameController) WithStore(store msgServer) *DomainNameController {
	c.store = store
	return c
}

func (c *DomainNameController) WithDomainName(domName types.DomainName) *DomainNameController {
	c.domainName = &domName
	c.msgDomain = domName.Domain
	c.msgDomainName = domName.DomainName
	return c
}

func (c *DomainNameController) WithConfiguration(cfg types.Params) *DomainNameController {
	c.conf = &cfg
	return c
}

func (c *DomainNameController) Validate() error {
	for _, check := range c.validators {
		if err := check(c); err != nil {
			return err
		}
	}
	return nil
}

func (c *DomainNameController) MustNotExist() *DomainNameController {
	c.validators = append(c.validators, func(controller *DomainNameController) error {
		return controller.mustNotExist()
	})
	return c
}

func (c *DomainNameController) MustExist() *DomainNameController {
	c.validators = append(c.validators, func(controller *DomainNameController) error {
		return controller.mustExist()
	})
	return c
}

func (c *DomainNameController) ValidDomainName() *DomainNameController {
	c.validators = append(c.validators, func(controller *DomainNameController) error {
		return controller.validDomainName()
	})
	return c
}

func (c *DomainNameController) NotExpired() *DomainNameController {
	c.validators = append(c.validators, func(controller *DomainNameController) error {
		return controller.notExpired()
	})
	return c
}

func (c *DomainNameController) notExpired() error {
	// assert domain exists
	if err := c.requireDomainName(); err != nil {
		return sdkerrors.Wrap(err, "validation check is not allowed on a non existing domain")
	}
	// if domain has no expiration time
	if c.domainName.ExpireAt == 0 {
		return nil
	}
	// check if domain has expired
	expireTime := utils.SecondsToTime(c.domainName.ExpireAt)
	// if block time is before expiration, return nil
	if c.ctx.BlockTime().Before(expireTime) {
		return nil
	}
	// if it has expired return error
	return sdkerrors.Wrapf(types.ErrDomainExpired, "%s has expired", c.msgDomainName)
}

func (c *DomainNameController) mustExist() error {
	return c.requireDomainName()
}

func (c *DomainNameController) requireDomainName() error {
	if c.domainName != nil {
		return nil
	}
	domainName, isFound := c.store.GetDomainName(c.ctx, c.msgDomain, c.msgDomainName)
	if !isFound {
		return sdkerrors.Wrapf(types.ErrDomainDoesNotExist, "not found: %s", c.msgDomainName)
	}
	c.domainName = &domainName
	return nil
}

func (c *DomainNameController) mustNotExist() error {
	err := c.requireDomainName()
	if err == nil {
		return sdkerrors.Wrapf(types.ErrDomainAlreadyExists, c.msgDomainName)
	}
	return nil
}

func (c *DomainNameController) validDomainName() error {
	if c.conf == nil {
		panic("configuration is missing")
	}

	validator := regexp.MustCompile(c.conf.ValidDomainName)

	if !validator.MatchString(c.msgDomainName) {
		return sdkerrors.Wrap(types.ErrInvalidDomainName, c.msgDomainName)
	}

	return nil
}
