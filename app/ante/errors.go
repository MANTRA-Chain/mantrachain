package ante

import (
	"errors"
	"fmt"
)

var (
	ErrMissingAccountKeeper    = errors.New("account keeper is required for ante builder")
	ErrMissingBankKeeper       = errors.New("bank keeper is required for ante builder")
	ErrMissingSignModeHandler  = errors.New("sign mode handler is required for ante builder")
	ErrMissingWasmConfig       = errors.New("wasm config is required for ante builder")
	ErrMissingWasmStoreService = errors.New("wasm store service is required for ante builder")
	ErrMissingCircuitKeeper    = errors.New("circuit keeper is required for ante builder")
)

func ErrNeitherNativeDenom(coinDenom, denom string) error {
	return fmt.Errorf("neither of coin.Denom %s and denom %s is the native denom of the chain", coinDenom, denom)
}

func ErrDenomNotRegistered(denom string) error {
	return fmt.Errorf("denom %s not registered in host zone", denom)
}

func ErrExpectedOneCoin(count int) error {
	return fmt.Errorf("expected exactly one native coin, got %d", count)
}
