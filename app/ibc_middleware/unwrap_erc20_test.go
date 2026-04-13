package ibc_middleware

import (
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/MANTRA-Chain/mantrachain/v8/app/evmutil"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/evm/x/vm/statedb"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	vmmocks "github.com/cosmos/evm/x/vm/types/mocks"
)

type mockEVMKeeper struct {
	*vmmocks.EVMKeeper

	underlyingErr error
	withdrawErr   error
}

var (
	testReceiver = common.HexToAddress("0x0000000000000000000000000000000000000001")
	testWrapper  = common.HexToAddress("0x0000000000000000000000000000000000000002")
)

const testDenom = "ibc/ABC123"

func minWithdrawAmount() *big.Int {
	return new(big.Int).Set(evmutil.MinWithdrawAmountWad)
}

func newMockEVMKeeper() mockEVMKeeper {
	return mockEVMKeeper{EVMKeeper: vmmocks.NewEVMKeeper()}
}

func (m mockEVMKeeper) CallEVM(
	ctx sdk.Context,
	stateDB *statedb.StateDB,
	contractABI abi.ABI,
	from common.Address,
	contract common.Address,
	commit bool,
	callFromPrecompile bool,
	gasCap *big.Int,
	method string,
	args ...interface{},
) (*evmtypes.MsgEthereumTxResponse, error) {
	return nil, errors.New("unexpected CallEVM")
}

func (m mockEVMKeeper) GetAccountOrEmpty(ctx sdk.Context, addr common.Address) statedb.Account {
	acc := m.GetAccount(ctx, addr)
	if acc == nil {
		return statedb.Account{}
	}
	return *acc
}

func (m mockEVMKeeper) IsContract(ctx sdk.Context, addr common.Address) bool {
	return m.GetAccount(ctx, addr) != nil
}

func (m mockEVMKeeper) CallEVMWithData(
	ctx sdk.Context,
	stateDB *statedb.StateDB,
	from common.Address,
	contract *common.Address,
	data []byte,
	commit bool,
	callFromPrecompile bool,
	gasCap *big.Int,
) (*evmtypes.MsgEthereumTxResponse, error) {
	if len(data) < 4 {
		return nil, errors.New("short call data")
	}

	switch [4]byte(data[:4]) {
	case evmutil.SelectorUnderlyingGetter:
		if m.underlyingErr != nil {
			return nil, m.underlyingErr
		}

		ret := make([]byte, 32)
		copy(ret[12:], common.HexToAddress("0x000000000000000000000000000000000000dead").Bytes())
		return &evmtypes.MsgEthereumTxResponse{Ret: ret}, nil
	case evmutil.SelectorWithdrawToAddressUint256:
		if m.withdrawErr != nil {
			return nil, m.withdrawErr
		}
		return &evmtypes.MsgEthereumTxResponse{Ret: []byte{}}, nil
	default:
		return nil, errors.New("unexpected selector")
	}
}

func TestTryUnwrapEmitsOutcomeEvents(t *testing.T) {
	tests := []struct {
		name            string
		amount          *big.Int
		keeper          mockEVMKeeper
		expectSuccess   string
		expectFailStep  string
		expectErrorAttr bool
	}{
		{
			name:           "invalid amount",
			amount:         nil,
			keeper:         newMockEVMKeeper(),
			expectSuccess:  "false",
			expectFailStep: "invalid_amount",
		},
		{
			name:           "amount below minimum",
			amount:         new(big.Int).Sub(evmutil.MinWithdrawAmountWad, big.NewInt(1)),
			keeper:         newMockEVMKeeper(),
			expectSuccess:  "false",
			expectFailStep: "amount_below_min",
		},
		{
			name:            "underlying probe fails",
			amount:          minWithdrawAmount(),
			keeper:          mockEVMKeeper{EVMKeeper: vmmocks.NewEVMKeeper(), underlyingErr: errors.New("probe failed")},
			expectSuccess:   "false",
			expectFailStep:  "underlying",
			expectErrorAttr: true,
		},
		{
			name:            "withdraw fails",
			amount:          minWithdrawAmount(),
			keeper:          mockEVMKeeper{EVMKeeper: vmmocks.NewEVMKeeper(), withdrawErr: errors.New("withdraw failed")},
			expectSuccess:   "false",
			expectFailStep:  "withdraw_to",
			expectErrorAttr: true,
		},
		{
			name:           "successful unwrap",
			amount:         minWithdrawAmount(),
			keeper:         newMockEVMKeeper(),
			expectSuccess:  "true",
			expectFailStep: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := sdk.Context{}.WithEventManager(sdk.NewEventManager())
			im := UnwrapERC20IBCModule{evmCaller: tc.keeper}

			im.tryUnwrap(
				ctx,
				testDenom,
				testReceiver,
				testWrapper,
				tc.amount,
			)

			events := ctx.EventManager().Events()
			if len(events) != 1 {
				t.Fatalf("expected exactly one event, got %d", len(events))
			}

			evt := events[0]
			if evt.Type != EventTypeUnwrapERC20 {
				t.Fatalf("expected event type %q, got %q", EventTypeUnwrapERC20, evt.Type)
			}

			attrs := attributesToMap(evt)
			if attrs[AttributeKeyUnwrapSuccess] != tc.expectSuccess {
				t.Fatalf("expected success=%q, got %q", tc.expectSuccess, attrs[AttributeKeyUnwrapSuccess])
			}

			if tc.expectFailStep == "" {
				if _, ok := attrs[AttributeKeyUnwrapFailureStep]; ok {
					t.Fatalf("did not expect %q attribute", AttributeKeyUnwrapFailureStep)
				}
			} else if attrs[AttributeKeyUnwrapFailureStep] != tc.expectFailStep {
				t.Fatalf("expected failure_step=%q, got %q", tc.expectFailStep, attrs[AttributeKeyUnwrapFailureStep])
			}

			_, hasErr := attrs[AttributeKeyUnwrapError]
			if hasErr != tc.expectErrorAttr {
				t.Fatalf("expected error attr presence=%t, got %t", tc.expectErrorAttr, hasErr)
			}
		})
	}
}

func attributesToMap(event sdk.Event) map[string]string {
	attrs := make(map[string]string, len(event.Attributes))
	for _, attr := range event.Attributes {
		attrs[attr.Key] = attr.Value
	}
	return attrs
}
