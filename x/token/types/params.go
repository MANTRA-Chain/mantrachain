package types

import (
	"fmt"
	"regexp"

	"github.com/cosmos/cosmos-sdk/codec"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	DefaultValidNftCollectionId                               string = "^[a-zA-Z0-9_/:-]{0,100}$"
	DefaultNftCollectionDefaultId                             string = "default"
	DefaultNftCollectionDefaultName                           string = "default"
	DefaultValidNftCollectionMetadataSymbolMinLength          int32  = 2
	DefaultValidNftCollectionMetadataSymbolMaxLength          int32  = 5
	DefaultValidNftCollectionMetadataDescriptionMaxLength     int32  = 1000
	DefaultValidNftCollectionMetadataNameMaxLength            int32  = 100
	DefaultValidNftCollectionMetadataImagesMaxCount           int32  = 10
	DefaultValidNftCollectionMetadataImagesTypeMaxLength      int32  = 25
	DefaultValidNftCollectionMetadataLinksMaxCount            int32  = 10
	DefaultValidNftCollectionMetadataLinksTypeMaxLength       int32  = 25
	DefaultValidNftCollectionMetadataOptionsMaxCount          int32  = 15
	DefaultValidNftCollectionMetadataOptionsTypeMaxLength     int32  = 25
	DefaultValidNftCollectionMetadataOptionsValueMaxLength    int32  = 25
	DefaultValidNftCollectionMetadataOptionsSubValueMaxLength int32  = 50

	DefaultValidNftId                                  string = "^[a-zA-Z0-9_/:-]{0,100}$"
	DefaultValidNftMetadataMaxCount                    int32  = 100
	DefaultValidNftMetadataTitleMaxLength              int32  = 100
	DefaultValidNftMetadataDescriptionMaxLength        int32  = 1000
	DefaultValidNftMetadataImagesMaxCount              int32  = 15
	DefaultValidNftMetadataImagesTypeMaxLength         int32  = 25
	DefaultValidNftMetadataLinksMaxCount               int32  = 15
	DefaultValidNftMetadataLinksTypeMaxLength          int32  = 25
	DefaultValidNftMetadataAttributesMaxCount          int32  = 300
	DefaultValidNftMetadataAttributesTypeMaxLength     int32  = 25
	DefaultValidNftMetadataAttributesValueMaxLength    int32  = 25
	DefaultValidNftMetadataAttributesSubValueMaxLength int32  = 50
	DefaultValidBurnNftMaxCount                        int32  = 50
)

var (
	KeyValidNftCollectionId                               = []byte("ValidNftCollectionId")
	KeyNftCollectionDefaultId                             = []byte("NftCollectionDefaultId")
	KeyNftCollectionDefaultName                           = []byte("NftCollectionDefaultName")
	KeyValidNftCollectionMetadataSymbolMinLength          = []byte("ValidNftCollectionMetadataSymbolMinLength")
	KeyValidNftCollectionMetadataSymbolMaxLength          = []byte("ValidNftCollectionMetadataSymbolMaxLength")
	KeyValidNftCollectionMetadataDescriptionMaxLength     = []byte("ValidNftCollectionMetadataDescriptionMaxLength")
	KeyValidNftCollectionMetadataNameMaxLength            = []byte("ValidNftCollectionMetadataNameMaxLength")
	KeyValidNftCollectionMetadataImagesMaxCount           = []byte("ValidNftCollectionMetadataImagesMaxCount")
	KeyValidNftCollectionMetadataImagesTypeMaxLength      = []byte("ValidNftCollectionMetadataImagesTypeMaxLength")
	KeyValidNftCollectionMetadataLinksMaxCount            = []byte("ValidNftCollectionMetadataLinksMaxCount")
	KeyValidNftCollectionMetadataLinksTypeMaxLength       = []byte("ValidNftCollectionMetadataLinksTypeMaxLength")
	KeyValidNftCollectionMetadataOptionsMaxCount          = []byte("ValidNftCollectionMetadataOptionsMaxCount")
	KeyValidNftCollectionMetadataOptionsTypeMaxLength     = []byte("ValidNftCollectionMetadataOptionsTypeMaxLength")
	KeyValidNftCollectionMetadataOptionsValueMaxLength    = []byte("ValidNftCollectionMetadataOptionsValueMaxLength")
	KeyValidNftCollectionMetadataOptionsSubValueMaxLength = []byte("ValidNftCollectionMetadataOptionsSubValueMaxLength")

	KeyValidNftId                                  = []byte("ValidNftId")
	KeyValidNftMetadataMaxCount                    = []byte("ValidNftMetadataMaxCount")
	KeyValidNftMetadataTitleMaxLength              = []byte("ValidNftMetadataTitleMaxLength")
	KeyValidNftMetadataDescriptionMaxLength        = []byte("ValidNftMetadataDescriptionMaxLength")
	KeyValidNftMetadataImagesMaxCount              = []byte("ValidNftMetadataImagesMaxCount")
	KeyValidNftMetadataImagesTypeMaxLength         = []byte("ValidNftMetadataImagesTypeMaxLength")
	KeyValidNftMetadataLinksMaxCount               = []byte("ValidNftMetadataLinksMaxCount")
	KeyValidNftMetadataLinksTypeMaxLength          = []byte("ValidNftMetadataLinksTypeMaxLength")
	KeyValidNftMetadataAttributesMaxCount          = []byte("ValidNftMetadataAttributesMaxCount")
	KeyValidNftMetadataAttributesTypeMaxLength     = []byte("ValidNftMetadataAttributesTypeMaxLength")
	KeyValidNftMetadataAttributesValueMaxLength    = []byte("ValidNftMetadataAttributesValueMaxLength")
	KeyValidNftMetadataAttributesSubValueMaxLength = []byte("ValidNftMetadataAttributesSubValueMaxLength")
	KeyValidBurnNftMaxCount                        = []byte("ValidBurnNftMaxCount")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	validNftCollectionId string,
	nftCollectionDefaultId string,
	nftCollectionDefaultName string,
	validNftCollectionMetadataSymbolMinLength int32,
	validNftCollectionMetadataSymbolMaxLength int32,
	validNftCollectionMetadataDescriptionMaxLength int32,
	validNftCollectionMetadataNameMaxLength int32,
	validNftCollectionMetadataImagesMaxCount int32,
	validNftCollectionMetadataImagesTypeMaxLength int32,
	validNftCollectionMetadataLinksMaxCount int32,
	validNftCollectionMetadataLinksTypeMaxLength int32,
	validNftCollectionMetadataOptionsMaxCount int32,
	validNftCollectionMetadataOptionsTypeMaxLength int32,
	validNftCollectionMetadataOptionsValueMaxLength int32,
	validNftCollectionMetadataOptionsSubValueMaxLength int32,

	validNftId string,
	validNftMetadataMaxCount int32,
	validNftMetadataTitleMaxLength int32,
	validNftMetadataDescriptionMaxLength int32,
	validNftMetadataImagesMaxCount int32,
	validNftMetadataImagesTypeMaxLength int32,
	validNftMetadataLinksMaxCount int32,
	validNftMetadataLinksTypeMaxLength int32,
	validNftMetadataAttributesMaxCount int32,
	validNftMetadataAttributesTypeMaxLength int32,
	validNftMetadataAttributesValueMaxLength int32,
	validNftMetadataAttributesSubValueMaxLength int32,
	validBurnNftMaxCount int32,
) Params {
	return Params{
		ValidNftCollectionId:                               validNftCollectionId,
		NftCollectionDefaultId:                             nftCollectionDefaultId,
		NftCollectionDefaultName:                           nftCollectionDefaultName,
		ValidNftCollectionMetadataSymbolMinLength:          validNftCollectionMetadataSymbolMinLength,
		ValidNftCollectionMetadataSymbolMaxLength:          validNftCollectionMetadataSymbolMaxLength,
		ValidNftCollectionMetadataDescriptionMaxLength:     validNftCollectionMetadataDescriptionMaxLength,
		ValidNftCollectionMetadataNameMaxLength:            validNftCollectionMetadataNameMaxLength,
		ValidNftCollectionMetadataImagesMaxCount:           validNftCollectionMetadataImagesMaxCount,
		ValidNftCollectionMetadataImagesTypeMaxLength:      validNftCollectionMetadataImagesTypeMaxLength,
		ValidNftCollectionMetadataLinksMaxCount:            validNftCollectionMetadataLinksMaxCount,
		ValidNftCollectionMetadataLinksTypeMaxLength:       validNftCollectionMetadataLinksTypeMaxLength,
		ValidNftCollectionMetadataOptionsMaxCount:          validNftCollectionMetadataOptionsMaxCount,
		ValidNftCollectionMetadataOptionsTypeMaxLength:     validNftCollectionMetadataOptionsTypeMaxLength,
		ValidNftCollectionMetadataOptionsValueMaxLength:    validNftCollectionMetadataOptionsValueMaxLength,
		ValidNftCollectionMetadataOptionsSubValueMaxLength: validNftCollectionMetadataOptionsSubValueMaxLength,

		ValidNftId:                                  validNftId,
		ValidNftMetadataMaxCount:                    validNftMetadataMaxCount,
		ValidNftMetadataTitleMaxLength:              validNftMetadataTitleMaxLength,
		ValidNftMetadataDescriptionMaxLength:        validNftMetadataDescriptionMaxLength,
		ValidNftMetadataImagesMaxCount:              validNftMetadataImagesMaxCount,
		ValidNftMetadataImagesTypeMaxLength:         validNftMetadataImagesTypeMaxLength,
		ValidNftMetadataLinksMaxCount:               validNftMetadataLinksMaxCount,
		ValidNftMetadataLinksTypeMaxLength:          validNftMetadataLinksTypeMaxLength,
		ValidNftMetadataAttributesMaxCount:          validNftMetadataAttributesMaxCount,
		ValidNftMetadataAttributesTypeMaxLength:     validNftMetadataAttributesTypeMaxLength,
		ValidNftMetadataAttributesValueMaxLength:    validNftMetadataAttributesValueMaxLength,
		ValidNftMetadataAttributesSubValueMaxLength: validNftMetadataAttributesSubValueMaxLength,
		ValidBurnNftMaxCount:                        validBurnNftMaxCount,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultValidNftCollectionId,
		DefaultNftCollectionDefaultId,
		DefaultNftCollectionDefaultName,
		DefaultValidNftCollectionMetadataSymbolMinLength,
		DefaultValidNftCollectionMetadataSymbolMaxLength,
		DefaultValidNftCollectionMetadataDescriptionMaxLength,
		DefaultValidNftCollectionMetadataNameMaxLength,
		DefaultValidNftCollectionMetadataImagesMaxCount,
		DefaultValidNftCollectionMetadataImagesTypeMaxLength,
		DefaultValidNftCollectionMetadataLinksMaxCount,
		DefaultValidNftCollectionMetadataLinksTypeMaxLength,
		DefaultValidNftCollectionMetadataOptionsMaxCount,
		DefaultValidNftCollectionMetadataOptionsTypeMaxLength,
		DefaultValidNftCollectionMetadataOptionsValueMaxLength,
		DefaultValidNftCollectionMetadataOptionsSubValueMaxLength,

		DefaultValidNftId,
		DefaultValidNftMetadataMaxCount,
		DefaultValidNftMetadataTitleMaxLength,
		DefaultValidNftMetadataDescriptionMaxLength,
		DefaultValidNftMetadataImagesMaxCount,
		DefaultValidNftMetadataImagesTypeMaxLength,
		DefaultValidNftMetadataLinksMaxCount,
		DefaultValidNftMetadataLinksTypeMaxLength,
		DefaultValidNftMetadataAttributesMaxCount,
		DefaultValidNftMetadataAttributesTypeMaxLength,
		DefaultValidNftMetadataAttributesValueMaxLength,
		DefaultValidNftMetadataAttributesSubValueMaxLength,
		DefaultValidBurnNftMaxCount,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyValidNftCollectionId, &p.ValidNftCollectionId, validateValidNftCollectionId),
		paramtypes.NewParamSetPair(KeyNftCollectionDefaultId, &p.NftCollectionDefaultId, validateNftCollectionDefaultId),
		paramtypes.NewParamSetPair(KeyNftCollectionDefaultName, &p.NftCollectionDefaultName, validateNftCollectionDefaultName),
		paramtypes.NewParamSetPair(KeyValidNftCollectionMetadataSymbolMinLength, &p.ValidNftCollectionMetadataSymbolMinLength, validateValidNftCollectionMetadataSymbolMinLength),
		paramtypes.NewParamSetPair(KeyValidNftCollectionMetadataSymbolMaxLength, &p.ValidNftCollectionMetadataSymbolMaxLength, validateValidNftCollectionMetadataSymbolMaxLength),
		paramtypes.NewParamSetPair(KeyValidNftCollectionMetadataDescriptionMaxLength, &p.ValidNftCollectionMetadataDescriptionMaxLength, validateValidNftCollectionMetadataDescriptionMaxLength),
		paramtypes.NewParamSetPair(KeyValidNftCollectionMetadataNameMaxLength, &p.ValidNftCollectionMetadataNameMaxLength, validateValidNftCollectionMetadataNameMaxLength),
		paramtypes.NewParamSetPair(KeyValidNftCollectionMetadataImagesMaxCount, &p.ValidNftCollectionMetadataImagesMaxCount, validateValidNftCollectionMetadataImagesMaxCount),
		paramtypes.NewParamSetPair(KeyValidNftCollectionMetadataImagesTypeMaxLength, &p.ValidNftCollectionMetadataImagesTypeMaxLength, validateValidNftCollectionMetadataImagesTypeMaxLength),
		paramtypes.NewParamSetPair(KeyValidNftCollectionMetadataLinksMaxCount, &p.ValidNftCollectionMetadataLinksMaxCount, validateValidNftCollectionMetadataLinksMaxCount),
		paramtypes.NewParamSetPair(KeyValidNftCollectionMetadataLinksTypeMaxLength, &p.ValidNftCollectionMetadataLinksTypeMaxLength, validateValidNftCollectionMetadataLinksTypeMaxLength),
		paramtypes.NewParamSetPair(KeyValidNftCollectionMetadataOptionsMaxCount, &p.ValidNftCollectionMetadataOptionsMaxCount, validateValidNftCollectionMetadataOptionsMaxCount),
		paramtypes.NewParamSetPair(KeyValidNftCollectionMetadataOptionsTypeMaxLength, &p.ValidNftCollectionMetadataOptionsTypeMaxLength, validateValidNftCollectionMetadataOptionsTypeMaxLength),
		paramtypes.NewParamSetPair(KeyValidNftCollectionMetadataOptionsValueMaxLength, &p.ValidNftCollectionMetadataOptionsValueMaxLength, validateValidNftCollectionMetadataOptionsValueMaxLength),
		paramtypes.NewParamSetPair(KeyValidNftCollectionMetadataOptionsSubValueMaxLength, &p.ValidNftCollectionMetadataOptionsSubValueMaxLength, validateValidNftCollectionMetadataOptionsSubValueMaxLength),

		paramtypes.NewParamSetPair(KeyValidNftId, &p.ValidNftId, validateValidNftId),
		paramtypes.NewParamSetPair(KeyValidNftMetadataMaxCount, &p.ValidNftMetadataMaxCount, validateValidNftMetadataMaxCount),
		paramtypes.NewParamSetPair(KeyValidNftMetadataTitleMaxLength, &p.ValidNftMetadataTitleMaxLength, validateValidNftMetadataTitleMaxLength),
		paramtypes.NewParamSetPair(KeyValidNftMetadataDescriptionMaxLength, &p.ValidNftMetadataDescriptionMaxLength, validateValidNftMetadataDescriptionMaxLength),
		paramtypes.NewParamSetPair(KeyValidNftMetadataImagesMaxCount, &p.ValidNftMetadataImagesMaxCount, validateValidNftMetadataImagesMaxCount),
		paramtypes.NewParamSetPair(KeyValidNftMetadataImagesTypeMaxLength, &p.ValidNftMetadataImagesTypeMaxLength, validateValidNftMetadataImagesTypeMaxLength),
		paramtypes.NewParamSetPair(KeyValidNftMetadataLinksMaxCount, &p.ValidNftMetadataLinksMaxCount, validateValidNftMetadataLinksMaxCount),
		paramtypes.NewParamSetPair(KeyValidNftMetadataLinksTypeMaxLength, &p.ValidNftMetadataLinksTypeMaxLength, validateValidNftMetadataLinksTypeMaxLength),
		paramtypes.NewParamSetPair(KeyValidNftMetadataAttributesMaxCount, &p.ValidNftMetadataAttributesMaxCount, validateValidNftMetadataAttributesMaxCount),
		paramtypes.NewParamSetPair(KeyValidNftMetadataAttributesTypeMaxLength, &p.ValidNftMetadataAttributesTypeMaxLength, validateValidNftMetadataAttributesTypeMaxLength),
		paramtypes.NewParamSetPair(KeyValidNftMetadataAttributesValueMaxLength, &p.ValidNftMetadataAttributesValueMaxLength, validateValidNftMetadataAttributesValueMaxLength),
		paramtypes.NewParamSetPair(KeyValidNftMetadataAttributesSubValueMaxLength, &p.ValidNftMetadataAttributesSubValueMaxLength, validateValidNftMetadataAttributesSubValueMaxLength),
		paramtypes.NewParamSetPair(KeyValidBurnNftMaxCount, &p.ValidBurnNftMaxCount, validateValidBurnNftMaxCount),
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
	if err := validateValidNftCollectionId(p.ValidNftCollectionId); err != nil {
		return err
	}

	if err := validateNftCollectionDefaultId(p.NftCollectionDefaultId); err != nil {
		return err
	}

	if err := validateNftCollectionDefaultName(p.NftCollectionDefaultName); err != nil {
		return err
	}

	if err := validateValidNftCollectionMetadataSymbolMinLength(p.ValidNftCollectionMetadataSymbolMinLength); err != nil {
		return err
	}

	if err := validateValidNftCollectionMetadataSymbolMaxLength(p.ValidNftCollectionMetadataSymbolMaxLength); err != nil {
		return err
	}

	if err := validateValidNftCollectionMetadataDescriptionMaxLength(p.ValidNftCollectionMetadataDescriptionMaxLength); err != nil {
		return err
	}

	if err := validateValidNftCollectionMetadataNameMaxLength(p.ValidNftCollectionMetadataNameMaxLength); err != nil {
		return err
	}

	if err := validateValidNftCollectionMetadataImagesMaxCount(p.ValidNftCollectionMetadataImagesMaxCount); err != nil {
		return err
	}

	if err := validateValidNftCollectionMetadataImagesTypeMaxLength(p.ValidNftCollectionMetadataImagesTypeMaxLength); err != nil {
		return err
	}

	if err := validateValidNftCollectionMetadataLinksMaxCount(p.ValidNftCollectionMetadataLinksMaxCount); err != nil {
		return err
	}

	if err := validateValidNftCollectionMetadataLinksTypeMaxLength(p.ValidNftCollectionMetadataLinksTypeMaxLength); err != nil {
		return err
	}

	if err := validateValidNftCollectionMetadataOptionsMaxCount(p.ValidNftCollectionMetadataOptionsMaxCount); err != nil {
		return err
	}

	if err := validateValidNftCollectionMetadataOptionsTypeMaxLength(p.ValidNftCollectionMetadataOptionsTypeMaxLength); err != nil {
		return err
	}

	if err := validateValidNftCollectionMetadataOptionsValueMaxLength(p.ValidNftCollectionMetadataOptionsValueMaxLength); err != nil {
		return err
	}

	if err := validateValidNftCollectionMetadataOptionsSubValueMaxLength(p.ValidNftCollectionMetadataOptionsSubValueMaxLength); err != nil {
		return err
	}

	if err := validateValidNftId(p.ValidNftId); err != nil {
		return err
	}

	if err := validateValidNftMetadataMaxCount(p.ValidNftMetadataMaxCount); err != nil {
		return err
	}

	if err := validateValidNftMetadataTitleMaxLength(p.ValidNftMetadataTitleMaxLength); err != nil {
		return err
	}

	if err := validateValidNftMetadataDescriptionMaxLength(p.ValidNftMetadataDescriptionMaxLength); err != nil {
		return err
	}

	if err := validateValidNftMetadataImagesMaxCount(p.ValidNftMetadataImagesMaxCount); err != nil {
		return err
	}

	if err := validateValidNftMetadataImagesTypeMaxLength(p.ValidNftMetadataImagesTypeMaxLength); err != nil {
		return err
	}

	if err := validateValidNftMetadataLinksMaxCount(p.ValidNftMetadataLinksMaxCount); err != nil {
		return err
	}

	if err := validateValidNftMetadataLinksTypeMaxLength(p.ValidNftMetadataLinksTypeMaxLength); err != nil {
		return err
	}

	if err := validateValidNftMetadataAttributesMaxCount(p.ValidNftMetadataAttributesMaxCount); err != nil {
		return err
	}

	if err := validateValidNftMetadataAttributesTypeMaxLength(p.ValidNftMetadataAttributesTypeMaxLength); err != nil {
		return err
	}

	if err := validateValidNftMetadataAttributesValueMaxLength(p.ValidNftMetadataAttributesValueMaxLength); err != nil {
		return err
	}

	if err := validateValidNftMetadataAttributesSubValueMaxLength(p.ValidNftMetadataAttributesSubValueMaxLength); err != nil {
		return err
	}

	if err := validateValidBurnNftMaxCount(p.ValidBurnNftMaxCount); err != nil {
		return err
	}

	return nil
}

func validateValidNftCollectionId(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == "" {
		return fmt.Errorf("valid nft collection id param regex should not be empty")
	}

	_, err := regexp.Compile(v)
	if err != nil {
		return fmt.Errorf("valid nft collection id param is invalid regex %s", v)
	}

	return nil
}

func validateNftCollectionDefaultId(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == "" {
		return fmt.Errorf("valid nft collection default id param should not be empty")
	}

	return nil
}

func validateNftCollectionDefaultName(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == "" {
		return fmt.Errorf("valid nft collection default name param should not be empty")
	}

	return nil
}

func validateValidNftCollectionMetadataSymbolMinLength(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("valid nft collection metadata symbol min length param must be positive: %d", v)
	}

	return nil
}

func validateValidNftCollectionMetadataSymbolMaxLength(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("valid nft collection metadata symbol max length param must be positive: %d", v)
	}

	return nil
}

func validateValidNftCollectionMetadataDescriptionMaxLength(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("valid nft collection metadata description max length param must be positive: %d", v)
	}

	return nil
}

func validateValidNftCollectionMetadataNameMaxLength(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("valid nft collection metadata name max length param must be positive: %d", v)
	}

	return nil
}

func validateValidNftCollectionMetadataImagesMaxCount(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("valid nft collection metadata images max count param must be positive: %d", v)
	}

	return nil
}

func validateValidNftCollectionMetadataImagesTypeMaxLength(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("valid nft collection metadata images type max length param must be positive: %d", v)
	}

	return nil
}

func validateValidNftCollectionMetadataLinksMaxCount(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("valid nft collection metadata links max count param must be positive: %d", v)
	}

	return nil
}

func validateValidNftCollectionMetadataLinksTypeMaxLength(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("valid nft collection metadata links type max length param must be positive: %d", v)
	}

	return nil
}

func validateValidNftCollectionMetadataOptionsMaxCount(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("valid nft collection metadata options max count param must be positive: %d", v)
	}

	return nil
}

func validateValidNftCollectionMetadataOptionsTypeMaxLength(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("valid nft collection metadata options type max length param must be positive: %d", v)
	}

	return nil
}

func validateValidNftCollectionMetadataOptionsValueMaxLength(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("valid nft collection metadata options value max length param must be positive: %d", v)
	}

	return nil
}

func validateValidNftCollectionMetadataOptionsSubValueMaxLength(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("valid nft collection metadata options sub value max length param must be positive: %d", v)
	}

	return nil
}

func validateValidNftId(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == "" {
		return fmt.Errorf("valid nft id param regex should not be empty")
	}

	_, err := regexp.Compile(v)
	if err != nil {
		return fmt.Errorf("valid nft id param is invalid regex %s", v)
	}

	return nil
}

func validateValidNftMetadataMaxCount(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("valid nft metadata max count param must be positive: %d", v)
	}

	return nil
}

func validateValidNftMetadataTitleMaxLength(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("valid nft metadata title max length param must be positive: %d", v)
	}

	return nil
}

func validateValidNftMetadataDescriptionMaxLength(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("valid nft metadata description max length param must be positive: %d", v)
	}

	return nil
}

func validateValidNftMetadataImagesMaxCount(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("valid nft metadata images max count param must be positive: %d", v)
	}

	return nil
}

func validateValidNftMetadataImagesTypeMaxLength(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("valid nft metadata images type max length param must be positive: %d", v)
	}

	return nil
}

func validateValidNftMetadataLinksMaxCount(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("valid nft metadata links max count param must be positive: %d", v)
	}

	return nil
}

func validateValidNftMetadataLinksTypeMaxLength(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("valid nft metadata links type max length param must be positive: %d", v)
	}

	return nil
}

func validateValidNftMetadataAttributesMaxCount(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("valid nft metadata attributes max count param must be positive: %d", v)
	}

	return nil
}

func validateValidNftMetadataAttributesTypeMaxLength(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("valid nft metadata attributes type max length param must be positive: %d", v)
	}

	return nil
}

func validateValidNftMetadataAttributesValueMaxLength(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("valid nft metadata attributes value max length param must be positive: %d", v)
	}

	return nil
}

func validateValidNftMetadataAttributesSubValueMaxLength(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("valid nft metadata attributes sub value max length param must be positive: %d", v)
	}

	return nil
}

func validateValidBurnNftMaxCount(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("valid burn nft max count param must be positive: %d", v)
	}

	return nil
}
