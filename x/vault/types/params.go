package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

const (
	DefaultAdminAccount            = ""
	DefaultStakingValidatorAddress = ""
	DefaultEpochBlockHeightOffset  = int64(100)
)

var (
	KeyAdminAccount            = []byte("AdminAccount")
	KeyStakingValidatorAddress = []byte("StakingValidatorAddress")
	KeyEpochBlockHeightOffset  = []byte("EpochBlockHeightOffset")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	adminAccount string,
	stakingValidatorAddress string,
	epochBlockHeightOffset int64,
) Params {
	return Params{
		AdminAccount:            adminAccount,
		StakingValidatorAddress: stakingValidatorAddress,
		EpochBlockHeightOffset:  epochBlockHeightOffset,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultAdminAccount,
		DefaultStakingValidatorAddress,
		DefaultEpochBlockHeightOffset,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyAdminAccount, &p.AdminAccount, validateAdminAccount),
		paramtypes.NewParamSetPair(KeyStakingValidatorAddress, &p.StakingValidatorAddress, validateStakingValidatorAddress),
		paramtypes.NewParamSetPair(KeyEpochBlockHeightOffset, &p.EpochBlockHeightOffset, validateEpochBlockHeightOffset),
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
	if err := validateStakingValidatorAddress(p.StakingValidatorAddress); err != nil {
		return err
	}
	if err := validateEpochBlockHeightOffset(p.EpochBlockHeightOffset); err != nil {
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

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func validateStakingValidatorAddress(i interface{}) error {
	_, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateEpochBlockHeightOffset(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("epoch block height offset param should be positive")
	}

	return nil
}
