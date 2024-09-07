package app

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/assert"
)

func TestRegisterIBC(t *testing.T) {
	registry := types.NewInterfaceRegistry()
	modules := RegisterIBC(registry)

	// Test that the correct number of modules are registered
	assert.Len(t, modules, 8, "Expected 8 modules to be registered")

	// Test that specific modules are present
	expectedModules := []string{
		"ibc",
		"transfer",
		"capability",
		"07-tendermint",
		"solomachine",
		"wasm",
	}

	for _, moduleName := range expectedModules {
		_, exists := modules[moduleName]
		assert.True(t, exists, "Expected module %s to be registered", moduleName)
	}

	// Test that the registry has interfaces registered
	// This is a basic check and might need to be adjusted based on the actual interfaces registered
	assert.NotEmpty(t, registry.ListAllInterfaces(), "Expected interfaces to be registered")
}
