package keeper

import (
	"regexp"

	"github.com/LimeChain/mantrachain/x/mns/types"
	"github.com/LimeChain/mantrachain/x/mns/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// DomainControllerFunc is the function signature for domain validation functions
type DomainControllerFunc func(controller *DomainController) error

type DomainController struct {
	validators []DomainControllerFunc
	msgDomain  string
	domain     *types.Domain
	store      msgServer
	conf       *types.Params
	ctx        sdk.Context
}

func NewDomainController(ctx sdk.Context, msgDomain string) *DomainController {
	return &DomainController{
		msgDomain: msgDomain,
		ctx:       ctx,
	}
}

func (c *DomainController) WithStore(store msgServer) *DomainController {
	c.store = store
	return c
}

func (c *DomainController) WithDomain(dom types.Domain) *DomainController {
	c.domain = &dom
	c.msgDomain = dom.Domain
	return c
}

func (c *DomainController) WithConfiguration(cfg types.Params) *DomainController {
	c.conf = &cfg
	return c
}

func (c *DomainController) Validate() error {
	for _, check := range c.validators {
		if err := check(c); err != nil {
			return err
		}
	}
	return nil
}

func (c *DomainController) Type(Type types.DomainType) *DomainController {
	c.validators = append(c.validators, func(controller *DomainController) error {
		return controller.dType(Type)
	})
	return c
}

func (c *DomainController) MustNotExist() *DomainController {
	c.validators = append(c.validators, func(controller *DomainController) error {
		return controller.mustNotExist()
	})
	return c
}

func (c *DomainController) MustExist() *DomainController {
	c.validators = append(c.validators, func(controller *DomainController) error {
		return controller.mustExist()
	})
	return c
}

func (c *DomainController) ValidDomain() *DomainController {
	c.validators = append(c.validators, func(controller *DomainController) error {
		return controller.validDomain()
	})
	return c
}

func (c *DomainController) NotExpired() *DomainController {
	c.validators = append(c.validators, func(controller *DomainController) error {
		return controller.notExpired()
	})
	return c
}

func (c *DomainController) HasOwnerIfClosed(owner sdk.AccAddress) *DomainController {
	c.validators = append(c.validators, func(controller *DomainController) error {
		return controller.hasOwnerIfClosed(owner)
	})
	return c
}

func (c *DomainController) hasOwnerIfClosed(owner sdk.AccAddress) error {
	// assert domain exists
	if err := c.requireDomain(); err != nil {
		panic("validation check is not allowed on a non existing domain")
	}
	if types.DomainType(c.domain.DomainType) == types.OpenDomain {
		return nil
	}
	if owner.Equals(c.domain.Owner) {
		return nil
	}
	// if it has expired return error
	return sdkerrors.Wrapf(types.ErrDomainClosed, "%s closed", c.msgDomain)
}

func (c *DomainController) notExpired() error {
	// assert domain exists
	if err := c.requireDomain(); err != nil {
		panic("validation check is not allowed on a non existing domain")
	}
	// if domain has no expiration time
	if c.domain.ExpireAt == 0 {
		return nil
	}
	// check if domain has expired
	expireTime := utils.SecondsToTime(c.domain.ExpireAt)
	// if block time is before expiration, return nil
	if c.ctx.BlockTime().Before(expireTime) {
		return nil
	}
	// if it has expired return error
	return sdkerrors.Wrapf(types.ErrDomainExpired, "%s has expired", c.msgDomain)
}

func (c *DomainController) mustExist() error {
	return c.requireDomain()
}

func (c *DomainController) dType(Type types.DomainType) error {
	if err := c.requireDomain(); err != nil {
		panic("validation check is not allowed on a non existing domain")
	}
	if types.DomainType(c.domain.DomainType) != Type {
		return sdkerrors.Wrapf(types.ErrInvalidDomainType, "operation not allowed on invalid domain type %s, expected %s", c.domain.DomainType, Type)
	}
	return nil
}

func (c *DomainController) requireDomain() error {
	if c.domain != nil {
		return nil
	}
	domain, isFound := c.store.GetDomain(c.ctx, c.msgDomain)
	if !isFound {
		return sdkerrors.Wrapf(types.ErrDomainDoesNotExist, "not found: %s", c.msgDomain)
	}
	c.domain = &domain
	return nil
}

func (c *DomainController) mustNotExist() error {
	err := c.requireDomain()
	if err == nil {
		return sdkerrors.Wrapf(types.ErrDomainAlreadyExists, c.msgDomain)
	}
	return nil
}

func (c *DomainController) validDomain() error {
	if c.conf == nil {
		panic("configuration is missing")
	}

	validator := regexp.MustCompile(c.conf.ValidDomain)

	if !validator.MatchString(c.msgDomain) {
		return sdkerrors.Wrap(types.ErrInvalidDomainName, c.msgDomain)
	}

	return nil
}
