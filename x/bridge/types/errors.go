package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/bridge module sentinel errors
var (
	ErrInvalidSigner             = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrNoInput                   = errors.Register(ModuleName, 1101, "no input")
	ErrNoOutputs                 = errors.Register(ModuleName, 1102, "no outputs")
	ErrInputOutputsCoinsMismatch = errors.Register(ModuleName, 1103, "input coins and output coins mismatch")
	ErrNoEthTxHashes             = errors.Register(ModuleName, 1104, "no eth tx hashes")
	ErrOutputEthTxHashMismatch   = errors.Register(ModuleName, 1105, "outputs ethTxHashes length mismatch")
	ErrNoInputCoins              = errors.Register(ModuleName, 1106, "no input coins")
	ErrNoInputCoin               = errors.Register(ModuleName, 1107, "no input coin")
	ErrEmptyEthTxHash            = errors.Register(ModuleName, 1108, "empty eth tx hash")
	ErrMultipleOutputCoins       = errors.Register(ModuleName, 1109, "multiple output coins")
)
