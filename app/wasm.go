package app

import (
	"context"

	circuitkeeper "cosmossdk.io/x/circuit/keeper"
)

// AllCapabilities returns all capabilities available with the current wasmvm
// See https://github.com/CosmWasm/cosmwasm/blob/main/docs/CAPABILITIES-BUILT-IN.md
// This functionality is going to be moved upstream: https://github.com/CosmWasm/wasmvm/issues/425
func AllCapabilities() []string {
	return []string{
		"iterator",
		"staking",
		"stargate",
		"tokenfactory",
		"cosmwasm_1_1",
		"cosmwasm_1_2",
		"cosmwasm_1_3",
		"cosmwasm_1_4",
		"cosmwasm_2_0",
		"cosmwasm_2_1",
		"cosmwasm_2_2",
	}
}

// wasmCircuitBreaker implements the circuit breaker for messages dispatched by wasm contracts.
// It checks the main circuit breaker first, and then appends 'wasm' to typeUrl for additional check.
type wasmCircuitBreaker struct {
	circuitKeeper circuitkeeper.Keeper
}

func (wcb wasmCircuitBreaker) IsAllowed(ctx context.Context, typeURL string) (bool, error) {
	isAllowed, err := wcb.circuitKeeper.IsAllowed(ctx, typeURL)
	if err != nil {
		return false, err
	}
	if !isAllowed {
		return false, nil // Deny if the main circuit breaker says so
	}

	// check additional static blacklist for wasm
	wasmTypeUrl := "wasm" + typeURL
	return wcb.circuitKeeper.IsAllowed(ctx, wasmTypeUrl)
}
