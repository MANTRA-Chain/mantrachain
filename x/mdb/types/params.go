package types

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		ValidNftCollectionId:                               "[a-zA-Z][a-zA-Z0-9/:-]{2,100}",
		NftCollectionDefaultId:                             "default",
		NftCollectionDefaultName:                           "default",
		ValidNftCollectionMetadataSymbolMinLength:          2,
		ValidNftCollectionMetadataSymbolMaxLength:          5,
		ValidNftCollectionMetadataDescriptionMaxLength:     1000,
		ValidNftCollectionMetadataNameMaxLength:            100,
		ValidNftCollectionMetadataImagesMaxCount:           10,
		ValidNftCollectionMetadataImagesTypeMaxLength:      25,
		ValidNftCollectionMetadataLinksMaxCount:            10,
		ValidNftCollectionMetadataLinksTypeMaxLength:       25,
		ValidNftCollectionMetadataOptionsMaxCount:          10,
		ValidNftCollectionMetadataOptionsTypeMaxLength:     25,
		ValidNftCollectionMetadataOptionsValueMaxLength:    25,
		ValidNftCollectionMetadataOptionsSubValueMaxLength: 50,

		ValidNftId:                                  "[a-zA-Z][a-zA-Z0-9/:-]{2,100}",
		ValidNftMetadataMaxCount:                    50,
		ValidNftMetadataTitleMaxLength:              100,
		ValidNftMetadataDescriptionMaxLength:        1000,
		ValidNftMetadataImagesMaxCount:              10,
		ValidNftMetadataImagesTypeMaxLength:         25,
		ValidNftMetadataLinksMaxCount:               10,
		ValidNftMetadataLinksTypeMaxLength:          25,
		ValidNftMetadataAttributesMaxCount:          50,
		ValidNftMetadataAttributesTypeMaxLength:     25,
		ValidNftMetadataAttributesValueMaxLength:    25,
		ValidNftMetadataAttributesSubValueMaxLength: 50,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

// Validate validates the set of params
func (p Params) Validate() error {
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
