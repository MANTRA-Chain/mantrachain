package types

import (
	"fmt"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

const (
	DefaultMaxClaimedRangeLength uint64 = 0
	DefaultMaxPurgedRangeLength  uint64 = 0
	DefaultMinDepositTime        uint64 = 0
	DefaultDistributionPeriod    uint64 = 0
	DefaultPairsCycleMaxCount    uint64 = math.MaxUint64
	DefaultMaxSnapshotsCount     uint64 = math.MaxUint64
)

var (
	DefaultDistributionFeeRate = sdk.NewDec(0)
)

var (
	KeyMaxClaimedRangeLength = []byte("MaxClaimedRangeLength")
	KeyMaxPurgedRangeLength  = []byte("MaxPurgedRangeLength")
	KeyMinDepositTime        = []byte("MinDepositTime")
	KeyDistributionPeriod    = []byte("DistributionPeriod")
	KeyPairsCycleMaxCount    = []byte("PairsCycleMaxCount")
	KeyDistributionFeeRate   = []byte("DistributionFeeRate")
	KeyMaxSnapshotsCount     = []byte("MaxSnapshotsCount")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	maxClaimedRangeLength uint64,
	maxPurgedRangeLength uint64,
	minDepositTime uint64,
	distributionPeriod uint64,
	pairsCycleMaxCount uint64,
	distributionFeeRate sdk.Dec,
	maxSnapshotsCount uint64,
) Params {
	return Params{
		MaxClaimedRangeLength: maxClaimedRangeLength,
		MaxPurgedRangeLength:  maxPurgedRangeLength,
		MinDepositTime:        minDepositTime,
		DistributionPeriod:    distributionPeriod,
		PairsCycleMaxCount:    pairsCycleMaxCount,
		DistributionFeeRate:   distributionFeeRate,
		MaxSnapshotsCount:     maxSnapshotsCount,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultMaxClaimedRangeLength,
		DefaultMaxPurgedRangeLength,
		DefaultMinDepositTime,
		DefaultDistributionPeriod,
		DefaultPairsCycleMaxCount,
		DefaultDistributionFeeRate,
		DefaultMaxSnapshotsCount,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMaxClaimedRangeLength, &p.MaxClaimedRangeLength, validateMaxClaimedRangeLength),
		paramtypes.NewParamSetPair(KeyMaxPurgedRangeLength, &p.MaxPurgedRangeLength, validateMaxPurgedRangeLength),
		paramtypes.NewParamSetPair(KeyMinDepositTime, &p.MinDepositTime, validateMinDepositTime),
		paramtypes.NewParamSetPair(KeyDistributionPeriod, &p.DistributionPeriod, validateDistributionPeriod),
		paramtypes.NewParamSetPair(KeyPairsCycleMaxCount, &p.PairsCycleMaxCount, validatePairsCycleMaxCount),
		paramtypes.NewParamSetPair(KeyDistributionFeeRate, &p.DistributionFeeRate, validateDistributionFeeRate),
		paramtypes.NewParamSetPair(KeyMaxSnapshotsCount, &p.MaxSnapshotsCount, validateMaxSnapshotsCount),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateMaxClaimedRangeLength(p.MaxClaimedRangeLength); err != nil {
		return err
	}

	if err := validateMaxPurgedRangeLength(p.MaxPurgedRangeLength); err != nil {
		return err
	}

	if err := validateMinDepositTime(p.MinDepositTime); err != nil {
		return err
	}

	if err := validateDistributionPeriod(p.DistributionPeriod); err != nil {
		return err
	}

	if err := validatePairsCycleMaxCount(p.PairsCycleMaxCount); err != nil {
		return err
	}

	if err := validateDistributionFeeRate(p.DistributionFeeRate); err != nil {
		return err
	}

	if err := validateMaxSnapshotsCount(p.MaxSnapshotsCount); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func validateMaxClaimedRangeLength(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateMaxPurgedRangeLength(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateMinDepositTime(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateDistributionPeriod(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validatePairsCycleMaxCount(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateDistributionFeeRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("distribution fee rate must not be negative: %s", v)
	}

	return nil
}

func validateMaxSnapshotsCount(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
