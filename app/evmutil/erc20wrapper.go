package evmutil

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/evm/x/ibc/callbacks/types"
	"github.com/cosmos/evm/x/vm/statedb"
)

// Helpers for interacting with ERC20 wrappers (underlying/withdrawTo).
var (
	MinWithdrawAmountWad = big.NewInt(1_000_000_000_000) // 1e12

	GasCapERC20WrapperUnderlying = big.NewInt(200_000)
	GasCapERC20WrapperWithdrawTo = big.NewInt(500_000)

	SelectorUnderlyingGetter         = [4]byte{0x6f, 0x30, 0x7d, 0xc3} // underlying()
	SelectorWithdrawToAddressUint256 = [4]byte{0x20, 0x5c, 0x28, 0x78} // withdrawTo(address,uint256)
)

func ERC20WrapperUnderlyingCallData() []byte {
	data := make([]byte, 4)
	copy(data[:4], SelectorUnderlyingGetter[:])
	return data
}

func ERC20WrapperWithdrawToCallData(to common.Address, amountWad *big.Int) ([]byte, error) {
	if amountWad == nil || amountWad.Sign() <= 0 {
		return nil, fmt.Errorf("invalid withdraw amount")
	}
	u, overflow := uint256.FromBig(amountWad)
	if overflow {
		return nil, fmt.Errorf("withdraw amount overflows uint256")
	}

	data := make([]byte, 4+32+32)
	copy(data[:4], SelectorWithdrawToAddressUint256[:])
	copy(data[4:4+32], common.LeftPadBytes(to.Bytes(), 32))
	amt32 := u.Bytes32()
	copy(data[4+32:], amt32[:])

	return data, nil
}

func DecodeABIAddress32(ret []byte) (common.Address, error) {
	if len(ret) < 32 {
		return common.Address{}, fmt.Errorf("returned short data")
	}
	addr := common.BytesToAddress(ret[12:32])
	if addr == (common.Address{}) {
		return common.Address{}, fmt.Errorf("returned zero address")
	}
	return addr, nil
}

func ERC20WrapperUnderlyingViaEVMCaller(ctx sdk.Context, caller types.EVMKeeper, from common.Address, wrapper common.Address) (common.Address, error) {
	data := ERC20WrapperUnderlyingCallData()
	stateDB := statedb.New(ctx, caller, statedb.NewEmptyTxConfig())
	res, err := caller.CallEVMWithData(ctx, stateDB, from, &wrapper, data, false, false, GasCapERC20WrapperUnderlying)
	if res == nil {
		if err != nil {
			return common.Address{}, WrapERC20WrapperUnderlyingError(err)
		}
		return common.Address{}, WrapERC20WrapperUnderlyingError(fmt.Errorf("nil response from EVM call"))
	}
	if err != nil {
		return common.Address{}, WrapERC20WrapperUnderlyingError(err)
	}
	underlying, err := DecodeABIAddress32(res.Ret)
	if err != nil {
		return common.Address{}, WrapERC20WrapperUnderlyingError(err)
	}
	return underlying, nil
}

func ERC20WrapperWithdrawToViaEVMCaller(ctx sdk.Context, caller types.EVMKeeper, from common.Address, wrapper common.Address, to common.Address, amountWad *big.Int) ([]byte, error) {
	data, err := ERC20WrapperWithdrawToCallData(to, amountWad)
	if err != nil {
		return nil, err
	}
	stateDB := statedb.New(ctx, caller, statedb.NewEmptyTxConfig())
	res, err := caller.CallEVMWithData(ctx, stateDB, from, &wrapper, data, true, false, GasCapERC20WrapperWithdrawTo)
	if res == nil {
		if err != nil {
			return nil, WrapERC20WrapperWithdrawToError(err)
		}
		return nil, WrapERC20WrapperWithdrawToError(fmt.Errorf("nil response from EVM call"))
	}
	if err != nil {
		return res.Ret, WrapERC20WrapperWithdrawToError(err)
	}
	return res.Ret, nil
}

type ERC20WrapperMethodError struct {
	Method string
	Err    error
}

func (e ERC20WrapperMethodError) Error() string {
	if e.Method == "" {
		return fmt.Sprintf("erc20 wrapper call failed: %s", e.Err)
	}
	return fmt.Sprintf("erc20 wrapper %s failed: %s", e.Method, e.Err)
}

func (e ERC20WrapperMethodError) Unwrap() error { return e.Err }

func WrapERC20WrapperMethodError(method string, err error) error {
	if err == nil {
		return nil
	}
	return ERC20WrapperMethodError{Method: method, Err: err}
}

func WrapERC20WrapperUnderlyingError(err error) error {
	return WrapERC20WrapperMethodError("underlying()", err)
}

func WrapERC20WrapperWithdrawToError(err error) error {
	return WrapERC20WrapperMethodError("withdrawTo(address,uint256)", err)
}
