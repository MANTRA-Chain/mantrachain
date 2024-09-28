package types

import (
	"fmt"
	"strings"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgUpdateParams) ValidateBasic() error {
	// Validate Authority address
	if msg.Authority == "" {
		return fmt.Errorf("Authority address cannot be empty")
	}
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return fmt.Errorf("invalid Authority address: %w", err)
	}

	// Validate McaTax
	if msg.McaTax != "" {
		mcaTax, err := math.LegacyNewDecFromStr(msg.McaTax)
		if err != nil {
			return fmt.Errorf("invalid mca tax: %w", err)
		}
		if mcaTax.IsNegative() {
			return fmt.Errorf("mca tax cannot be negative")
		}
		if mcaTax.GT(MaxMcaTax) {
			return fmt.Errorf("mca tax %s cannot exceed maximum of %s", mcaTax, MaxMcaTax)
		}
	}

	// Validate McaAddress
	if msg.McaAddress != "" {
		_, err = sdk.AccAddressFromBech32(msg.McaAddress)
		if err != nil {
			return fmt.Errorf("invalid mca address: %w", err)
		}
		if !strings.HasPrefix(msg.McaAddress, "mantra") {
			return fmt.Errorf("mca address must have 'mantra' prefix")
		}
	}

	return nil
}
