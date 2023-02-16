package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	DefaultAdminAccount           = ""
	DefaultTokenCollectionCreator = ""
	DefaultTokenCollectionId      = ""
	DefaultPriviliges             = 0x000000000000000000000000000000000000000000000000ffffffffffffffff
)

var (
	KeyAdminAccount           = []byte("AdminAccount")
	KeyTokenCollectionCreator = []byte("TokenCollectionCreator")
	KeyTokenCollectionId      = []byte("TokenCollectionId")
	KeyDefaultPriviliges      = []byte("DefaultPriviliges")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	adminAccount string,
	tokenCollectionCreator string,
	tokenCollectionId string,
	defaultPriviliges uint64,
) Params {
	return Params{
		AdminAccount:           adminAccount,
		TokenCollectionCreator: tokenCollectionCreator,
		TokenCollectionId:      tokenCollectionId,
		DefaultPriviliges:      defaultPriviliges,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultAdminAccount,
		DefaultTokenCollectionCreator,
		DefaultTokenCollectionId,
		DefaultPriviliges,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyAdminAccount, &p.AdminAccount, validateAdminAccount),
		paramtypes.NewParamSetPair(KeyTokenCollectionCreator, &p.TokenCollectionCreator, validateTokenCollectionCreator),
		paramtypes.NewParamSetPair(KeyTokenCollectionId, &p.TokenCollectionId, validateTokenCollectionId),
		paramtypes.NewParamSetPair(KeyDefaultPriviliges, &p.DefaultPriviliges, validateDefaultPriviliges),
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
	if err := validateAdminAccount(p.AdminAccount); err != nil {
		return err
	}

	if err := validateTokenCollectionCreator(p.TokenCollectionCreator); err != nil {
		return err
	}
	if err := validateTokenCollectionId(p.TokenCollectionId); err != nil {
		return err
	}

	return nil
}

func validateAdminAccount(i interface{}) error {
	_, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateTokenCollectionCreator(i interface{}) error {
	_, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateTokenCollectionId(i interface{}) error {
	_, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateDefaultPriviliges(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
