package types

import (
	"fmt"
	"regexp"

	"github.com/cosmos/cosmos-sdk/codec"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	DefaultValidDomain     string = "^[-_a-z0-9]{4,16}$"
	DefaultValidDomainName string = "^[-_\\.a-z0-9]{1,64}$"
)

var (
	KeyValidDomain     = []byte("ValidDomain")
	KeyValidDomainName = []byte("ValidDomainName")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	validDomain string,
	validDomainName string,
) Params {
	return Params{
		ValidDomain:     validDomain,
		ValidDomainName: validDomainName,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultValidDomain,
		DefaultValidDomainName,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyValidDomain, &p.ValidDomain, validateValidDomain),
		paramtypes.NewParamSetPair(KeyValidDomainName, &p.ValidDomainName, validateValidDomainName),
	}
}

// unmarshal the current staking params value from store key or panic
func MustUnmarshalParams(cdc *codec.LegacyAmino, value []byte) Params {
	params, err := UnmarshalParams(cdc, value)
	if err != nil {
		panic(err)
	}

	return params
}

// unmarshal the current staking params value from store key
func UnmarshalParams(cdc *codec.LegacyAmino, value []byte) (params Params, err error) {
	err = cdc.Unmarshal(value, &params)
	if err != nil {
		return
	}

	return
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateValidDomain(p.ValidDomain); err != nil {
		return err
	}

	if err := validateValidDomainName(p.ValidDomainName); err != nil {
		return err
	}

	return nil
}

func validateValidDomain(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == "" {
		return fmt.Errorf("valid domain param regex should not be empty")
	}

	_, err := regexp.Compile(v)
	if err != nil {
		return fmt.Errorf("valid domain param is invalid regex %s", v)
	}

	return nil
}

func validateValidDomainName(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == "" {
		return fmt.Errorf("valid domain name param regex should not be empty")
	}

	_, err := regexp.Compile(v)
	if err != nil {
		return fmt.Errorf("valid domain name param is invalid regex %s", v)
	}

	return nil
}
