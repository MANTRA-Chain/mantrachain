package types

import (
	"fmt"
	"math/big"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	DefaultAdminAccount                            = ""
	DefaultAccountPrivilegesTokenCollectionCreator = ""
	DefaultAccountPrivilegesTokenCollectionId      = ""
	DefaultPrivileges                              = append(make([]byte, 24), big.NewInt(0).Sub(
		big.NewInt(0).
			Exp(
				big.NewInt(2),
				big.NewInt(64),
				nil,
			),
		big.NewInt(1),
	).Bytes()[:]...)
)

var (
	KeyAdminAccount                            = []byte("AdminAccount")
	KeyAccountPrivilegesTokenCollectionCreator = []byte("AccountPrivilegesTokenCollectionCreator")
	KeyAccountPrivilegesTokenCollectionId      = []byte("AccountPrivilegesTokenCollectionId")
	KeyDefaultPrivileges                       = []byte("DefaultPrivileges")
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
	defaultPrivileges []byte,
) Params {
	return Params{
		AdminAccount:                            adminAccount,
		AccountPrivilegesTokenCollectionCreator: accountPrivilegesTokenCollectionCreator,
		AccountPrivilegesTokenCollectionId:      accountPrivilegesTokenCollectionId,
		DefaultPrivileges:                       defaultPrivileges,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultAdminAccount,
		DefaultAccountPrivilegesTokenCollectionCreator,
		DefaultAccountPrivilegesTokenCollectionId,
		DefaultPrivileges,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyAdminAccount, &p.AdminAccount, validateAdminAccount),
		paramtypes.NewParamSetPair(KeyAccountPrivilegesTokenCollectionCreator, &p.AccountPrivilegesTokenCollectionCreator, validateAccountPrivilegesTokenCollectionCreator),
		paramtypes.NewParamSetPair(KeyAccountPrivilegesTokenCollectionId, &p.AccountPrivilegesTokenCollectionId, validateAccountPrivilegesTokenCollectionId),
		paramtypes.NewParamSetPair(KeyDefaultPrivileges, &p.DefaultPrivileges, validateDefaultPrivileges),
	}
}

func MustUnmarshalParams(cdc *codec.LegacyAmino, value []byte) Params {
	params, err := UnmarshalParams(cdc, value)
	if err != nil {
		panic(err)
	}

	return params
}

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

	if err := validateDefaultPrivileges(p.DefaultPrivileges); err != nil {
		return err
	}

	return nil
}

func validateAdminAccount(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	_, err := sdk.AccAddressFromBech32(v)
	if err != nil {
		return fmt.Errorf("invalid account address (%s)", err)
	}

	return nil
}

func validateAccountPrivilegesTokenCollectionCreator(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == "" {
		return fmt.Errorf("valid account privileges token collection creator param should not be empty")
	}

	return nil
}

func validateAccountPrivilegesTokenCollectionId(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == "" {
		return fmt.Errorf("valid account privileges token collection id param should not be empty")
	}

	return nil
}

func validateDefaultPrivileges(i interface{}) error {
	v, ok := i.([]byte)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if len(v) != 32 {
		return fmt.Errorf("valid default privileges param should have length of 32")
	}

	return nil
}
