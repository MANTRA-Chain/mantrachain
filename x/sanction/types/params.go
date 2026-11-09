package types

import "fmt"

const (
	// DefaultMaxAddressesPerMsg is the default cap on addresses per
	// sanctioner message.
	DefaultMaxAddressesPerMsg uint32 = 50
	// DefaultMaxReasonLength is the default cap, in bytes, on the reason of a
	// sanctioner message.
	DefaultMaxReasonLength uint32 = 256
)

// NewParams creates a new Params instance.
func NewParams(maxAddressesPerMsg, maxReasonLength uint32) Params {
	return Params{
		MaxAddressesPerMsg: maxAddressesPerMsg,
		MaxReasonLength:    maxReasonLength,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(DefaultMaxAddressesPerMsg, DefaultMaxReasonLength)
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if p.MaxAddressesPerMsg == 0 {
		return fmt.Errorf("max_addresses_per_msg must be positive")
	}
	if p.MaxReasonLength == 0 {
		return fmt.Errorf("max_reason_length must be positive")
	}
	return nil
}
