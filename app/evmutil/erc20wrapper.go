package evmutil

import "math/big"

// Helpers for interacting with ERC20 wrappers (underlying/withdrawTo).
var (
	MinWithdrawAmountWad = big.NewInt(1_000_000_000_000) // 1e12

	SelectorUnderlyingGetter         = [4]byte{0x6f, 0x30, 0x7d, 0xc3} // underlying()
	SelectorWithdrawToAddressUint256 = [4]byte{0x20, 0x5c, 0x28, 0x78} // withdrawTo(address,uint256)
)
