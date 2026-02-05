package distrclaim

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/holiman/uint256"

	cmn "github.com/cosmos/evm/precompiles/common"
	erc20types "github.com/cosmos/evm/x/erc20/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

type delegatorWithdrawAddrGetter interface {
	GetDelegatorWithdrawAddr(ctx context.Context, delAddr sdk.AccAddress) (sdk.AccAddress, error)
}

const (
	maxGasUnderlyingGetter = uint64(200_000)
)

var minWithdrawAmountWad = big.NewInt(1_000_000_000_000) // wmantraUSD SCALAR (1e12)

var (
	selectorUnderlyingGetter         = [4]byte{0x6f, 0x30, 0x7d, 0xc3} // underlying()
	selectorWithdrawToAddressUint256 = [4]byte{0x20, 0x5c, 0x28, 0x78} // withdrawTo(address,uint256)
)

func (p *Precompile) runClaimRewardsAndConvertCoin(
	ctx sdk.Context,
	input []byte,
	evm *vm.EVM,
	contract *vm.Contract,
) ([]byte, error) {
	var in ClaimRewardsAndConvertCoinCall
	if _, err := (&in).Decode(input); err != nil {
		return nil, err
	}

	delegatorAddr := in.Delegator
	maxRetrieve := in.MaxRetrieve
	denom := in.Denom

	if delegatorAddr == (common.Address{}) {
		return nil, fmt.Errorf(cmn.ErrInvalidDelegator, delegatorAddr)
	}
	if denom == "" {
		return nil, fmt.Errorf(cmn.ErrInvalidType, "denom", "string", denom)
	}

	out, err := p.claimRewardsAndConvertCoin(ctx, delegatorAddr, maxRetrieve, denom, evm, contract)
	if err != nil {
		return nil, err
	}

	return out.Encode()
}

func (p *Precompile) claimRewardsAndConvertCoin(
	ctx sdk.Context,
	delegatorAddr common.Address,
	maxRetrieve uint32,
	denom string,
	evm *vm.EVM,
	contract *vm.Contract,
) (*ClaimRewardsAndConvertCoinReturn, error) {
	maxVals, err := p.stakingKeeper.MaxValidators(ctx)
	if err != nil {
		return nil, err
	}
	if maxRetrieve > maxVals {
		return nil, fmt.Errorf(
			"maxRetrieve (%d) parameter exceeds the maximum number of validators (%d)",
			maxRetrieve,
			maxVals,
		)
	}

	msgSender := contract.Caller()
	if msgSender != delegatorAddr {
		return nil, fmt.Errorf(cmn.ErrRequesterIsNotMsgSender, msgSender.String(), delegatorAddr.String())
	}

	delegatorBech32, err := p.addrCdc.BytesToString(delegatorAddr.Bytes())
	if err != nil {
		return nil, err
	}

	withdrawAddrGetter, ok := p.distributionKeeper.(delegatorWithdrawAddrGetter)
	if !ok {
		return nil, fmt.Errorf("distribution keeper does not support GetDelegatorWithdrawAddr")
	}

	prevWithdrawAddr, err := withdrawAddrGetter.GetDelegatorWithdrawAddr(ctx, sdk.AccAddress(delegatorAddr.Bytes()))
	if err != nil {
		return nil, err
	}
	prevWithdrawBech32, err := p.addrCdc.BytesToString(prevWithdrawAddr)
	if err != nil {
		return nil, err
	}
	changedWithdrawAddr := prevWithdrawBech32 != delegatorBech32

	// Ensure rewards are paid to the delegator account itself, so the subsequent
	// conversion burns from the same account.
	if changedWithdrawAddr {
		_, err = p.distributionMsgServer.SetWithdrawAddress(ctx, &distrtypes.MsgSetWithdrawAddress{
			DelegatorAddress: delegatorBech32,
			WithdrawAddress:  delegatorBech32,
		})
		if err != nil {
			return nil, err
		}
	}

	res, err := p.stakingKeeper.GetDelegatorValidators(ctx, delegatorAddr.Bytes(), maxRetrieve)
	if err != nil {
		return nil, err
	}

	totalCoins := sdk.Coins{}
	for _, validator := range res.Validators {
		valAddr, err := sdk.ValAddressFromBech32(validator.OperatorAddress)
		if err != nil {
			return nil, err
		}

		coins, err := p.distributionKeeper.WithdrawDelegationRewards(ctx, delegatorAddr.Bytes(), valAddr)
		if err != nil {
			return nil, err
		}

		totalCoins = totalCoins.Add(coins...)
	}

	converted := big.NewInt(0)
	amount := totalCoins.AmountOf(denom)
	if amount.IsPositive() {
		amtBig := amount.BigInt()
		_, err := p.erc20MsgServer.ConvertCoin(ctx, &erc20types.MsgConvertCoin{
			Coin: sdk.Coin{
				Denom:  denom,
				Amount: amount,
			},
			Receiver: delegatorAddr.Hex(),
			Sender:   delegatorBech32,
		})
		if err != nil {
			return nil, err
		}
		converted = amtBig
		p.tryUnwrapWrapper(ctx, evm, msgSender, contract, denom, amtBig)
	}

	if changedWithdrawAddr {
		_, err := p.distributionMsgServer.SetWithdrawAddress(ctx, &distrtypes.MsgSetWithdrawAddress{
			DelegatorAddress: delegatorBech32,
			WithdrawAddress:  prevWithdrawBech32,
		})
		if err != nil {
			return nil, err
		}
	}

	return &ClaimRewardsAndConvertCoinReturn{Amount: converted}, nil
}

func parseERC20DenomAddress(denom string) (common.Address, error) {
	addrStr := strings.TrimPrefix(denom, "erc20:")
	if addrStr == "" {
		return common.Address{}, fmt.Errorf("invalid erc20 denom: %q", denom)
	}
	if !strings.HasPrefix(addrStr, "0x") {
		addrStr = "0x" + addrStr
	}
	if !common.IsHexAddress(addrStr) {
		return common.Address{}, fmt.Errorf("invalid erc20 denom address: %q", denom)
	}
	addr := common.HexToAddress(addrStr)
	if addr == (common.Address{}) {
		return common.Address{}, fmt.Errorf("invalid erc20 denom address: %q", denom)
	}
	return addr, nil
}

func (p *Precompile) tryUnwrapWrapper(ctx sdk.Context, evm *vm.EVM, caller common.Address, contract *vm.Contract, denom string, amountWad *big.Int) {
	if evm == nil || amountWad == nil {
		return
	}
	if !strings.HasPrefix(denom, "erc20:") {
		return
	}
	if amountWad.Cmp(minWithdrawAmountWad) < 0 {
		return
	}

	wrapperAddr, err := parseERC20DenomAddress(denom)
	if err != nil {
		return
	}

	// Only attempt for wrappers that look like ERC20Wrapper (underlying() + withdrawTo(address,uint256)).
	if _, err := evmGetUnderlyingERC20Wrapper(evm, caller, wrapperAddr, contract); err != nil {
		return
	}

	// Best-effort unwrap via x/vm keeper (sees ConvertCoin state updates).
	_, _ = p.evmWithdrawToAddressUint256(ctx, caller, wrapperAddr, caller, amountWad)
}

func evmGetUnderlyingERC20Wrapper(evm *vm.EVM, caller common.Address, wrapper common.Address, contract *vm.Contract) (common.Address, error) {
	gas := contract.Gas
	if gas > maxGasUnderlyingGetter {
		gas = maxGasUnderlyingGetter
	}
	data := make([]byte, 4)
	copy(data[:4], selectorUnderlyingGetter[:])
	ret, left, err := evm.StaticCall(caller, wrapper, data, gas)
	if err != nil {
		return common.Address{}, fmt.Errorf("wrapper underlying() call failed: %w", err)
	}
	used := gas - left
	if used > 0 {
		if !contract.UseGas(used, nil, tracing.GasChangeCallPrecompiledContract) {
			return common.Address{}, vm.ErrOutOfGas
		}
	}
	if len(ret) < 32 {
		return common.Address{}, fmt.Errorf("wrapper underlying() returned short data")
	}
	underlying := common.BytesToAddress(ret[12:32])
	if underlying == (common.Address{}) {
		return common.Address{}, fmt.Errorf("wrapper underlying() returned zero address")
	}
	return underlying, nil
}

func (p *Precompile) evmWithdrawToAddressUint256(ctx sdk.Context, caller common.Address, wrapper common.Address, to common.Address, amount *big.Int) ([]byte, error) {
	if amount == nil || amount.Sign() <= 0 {
		return nil, fmt.Errorf("invalid withdraw amount")
	}
	u, overflow := uint256.FromBig(amount)
	if overflow {
		return nil, fmt.Errorf("withdraw amount overflows uint256")
	}
	data := make([]byte, 4+32+32)
	copy(data[:4], selectorWithdrawToAddressUint256[:])
	copy(data[4:4+32], common.LeftPadBytes(to.Bytes(), 32))
	amt32 := u.Bytes32()
	copy(data[4+32:], amt32[:])

	res, err := p.evmKeeper.CallEVMWithData(ctx, caller, &wrapper, data, true, nil)
	if res == nil {
		return nil, err
	}
	return res.Ret, err
}
