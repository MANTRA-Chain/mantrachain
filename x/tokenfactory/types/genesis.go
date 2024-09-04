package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:        DefaultParams(),
		FactoryDenoms: []GenesisDenom{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	err := gs.Params.Validate()
	if err != nil {
		return err
	}

	seenDenoms := make(map[string]bool)

	for _, denom := range gs.GetFactoryDenoms() {
		if err := validateDenom(denom, seenDenoms); err != nil {
			return err
		}
	}

	return nil
}

// validateDenom validates a single denom in the genesis state.
func validateDenom(denom GenesisDenom, seenDenoms map[string]bool) error {
	switch {
	case seenDenoms[denom.GetDenom()]:
		return errorsmod.Wrapf(ErrInvalidGenesis, "duplicate denom: %s", denom.GetDenom())
	case denom.AuthorityMetadata.Admin != "":
		if _, err := sdk.AccAddressFromBech32(denom.AuthorityMetadata.Admin); err != nil {
			return errorsmod.Wrapf(ErrInvalidAuthorityMetadata, "invalid admin address: %s", err)
		}
	case denom.HookContractAddress != "":
		if _, err := sdk.AccAddressFromBech32(denom.HookContractAddress); err != nil {
			return errorsmod.Wrapf(ErrInvalidHookContractAddress, "invalid hook contract address: %s", err)
		}
	}

	seenDenoms[denom.GetDenom()] = true

	if _, _, err := DeconstructDenom(denom.GetDenom()); err != nil {
		return err
	}

	return nil
}
