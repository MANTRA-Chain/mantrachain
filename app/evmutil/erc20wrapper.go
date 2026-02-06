package evmutil

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"

	sdk "github.com/cosmos/cosmos-sdk/types"

	evmtypes "github.com/cosmos/evm/x/vm/types"
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

type EVMCaller interface {
	CallEVMWithData(
		ctx sdk.Context,
		from common.Address,
		contract *common.Address,
		data []byte,
		commit bool,
		gasCap *big.Int,
	) (*evmtypes.MsgEthereumTxResponse, error)
}

func ERC20WrapperUnderlyingViaEVMCaller(ctx sdk.Context, caller EVMCaller, from common.Address, wrapper common.Address) (common.Address, error) {
	data := ERC20WrapperUnderlyingCallData()
	res, err := caller.CallEVMWithData(ctx, from, &wrapper, data, false, GasCapERC20WrapperUnderlying)
	if res == nil {
		return common.Address{}, err
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

func ERC20WrapperWithdrawToViaEVMCaller(ctx sdk.Context, caller EVMCaller, from common.Address, wrapper common.Address, to common.Address, amountWad *big.Int) ([]byte, error) {
	data, err := ERC20WrapperWithdrawToCallData(to, amountWad)
	if err != nil {
		return nil, err
	}
	res, err := caller.CallEVMWithData(ctx, from, &wrapper, data, true, GasCapERC20WrapperWithdrawTo)
	if res == nil {
		return nil, err
	}
	return res.Ret, err
}

type ERC20WrapperMethodError struct {
	Method string
	Err    error
}

func (e ERC20WrapperMethodError) Error() string {
	if e.Method == "" {
		return fmt.Sprintf("erc20 wrapper call failed: %v", e.Err)
	}
	return fmt.Sprintf("erc20 wrapper %s failed: %v", e.Method, e.Err)
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
