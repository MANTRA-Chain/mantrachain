package types

import (
	"fmt"
	"regexp"

	"github.com/cosmos/cosmos-sdk/codec"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	DefaultAdminAccount                           string = ""
	DefaultValidBridgeId                          string = "^[a-zA-Z0-9_/:-]{0,100}$"
	DefaultValidBridgeCw20ContractNameMinLength   int32  = 2
	DefaultValidBridgeCw20ContractNameMaxLength   int32  = 100
	DefaultValidBridgeCw20ContractSymbolMinLength int32  = 2
	DefaultValidBridgeCw20ContractSymbolMaxLength int32  = 6
	DefaultValidMintListMetadataMaxCount          int32  = 10
)

var (
	KeyAdminAccount                           = []byte("AdminAccount")
	KeyValidBridgeId                          = []byte("ValidBridgeId")
	KeyValidBridgeCw20ContractNameMinLength   = []byte("ValidBridgeCw20ContractNameMinLength")
	KeyValidBridgeCw20ContractNameMaxLength   = []byte("ValidBridgeCw20ContractNameMaxLength")
	KeyValidBridgeCw20ContractSymbolMinLength = []byte("ValidBridgeCw20ContractSymbolMinLength")
	KeyValidBridgeCw20ContractSymbolMaxLength = []byte("ValidBridgeCw20ContractSymbolMaxLength")
	KeyValidMintListMetadataMaxCount          = []byte("ValidMintListMetadataMaxCount")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	adminAccount string,
	validBridgeId string,
	validBridgeCw20ContractNameMinLength int32,
	validBridgeCw20ContractNameMaxLength int32,
	validBridgeCw20ContractSymbolMinLength int32,
	validBridgeCw20ContractSymbolMaxLength int32,
	validMintListMetadataMaxCount int32,
) Params {
	return Params{
		AdminAccount:                           adminAccount,
		ValidBridgeId:                          validBridgeId,
		ValidBridgeCw20ContractNameMinLength:   validBridgeCw20ContractNameMinLength,
		ValidBridgeCw20ContractNameMaxLength:   validBridgeCw20ContractNameMaxLength,
		ValidBridgeCw20ContractSymbolMinLength: validBridgeCw20ContractSymbolMinLength,
		ValidBridgeCw20ContractSymbolMaxLength: validBridgeCw20ContractSymbolMaxLength,
		ValidMintListMetadataMaxCount:          validMintListMetadataMaxCount,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultAdminAccount,
		DefaultValidBridgeId,
		DefaultValidBridgeCw20ContractNameMinLength,
		DefaultValidBridgeCw20ContractNameMaxLength,
		DefaultValidBridgeCw20ContractSymbolMinLength,
		DefaultValidBridgeCw20ContractSymbolMaxLength,
		DefaultValidMintListMetadataMaxCount,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyAdminAccount, &p.AdminAccount, validateAdminAccount),
		paramtypes.NewParamSetPair(KeyValidBridgeId, &p.ValidBridgeId, validateValidBridgeId),
		paramtypes.NewParamSetPair(KeyValidBridgeCw20ContractNameMinLength, &p.ValidBridgeCw20ContractNameMinLength, validateValidBridgeCw20ContractNameMinLength),
		paramtypes.NewParamSetPair(KeyValidBridgeCw20ContractNameMaxLength, &p.ValidBridgeCw20ContractNameMaxLength, validateValidBridgeCw20ContractNameMaxLength),
		paramtypes.NewParamSetPair(KeyValidBridgeCw20ContractSymbolMinLength, &p.ValidBridgeCw20ContractSymbolMinLength, validateValidBridgeCw20ContractSymbolMinLength),
		paramtypes.NewParamSetPair(KeyValidBridgeCw20ContractSymbolMaxLength, &p.ValidBridgeCw20ContractSymbolMaxLength, validateValidBridgeCw20ContractSymbolMaxLength),
		paramtypes.NewParamSetPair(KeyValidMintListMetadataMaxCount, &p.ValidMintListMetadataMaxCount, validateValidMintListMetadataMaxCount),
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

	if err := validateValidBridgeId(p.ValidBridgeId); err != nil {
		return err
	}

	if err := validateValidBridgeCw20ContractNameMinLength(p.ValidBridgeCw20ContractNameMinLength); err != nil {
		return err
	}

	if err := validateValidBridgeCw20ContractNameMaxLength(p.ValidBridgeCw20ContractNameMaxLength); err != nil {
		return err
	}

	if err := validateValidBridgeCw20ContractSymbolMinLength(p.ValidBridgeCw20ContractSymbolMinLength); err != nil {
		return err
	}

	if err := validateValidBridgeCw20ContractSymbolMaxLength(p.ValidBridgeCw20ContractSymbolMaxLength); err != nil {
		return err
	}

	if err := validateValidMintListMetadataMaxCount(p.ValidMintListMetadataMaxCount); err != nil {
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

func validateValidBridgeId(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == "" {
		return fmt.Errorf("valid bridge id param regex should not be empty")
	}

	_, err := regexp.Compile(v)
	if err != nil {
		return fmt.Errorf("valid bridge id param is invalid regex %s", v)
	}

	return nil
}

func validateValidMintListMetadataMaxCount(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("valid mint list metadata max count param must be positive: %d", v)
	}

	return nil
}

func validateValidBridgeCw20ContractNameMinLength(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("valid bridge cw20 contract name min length param must be positive: %d", v)
	}

	return nil
}

func validateValidBridgeCw20ContractNameMaxLength(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("valid bridge cw20 contract name max length param must be positive: %d", v)
	}

	return nil
}

func validateValidBridgeCw20ContractSymbolMinLength(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("valid bridge cw20 contract symbol min length param must be positive: %d", v)
	}

	return nil
}

func validateValidBridgeCw20ContractSymbolMaxLength(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("valid bridge cw20 contract symbol min length param must be positive: %d", v)
	}

	return nil
}
