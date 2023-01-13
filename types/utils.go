package types

import (
	"strings"
)

const (
	// MainnetChainID defines the Mantrachain EIP155 chain ID for mainnet
	MainnetChainID = "mantrachain_7001"
	// TestnetChainID defines the Mantrachain EIP155 chain ID for testnet
	TestnetChainID = "mantrachain_7000"
)

// IsMainnet returns true if the chain-id has the Mantrachain mainnet EIP155 chain prefix.
func IsMainnet(chainID string) bool {
	return strings.HasPrefix(chainID, MainnetChainID)
}

// IsTestnet returns true if the chain-id has the Mantrachain testnet EIP155 chain prefix.
func IsTestnet(chainID string) bool {
	return strings.HasPrefix(chainID, TestnetChainID)
}
