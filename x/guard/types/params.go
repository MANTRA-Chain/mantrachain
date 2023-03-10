package types

import (
	"fmt"
	"math/big"

	"github.com/cosmos/cosmos-sdk/codec"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	DefaultAdminAccount                            = ""
	DefaultAccountPrivilegesTokenCollectionCreator = ""
	DefaultAccountPrivilegesTokenCollectionId      = ""
	DefaultAccountPrivileges                       = big.NewInt(0).Sub(big.NewInt(0).Exp(big.NewInt(2), big.NewInt(64), nil), big.NewInt(1)).Bytes()
)

var (
	KeyAdminAccount                            = []byte("AdminAccount")
	KeyAccountPrivilegesTokenCollectionCreator = []byte("AccountPrivilegesTokenCollectionCreator")
	KeyAccountPrivilegesTokenCollectionId      = []byte("AccountPrivilegesTokenCollectionId")
	KeyDefaultAccountPrivileges                = []byte("DefaultAccountPrivileges")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	adminAccount string,
	accountPrivilegesTokenCollectionCreator string,
	accountPrivilegesTokenCollectionId string,
	defaultAccountPrivileges []byte,
) Params {
	return Params{
		AdminAccount:                            adminAccount,
		AccountPrivilegesTokenCollectionCreator: accountPrivilegesTokenCollectionCreator,
		AccountPrivilegesTokenCollectionId:      accountPrivilegesTokenCollectionId,
		DefaultAccountPrivileges:                defaultAccountPrivileges,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultAdminAccount,
		DefaultAccountPrivilegesTokenCollectionCreator,
		DefaultAccountPrivilegesTokenCollectionId,
		DefaultAccountPrivileges,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyAdminAccount, &p.AdminAccount, validateAdminAccount),
		paramtypes.NewParamSetPair(KeyAccountPrivilegesTokenCollectionCreator, &p.AccountPrivilegesTokenCollectionCreator, validateAccountPrivilegesTokenCollectionCreator),
		paramtypes.NewParamSetPair(KeyAccountPrivilegesTokenCollectionId, &p.AccountPrivilegesTokenCollectionId, validateAccountPrivilegesTokenCollectionId),
		paramtypes.NewParamSetPair(KeyDefaultAccountPrivileges, &p.DefaultAccountPrivileges, validateDefaultAccountPrivileges),
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

	if err := validateAccountPrivilegesTokenCollectionCreator(p.AccountPrivilegesTokenCollectionCreator); err != nil {
		return err
	}

	if err := validateAccountPrivilegesTokenCollectionId(p.AccountPrivilegesTokenCollectionId); err != nil {
		return err
	}

	if err := validateDefaultAccountPrivileges(p.DefaultAccountPrivileges); err != nil {
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

func validateAccountPrivilegesTokenCollectionCreator(i interface{}) error {
	_, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateAccountPrivilegesTokenCollectionId(i interface{}) error {
	_, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateDefaultAccountPrivileges(i interface{}) error {
	_, ok := i.([]byte)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
