package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

const (
	DefaultStakingValidatorAddress = ""
	DefaultStakingValidatorDenom   = "stake"
	DefaultEpochBlockHeightOffset  = int64(5)
	DefaultMinEpochWithdrawAmount  = int64(1)
	DefaultMinRewardWithdrawAmount = int64(1)
)

var (
	KeyStakingValidatorAddress = []byte("StakingValidatorAddress")
	KeyStakingValidatorDenom   = []byte("StakingValidatorDenom")
	KeyEpochBlockHeightOffset  = []byte("EpochBlockHeightOffset")
	KeyMinEpochWithdrawAmount  = []byte("MinEpochWithdrawAmount")
	KeyMinRewardWithdrawAmount = []byte("MinRewardWithdrawAmount")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	stakingValidatorAddress string,
	stakingValidatorDenom string,
	epochBlockHeightOffset int64,
	minEpochWithdrawAmount int64,
	minRewardWithdrawAmount int64,
) Params {
	return Params{
		StakingValidatorAddress: stakingValidatorAddress,
		StakingValidatorDenom:   stakingValidatorDenom,
		EpochBlockHeightOffset:  epochBlockHeightOffset,
		MinEpochWithdrawAmount:  minEpochWithdrawAmount,
		MinRewardWithdrawAmount: minRewardWithdrawAmount,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultStakingValidatorAddress,
		DefaultStakingValidatorDenom,
		DefaultEpochBlockHeightOffset,
		DefaultMinEpochWithdrawAmount,
		DefaultMinRewardWithdrawAmount,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyStakingValidatorAddress, &p.StakingValidatorAddress, validateStakingValidatorAddress),
		paramtypes.NewParamSetPair(KeyStakingValidatorDenom, &p.StakingValidatorDenom, validateStakingValidatorDenom),
		paramtypes.NewParamSetPair(KeyEpochBlockHeightOffset, &p.EpochBlockHeightOffset, validateEpochBlockHeightOffset),
		paramtypes.NewParamSetPair(KeyMinEpochWithdrawAmount, &p.MinEpochWithdrawAmount, validateMinEpochWithdrawAmount),
		paramtypes.NewParamSetPair(KeyMinRewardWithdrawAmount, &p.MinRewardWithdrawAmount, validateMinRewardWithdrawAmount),
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
	if err := validateStakingValidatorAddress(p.StakingValidatorAddress); err != nil {
		return err
	}
	if err := validateStakingValidatorDenom(p.StakingValidatorDenom); err != nil {
		return err
	}
	if err := validateEpochBlockHeightOffset(p.EpochBlockHeightOffset); err != nil {
		return err
	}
	if err := validateMinEpochWithdrawAmount(p.MinEpochWithdrawAmount); err != nil {
		return err
	}
	if err := validateMinRewardWithdrawAmount(p.MinRewardWithdrawAmount); err != nil {
		return err
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

func validateStakingValidatorDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == "" {
		return fmt.Errorf("valid staking validator denom param should not be empty")
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

func validateMinEpochWithdrawAmount(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("min epoch withdraw amount param should be positive")
	}

	return nil
}

func validateMinRewardWithdrawAmount(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("min reward withdraw amount param should be positive")
	}

	return nil
}
