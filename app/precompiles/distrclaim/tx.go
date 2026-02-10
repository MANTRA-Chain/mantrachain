package distrclaim

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/MANTRA-Chain/mantrachain/v8/app/evmutil"

	cmn "github.com/cosmos/evm/precompiles/common"
	erc20types "github.com/cosmos/evm/x/erc20/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
	if emitErr := p.emitClaimRewardsAndConvertCoinEvent(ctx, evm, delegatorAddr, denom, out.Amount); emitErr != nil {
		return nil, emitErr
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

	withdrawAddrSdk, err := p.distributionKeeper.GetDelegatorWithdrawAddr(ctx, sdk.AccAddress(delegatorAddr.Bytes()))
	if err != nil {
		return nil, err
	}
	withdrawAddrBech32, err := p.addrCdc.BytesToString(withdrawAddrSdk)
	if err != nil {
		return nil, err
	}
	if withdrawAddrBech32 != delegatorBech32 {
		return nil, fmt.Errorf(
			"withdraw address %s must equal delegator %s",
			withdrawAddrBech32,
			delegatorBech32,
		)
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
	if amountWad.Cmp(evmutil.MinWithdrawAmountWad) < 0 {
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
	_, _ = evmutil.ERC20WrapperWithdrawToViaEVMCaller(ctx, p.evmKeeper, caller, wrapperAddr, caller, amountWad)
}

func evmGetUnderlyingERC20Wrapper(evm *vm.EVM, caller common.Address, wrapper common.Address, contract *vm.Contract) (common.Address, error) {
	gas := contract.Gas
	if cap := evmutil.GasCapERC20WrapperUnderlying.Uint64(); gas > cap {
		gas = cap
	}
	data := evmutil.ERC20WrapperUnderlyingCallData()
	ret, left, err := evm.StaticCall(caller, wrapper, data, gas)
	if err != nil {
		return common.Address{}, evmutil.WrapERC20WrapperUnderlyingError(err)
	}
	used := gas - left
	if used > 0 {
		if !contract.UseGas(used, nil, tracing.GasChangeCallPrecompiledContract) {
			return common.Address{}, vm.ErrOutOfGas
		}
	}
	underlying, err := evmutil.DecodeABIAddress32(ret)
	if err != nil {
		return common.Address{}, evmutil.WrapERC20WrapperUnderlyingError(err)
	}
	return underlying, nil
}
