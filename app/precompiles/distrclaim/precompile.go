package distrclaim

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"cosmossdk.io/core/address"
	storetypes "cosmossdk.io/store/types"

	cmn "github.com/cosmos/evm/precompiles/common"
	erc20types "github.com/cosmos/evm/x/erc20/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

const (
	// DistributionClaimPrecompileAddress is a static precompile for claim+convert.
	DistributionClaimPrecompileAddress = "0x0000000000000000000000000000000000000a01"
)

//go:generate go run github.com/yihuang/go-abi/cmd -var=HumanABI -output distrclaim.abi.go
var HumanABI = []string{
	"function claimRewardsAndConvertCoin(address delegator, uint32 maxRetrieve, string denom) returns (uint256 amount)",
}

var _ vm.PrecompiledContract = &Precompile{}

// Precompile implements `claimRewardsAndConvertCoin`.
// It sets the caller's withdraw address, claims up to `maxRetrieve` rewards, then
// converts the specified `denom` into ERC20 via x/erc20 and returns the amount.
type Precompile struct {
	cmn.Precompile
	stakingKeeper         cmn.StakingKeeper
	distributionKeeper    cmn.DistributionKeeper
	distributionMsgServer distrtypes.MsgServer
	erc20MsgServer        erc20types.MsgServer
	addrCdc               address.Codec
}

func NewPrecompile(
	bankKeeper cmn.BankKeeper,
	stakingKeeper cmn.StakingKeeper,
	distributionKeeper cmn.DistributionKeeper,
	distributionMsgServer distrtypes.MsgServer,
	erc20MsgServer erc20types.MsgServer,
	addrCdc address.Codec,
) *Precompile {
	return &Precompile{
		Precompile: cmn.Precompile{
			KvGasConfig:           storetypes.KVGasConfig(),
			TransientKVGasConfig:  storetypes.TransientGasConfig(),
			ContractAddress:       common.HexToAddress(DistributionClaimPrecompileAddress),
			BalanceHandlerFactory: cmn.NewBalanceHandlerFactory(bankKeeper),
		},
		stakingKeeper:         stakingKeeper,
		distributionKeeper:    distributionKeeper,
		distributionMsgServer: distributionMsgServer,
		erc20MsgServer:        erc20MsgServer,
		addrCdc:               addrCdc,
	}
}

func (p Precompile) RequiredGas(input []byte) uint64 {
	if len(input) < 4 {
		return 0
	}

	methodID := binary.BigEndian.Uint32(input)
	return p.Precompile.RequiredGas(input, p.IsTransactionID(methodID))
}

func (p Precompile) Run(evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	return p.RunNativeAction(evm, contract, func(ctx sdk.Context) ([]byte, error) {
		return p.Execute(ctx, evm.StateDB, contract, readonly)
	})
}

func (p Precompile) Execute(ctx sdk.Context, stateDB vm.StateDB, contract *vm.Contract, readOnly bool) ([]byte, error) {
	methodID, input, err := splitMethodID(contract.Input)
	if err != nil {
		return nil, err
	}

	if readOnly && p.IsTransactionID(methodID) {
		return nil, vm.ErrWriteProtection
	}

	switch methodID {
	case ClaimRewardsAndConvertCoinID:
		return p.runClaimRewardsAndConvertCoin(ctx, input, stateDB, contract)
	default:
		return nil, fmt.Errorf("unknown method id: %d", methodID)
	}
}

func (Precompile) IsTransactionID(methodID uint32) bool {
	return methodID == ClaimRewardsAndConvertCoinID
}

func splitMethodID(input []byte) (uint32, []byte, error) {
	if len(input) < 4 {
		return 0, nil, errors.New("invalid input length")
	}

	methodID := binary.BigEndian.Uint32(input)
	return methodID, input[4:], nil
}
