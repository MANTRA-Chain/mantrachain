package app

import (
	"context"

	circuitkeeper "cosmossdk.io/x/circuit/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	erc20types "github.com/cosmos/evm/x/erc20/types"
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

var (
	// wasmBlacklistedMsgs contains the list of messages that are not allowed to be executed by wasm contracts.
	// This is a static blacklist that prevents contracts from executing potentially dangerous messages.
	wasmBlacklistedMsgs = map[string]struct{}{
		sdk.MsgTypeURL(&erc20types.MsgRegisterERC20{}): {},
	}
)

// wasmCircuitBreaker is a custom circuit breaker for wasm contracts.
// It enforces a static blacklist of messages and also respects the dynamic
// list from the main circuit breaker keeper.
type wasmCircuitBreaker struct {
	circuitKeeper circuitkeeper.Keeper
}

func (wcb wasmCircuitBreaker) IsAllowed(ctx context.Context, typeURL string) (bool, error) {
	if _, isBlacklisted := wasmBlacklistedMsgs[typeURL]; isBlacklisted {
		return false, nil // Deny if on the static blacklist
	}
	// Otherwise, defer to the main circuit breaker's list
	return wcb.circuitKeeper.IsAllowed(ctx, typeURL)
}
