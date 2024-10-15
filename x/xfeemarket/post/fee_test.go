package post_test

import (
	"fmt"
	"testing"

	"cosmossdk.io/math"
	"github.com/MANTRA-Chain/mantrachain/x/xfeemarket/post"
	postsuite "github.com/MANTRA-Chain/mantrachain/x/xfeemarket/post/suite"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/skip-mev/feemarket/x/feemarket/types"
	"github.com/stretchr/testify/mock"
)

func TestBurnCoins(t *testing.T) {
	tests := []struct {
		name           string
		coins          sdk.Coins
		distributeFees bool
		wantErr        bool
	}{
		{
			name:    "valid",
			coins:   sdk.NewCoins(sdk.NewCoin("test", math.NewInt(10))),
			wantErr: false,
		},
		{
			name:    "valid no coins",
			coins:   sdk.NewCoins(),
			wantErr: false,
		},
		{
			name:    "valid zero coin",
			coins:   sdk.NewCoins(sdk.NewCoin("test", math.ZeroInt())),
			wantErr: false,
		},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("Case %s", tc.name), func(t *testing.T) {
			s := postsuite.SetupTestSuite(t, true)
			// Set up mock expectations based on what BurnCoins is expected to call
			s.MockBankKeeper.On("BurnCoins", s.Ctx, types.FeeCollectorName, tc.coins).Return(nil).Once()
			if err := post.BurnCoins(s.MockBankKeeper, s.Ctx, tc.coins); (err != nil) != tc.wantErr {
				s.Errorf(err, "BurnCoins() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestBurnCoinsAndRefund(t *testing.T) {
	tests := []struct {
		name    string
		coins   sdk.Coins
		wantErr bool
	}{
		{
			name:    "valid",
			coins:   sdk.NewCoins(sdk.NewCoin("test", math.NewInt(10))),
			wantErr: false,
		},
		{
			name:    "valid no coins",
			coins:   sdk.NewCoins(),
			wantErr: false,
		},
		{
			name:    "valid zero coin",
			coins:   sdk.NewCoins(sdk.NewCoin("test", math.ZeroInt())),
			wantErr: false,
		},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("Case %s", tc.name), func(t *testing.T) {
			s := postsuite.SetupTestSuite(t, true)
			accs := s.CreateTestAccounts(1)

			s.MockBankKeeper.On("BurnCoins", s.Ctx, types.FeeCollectorName, tc.coins).Return(nil).Once()
			if err := post.BurnCoins(s.MockBankKeeper, s.Ctx, tc.coins); (err != nil) != tc.wantErr {
				s.Errorf(err, "DeductCoins() error = %v, wantErr %v", err, tc.wantErr)
			}

			s.MockBankKeeper.On("SendCoinsFromModuleToAccount", s.Ctx, types.FeeCollectorName, accs[0].Account.GetAddress(),
				tc.coins).Return(nil).Once()
			if err := post.RefundTip(s.MockBankKeeper, s.Ctx, accs[0].Account.GetAddress(), tc.coins); (err != nil) != tc.wantErr {
				s.Errorf(err, "DeductCoins() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestRefund(t *testing.T) {
	tests := []struct {
		name    string
		coins   sdk.Coins
		wantErr bool
	}{
		{
			name:    "valid",
			coins:   sdk.NewCoins(sdk.NewCoin("test", math.NewInt(10))),
			wantErr: false,
		},
		{
			name:    "valid no coins",
			coins:   sdk.NewCoins(),
			wantErr: false,
		},
		{
			name:    "valid zero coin",
			coins:   sdk.NewCoins(sdk.NewCoin("test", math.ZeroInt())),
			wantErr: false,
		},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("Case %s", tc.name), func(t *testing.T) {
			s := postsuite.SetupTestSuite(t, true)
			accs := s.CreateTestAccounts(1)
			s.MockBankKeeper.On("SendCoinsFromModuleToAccount", s.Ctx, types.FeeCollectorName, accs[0].Account.GetAddress(),
				tc.coins).Return(nil).Once()

			if err := post.RefundTip(s.MockBankKeeper, s.Ctx, accs[0].Account.GetAddress(), tc.coins); (err != nil) != tc.wantErr {
				s.Errorf(err, "SendTip() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

//nolint:maintidx
func TestPostHandleMock(t *testing.T) {
	// Same data for every test case
	const (
		baseDenom              = "stake"
		resolvableDenom        = "atom"
		expectedConsumedGas    = 9307
		expectedConsumedSimGas = expectedConsumedGas + post.BankSendGasConsumption
		gasLimit               = expectedConsumedSimGas // 152800 12490 + 36385
	)

	validFeeAmount := types.DefaultMinBaseGasPrice.MulInt64(int64(gasLimit))
	validFee := sdk.NewCoins(sdk.NewCoin(baseDenom, validFeeAmount.TruncateInt()))
	validResolvableFee := sdk.NewCoins(sdk.NewCoin(resolvableDenom, validFeeAmount.TruncateInt()))

	testCases := []postsuite.TestCase{
		{
			Name: "signer has no funds",
			Malleate: func(s *postsuite.TestSuite) postsuite.TestCaseArgs {
				accs := s.CreateTestAccounts(1)
				s.MockBankKeeper.On("SendCoinsFromAccountToModule", mock.Anything, accs[0].Account.GetAddress(),
					types.FeeCollectorName, mock.Anything).Return(sdkerrors.ErrInsufficientFunds).Once()

				return postsuite.TestCaseArgs{
					Msgs:      []sdk.Msg{testdata.NewTestMsg(accs[0].Account.GetAddress())},
					GasLimit:  gasLimit,
					FeeAmount: validFee,
				}
			},
			RunAnte:  true,
			RunPost:  true,
			Simulate: false,
			ExpPass:  false,
			ExpErr:   sdkerrors.ErrInsufficientFunds,
			Mock:     true,
		},
		{
			Name: "signer has no funds - simulate",
			Malleate: func(s *postsuite.TestSuite) postsuite.TestCaseArgs {
				accs := s.CreateTestAccounts(1)
				s.MockBankKeeper.On("SendCoinsFromAccountToModule", mock.Anything, accs[0].Account.GetAddress(),
					types.FeeCollectorName, mock.Anything).Return(sdkerrors.ErrInsufficientFunds).Once()

				return postsuite.TestCaseArgs{
					Msgs:      []sdk.Msg{testdata.NewTestMsg(accs[0].Account.GetAddress())},
					GasLimit:  gasLimit,
					FeeAmount: validFee,
				}
			},
			RunAnte:  true,
			RunPost:  true,
			Simulate: true,
			ExpPass:  false,
			ExpErr:   sdkerrors.ErrInsufficientFunds,
			Mock:     true,
		},
		{
			Name: "0 gas given should fail",
			Malleate: func(s *postsuite.TestSuite) postsuite.TestCaseArgs {
				accs := s.CreateTestAccounts(1)

				return postsuite.TestCaseArgs{
					Msgs:      []sdk.Msg{testdata.NewTestMsg(accs[0].Account.GetAddress())},
					GasLimit:  0,
					FeeAmount: validFee,
				}
			},
			RunAnte:  true,
			RunPost:  true,
			Simulate: false,
			ExpPass:  false,
			ExpErr:   sdkerrors.ErrOutOfGas,
			Mock:     true,
		},
		{
			Name: "0 gas given should pass - simulate",
			Malleate: func(s *postsuite.TestSuite) postsuite.TestCaseArgs {
				accs := s.CreateTestAccounts(1)
				s.MockBankKeeper.On("SendCoinsFromAccountToModule", mock.Anything, accs[0].Account.GetAddress(),
					types.FeeCollectorName, mock.Anything).Return(nil).Once()

				return postsuite.TestCaseArgs{
					Msgs:      []sdk.Msg{testdata.NewTestMsg(accs[0].Account.GetAddress())},
					GasLimit:  0,
					FeeAmount: validFee,
				}
			},
			RunAnte:           true,
			RunPost:           true,
			Simulate:          true,
			ExpPass:           true,
			ExpErr:            nil,
			ExpectConsumedGas: expectedConsumedSimGas,
			Mock:              true,
		},
		{
			Name: "signer has enough funds, should pass, no tip",
			Malleate: func(s *postsuite.TestSuite) postsuite.TestCaseArgs {
				accs := s.CreateTestAccounts(1)
				s.MockBankKeeper.On("SendCoinsFromAccountToModule", mock.Anything, accs[0].Account.GetAddress(),
					types.FeeCollectorName, mock.Anything).Return(nil)
				s.MockBankKeeper.On("BurnCoins", mock.Anything, types.FeeCollectorName, mock.Anything).Return(nil).Once()
				s.MockBankKeeper.On("SendCoinsFromModuleToAccount", mock.Anything, types.FeeCollectorName,
					accs[0].Account.GetAddress(), mock.Anything).Return(nil).Once()
				return postsuite.TestCaseArgs{
					Msgs:      []sdk.Msg{testdata.NewTestMsg(accs[0].Account.GetAddress())},
					GasLimit:  gasLimit,
					FeeAmount: validFee,
				}
			},
			RunAnte:           true,
			RunPost:           true,
			Simulate:          false,
			ExpPass:           true,
			ExpErr:            nil,
			ExpectConsumedGas: expectedConsumedGas,
			Mock:              true,
		},
		{
			Name: "fee market is enabled during the transaction - should pass and skip deduction until next block",
			Malleate: func(s *postsuite.TestSuite) postsuite.TestCaseArgs {
				accs := s.CreateTestAccounts(1)

				// disable fee market before tx
				s.Ctx = s.Ctx.WithBlockHeight(10)
				disabledParams := types.DefaultParams()
				disabledParams.Enabled = false
				err := s.FeeMarketKeeper.SetParams(s.Ctx, disabledParams)
				s.Require().NoError(err)

				s.MockBankKeeper.On("SendCoinsFromAccountToModule", mock.Anything, accs[0].Account.GetAddress(),
					authtypes.FeeCollectorName, mock.Anything).Return(nil).Once()

				return postsuite.TestCaseArgs{
					Msgs:      []sdk.Msg{testdata.NewTestMsg(accs[0].Account.GetAddress())},
					GasLimit:  gasLimit,
					FeeAmount: validResolvableFee,
				}
			},
			StateUpdate: func(s *postsuite.TestSuite) {
				// enable the fee market
				enabledParams := types.DefaultParams()
				req := &types.MsgParams{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Params:    enabledParams,
				}

				_, err := s.MsgServer.Params(s.Ctx, req)
				s.Require().NoError(err)

				height, err := s.FeeMarketKeeper.GetEnabledHeight(s.Ctx)
				s.Require().NoError(err)
				s.Require().Equal(int64(10), height)
			},
			RunAnte:           true,
			RunPost:           true,
			Simulate:          false,
			ExpPass:           true,
			ExpErr:            nil,
			ExpectConsumedGas: 15340, // extra gas consumed because msg server is run, but deduction is skipped
			Mock:              true,
		},
		{
			Name: "signer has enough funds, should pass, no tip - resolvable denom",
			Malleate: func(s *postsuite.TestSuite) postsuite.TestCaseArgs {
				accs := s.CreateTestAccounts(1)
				s.MockBankKeeper.On("SendCoinsFromAccountToModule", mock.Anything, accs[0].Account.GetAddress(),
					types.FeeCollectorName, mock.Anything).Return(nil).Once()
				s.MockBankKeeper.On("SendCoinsFromModuleToAccount", mock.Anything, types.FeeCollectorName,
					accs[0].Account.GetAddress(), mock.Anything).Return(nil).Once()

				return postsuite.TestCaseArgs{
					Msgs:      []sdk.Msg{testdata.NewTestMsg(accs[0].Account.GetAddress())},
					GasLimit:  gasLimit,
					FeeAmount: validResolvableFee,
				}
			},
			RunAnte:           true,
			RunPost:           true,
			Simulate:          false,
			ExpPass:           true,
			ExpErr:            nil,
			ExpectConsumedGas: expectedConsumedGas,
			Mock:              true,
		},
		{
			Name: "signer has enough funds, should pass, no tip - resolvable denom - simulate",
			Malleate: func(s *postsuite.TestSuite) postsuite.TestCaseArgs {
				accs := s.CreateTestAccounts(1)
				s.MockBankKeeper.On("SendCoinsFromAccountToModule", mock.Anything, accs[0].Account.GetAddress(),
					types.FeeCollectorName, mock.Anything).Return(nil).Once()

				return postsuite.TestCaseArgs{
					Msgs:      []sdk.Msg{testdata.NewTestMsg(accs[0].Account.GetAddress())},
					GasLimit:  gasLimit,
					FeeAmount: validResolvableFee,
				}
			},
			RunAnte:           true,
			RunPost:           true,
			Simulate:          true,
			ExpPass:           true,
			ExpErr:            nil,
			ExpectConsumedGas: expectedConsumedSimGas,
			Mock:              true,
		},
		{
			Name: "0 gas given should pass in simulate - no fee",
			Malleate: func(s *postsuite.TestSuite) postsuite.TestCaseArgs {
				accs := s.CreateTestAccounts(1)
				s.MockBankKeeper.On("SendCoinsFromAccountToModule", mock.Anything, accs[0].Account.GetAddress(),
					types.FeeCollectorName, mock.Anything).Return(nil).Once()
				return postsuite.TestCaseArgs{
					Msgs:      []sdk.Msg{testdata.NewTestMsg(accs[0].Account.GetAddress())},
					GasLimit:  0,
					FeeAmount: nil,
				}
			},
			RunAnte:           true,
			RunPost:           false,
			Simulate:          true,
			ExpPass:           true,
			ExpErr:            nil,
			ExpectConsumedGas: expectedConsumedSimGas,
			Mock:              true,
		},
		{
			Name: "0 gas given should pass in simulate - fee",
			Malleate: func(s *postsuite.TestSuite) postsuite.TestCaseArgs {
				accs := s.CreateTestAccounts(1)
				s.MockBankKeeper.On("SendCoinsFromAccountToModule", mock.Anything, accs[0].Account.GetAddress(),
					types.FeeCollectorName, mock.Anything).Return(nil).Once()
				return postsuite.TestCaseArgs{
					Msgs:      []sdk.Msg{testdata.NewTestMsg(accs[0].Account.GetAddress())},
					GasLimit:  0,
					FeeAmount: validFee,
				}
			},
			RunAnte:           true,
			RunPost:           false,
			Simulate:          true,
			ExpPass:           true,
			ExpErr:            nil,
			ExpectConsumedGas: expectedConsumedSimGas,
			Mock:              true,
		},
		{
			Name: "no fee - fail",
			Malleate: func(s *postsuite.TestSuite) postsuite.TestCaseArgs {
				accs := s.CreateTestAccounts(1)

				return postsuite.TestCaseArgs{
					Msgs:      []sdk.Msg{testdata.NewTestMsg(accs[0].Account.GetAddress())},
					GasLimit:  1000000000,
					FeeAmount: nil,
				}
			},
			RunAnte:  true,
			RunPost:  true,
			Simulate: false,
			ExpPass:  false,
			ExpErr:   types.ErrNoFeeCoins,
			Mock:     true,
		},
		{
			Name: "no gas limit - fail",
			Malleate: func(s *postsuite.TestSuite) postsuite.TestCaseArgs {
				accs := s.CreateTestAccounts(1)

				return postsuite.TestCaseArgs{
					Msgs:      []sdk.Msg{testdata.NewTestMsg(accs[0].Account.GetAddress())},
					GasLimit:  0,
					FeeAmount: nil,
				}
			},
			RunAnte:  true,
			RunPost:  true,
			Simulate: false,
			ExpPass:  false,
			ExpErr:   sdkerrors.ErrOutOfGas,
			Mock:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Case %s", tc.Name), func(t *testing.T) {
			s := postsuite.SetupTestSuite(t, tc.Mock)
			s.TxBuilder = s.ClientCtx.TxConfig.NewTxBuilder()
			args := tc.Malleate(s)
			s.RunTestCase(t, tc, args)
		})
	}
}

//nolint:maintidx
func TestPostHandle(t *testing.T) {
	// Same data for every test case
	const (
		baseDenom           = "stake"
		resolvableDenom     = "atom"
		expectedConsumedGas = 59122

		expectedConsumedGasResolve = 25360 // slight difference due to denom resolver

		gasLimit = 100000
	)

	validFeeAmount := types.DefaultMinBaseGasPrice.MulInt64(int64(gasLimit))
	validFee := sdk.NewCoins(sdk.NewCoin(baseDenom, validFeeAmount.TruncateInt()))
	validResolvableFee := sdk.NewCoins(sdk.NewCoin(resolvableDenom, validFeeAmount.TruncateInt()))

	testCases := []postsuite.TestCase{
		{
			Name: "signer has no funds",
			Malleate: func(s *postsuite.TestSuite) postsuite.TestCaseArgs {
				accs := s.CreateTestAccounts(1)

				return postsuite.TestCaseArgs{
					Msgs:      []sdk.Msg{testdata.NewTestMsg(accs[0].Account.GetAddress())},
					GasLimit:  gasLimit,
					FeeAmount: validFee,
				}
			},
			RunAnte:  true,
			RunPost:  true,
			Simulate: false,
			ExpPass:  false,
			ExpErr:   sdkerrors.ErrInsufficientFunds,
			Mock:     false,
		},
		{
			Name: "signer has no funds - simulate - pass",
			Malleate: func(s *postsuite.TestSuite) postsuite.TestCaseArgs {
				accs := s.CreateTestAccounts(1)

				return postsuite.TestCaseArgs{
					Msgs:      []sdk.Msg{testdata.NewTestMsg(accs[0].Account.GetAddress())},
					GasLimit:  gasLimit,
					FeeAmount: validFee,
				}
			},
			RunAnte:           true,
			RunPost:           true,
			Simulate:          true,
			ExpPass:           true,
			ExpErr:            nil,
			Mock:              false,
			ExpectConsumedGas: expectedConsumedGas,
		},
		{
			Name: "0 gas given should fail",
			Malleate: func(s *postsuite.TestSuite) postsuite.TestCaseArgs {
				accs := s.CreateTestAccounts(1)

				return postsuite.TestCaseArgs{
					Msgs:      []sdk.Msg{testdata.NewTestMsg(accs[0].Account.GetAddress())},
					GasLimit:  0,
					FeeAmount: validFee,
				}
			},
			RunAnte:  true,
			RunPost:  true,
			Simulate: false,
			ExpPass:  false,
			ExpErr:   sdkerrors.ErrOutOfGas,
			Mock:     false,
		},
		{
			Name: "0 gas given should pass - simulate",
			Malleate: func(s *postsuite.TestSuite) postsuite.TestCaseArgs {
				accs := s.CreateTestAccounts(1)

				return postsuite.TestCaseArgs{
					Msgs:      []sdk.Msg{testdata.NewTestMsg(accs[0].Account.GetAddress())},
					GasLimit:  0,
					FeeAmount: validFee,
				}
			},
			RunAnte:           true,
			RunPost:           true,
			Simulate:          true,
			ExpPass:           true,
			ExpErr:            nil,
			ExpectConsumedGas: expectedConsumedGas,
			Mock:              false,
		},
		{
			Name: "signer has enough funds, should pass, no tip",
			Malleate: func(s *postsuite.TestSuite) postsuite.TestCaseArgs {
				accs := s.CreateTestAccounts(1)

				balance := postsuite.TestAccountBalance{
					TestAccount: accs[0],
					Coins:       validFee,
				}
				s.SetAccountBalances([]postsuite.TestAccountBalance{balance})

				return postsuite.TestCaseArgs{
					Msgs:      []sdk.Msg{testdata.NewTestMsg(accs[0].Account.GetAddress())},
					GasLimit:  gasLimit,
					FeeAmount: validFee,
				}
			},
			RunAnte:           true,
			RunPost:           true,
			Simulate:          false,
			ExpPass:           true,
			ExpErr:            nil,
			ExpectConsumedGas: 34836,
			Mock:              false,
		},
		{
			Name: "fee market is enabled during the transaction - should pass and skip deduction until next block",
			Malleate: func(s *postsuite.TestSuite) postsuite.TestCaseArgs {
				accs := s.CreateTestAccounts(1)

				balance := postsuite.TestAccountBalance{
					TestAccount: accs[0],
					Coins:       validResolvableFee,
				}
				s.SetAccountBalances([]postsuite.TestAccountBalance{balance})

				// disable fee market before tx
				s.Ctx = s.Ctx.WithBlockHeight(10)
				disabledParams := types.DefaultParams()
				disabledParams.Enabled = false
				err := s.FeeMarketKeeper.SetParams(s.Ctx, disabledParams)
				s.Require().NoError(err)

				return postsuite.TestCaseArgs{
					Msgs:      []sdk.Msg{testdata.NewTestMsg(accs[0].Account.GetAddress())},
					GasLimit:  gasLimit,
					FeeAmount: validResolvableFee,
				}
			},
			StateUpdate: func(s *postsuite.TestSuite) {
				// enable the fee market
				enabledParams := types.DefaultParams()
				req := &types.MsgParams{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Params:    enabledParams,
				}

				_, err := s.MsgServer.Params(s.Ctx, req)
				s.Require().NoError(err)

				height, err := s.FeeMarketKeeper.GetEnabledHeight(s.Ctx)
				s.Require().NoError(err)
				s.Require().Equal(int64(10), height)
			},
			RunAnte:           true,
			RunPost:           true,
			Simulate:          false,
			ExpPass:           true,
			ExpErr:            nil,
			ExpectConsumedGas: 15340, // extra gas consumed because msg server is run, but bank keepers are skipped
			Mock:              false,
		},
		{
			Name: "signer has enough funds, should pass, no tip - resolvable denom",
			Malleate: func(s *postsuite.TestSuite) postsuite.TestCaseArgs {
				accs := s.CreateTestAccounts(1)

				balance := postsuite.TestAccountBalance{
					TestAccount: accs[0],
					Coins:       validResolvableFee,
				}
				s.SetAccountBalances([]postsuite.TestAccountBalance{balance})

				return postsuite.TestCaseArgs{
					Msgs:      []sdk.Msg{testdata.NewTestMsg(accs[0].Account.GetAddress())},
					GasLimit:  gasLimit,
					FeeAmount: validResolvableFee,
				}
			},
			RunAnte:           true,
			RunPost:           true,
			Simulate:          false,
			ExpPass:           true,
			ExpErr:            nil,
			ExpectConsumedGas: expectedConsumedGasResolve,
			Mock:              false,
		},
		{
			Name: "signer has enough funds, should pass, no tip - resolvable denom - simulate",
			Malleate: func(s *postsuite.TestSuite) postsuite.TestCaseArgs {
				accs := s.CreateTestAccounts(1)

				balance := postsuite.TestAccountBalance{
					TestAccount: accs[0],
					Coins:       validResolvableFee,
				}
				s.SetAccountBalances([]postsuite.TestAccountBalance{balance})

				return postsuite.TestCaseArgs{
					Msgs:      []sdk.Msg{testdata.NewTestMsg(accs[0].Account.GetAddress())},
					GasLimit:  gasLimit,
					FeeAmount: validResolvableFee,
				}
			},
			RunAnte:           true,
			RunPost:           true,
			Simulate:          true,
			ExpPass:           true,
			ExpErr:            nil,
			ExpectConsumedGas: expectedConsumedGas,
			Mock:              false,
		},
		{
			Name: "signer has no balance, should pass, no tip - resolvable denom - simulate",
			Malleate: func(s *postsuite.TestSuite) postsuite.TestCaseArgs {
				accs := s.CreateTestAccounts(1)

				return postsuite.TestCaseArgs{
					Msgs:      []sdk.Msg{testdata.NewTestMsg(accs[0].Account.GetAddress())},
					GasLimit:  gasLimit,
					FeeAmount: validResolvableFee,
				}
			},
			RunAnte:           true,
			RunPost:           true,
			Simulate:          true,
			ExpPass:           true,
			ExpErr:            nil,
			ExpectConsumedGas: expectedConsumedGas,
			Mock:              false,
		},
		{
			Name: "0 gas given should pass in simulate - no fee",
			Malleate: func(s *postsuite.TestSuite) postsuite.TestCaseArgs {
				accs := s.CreateTestAccounts(1)

				return postsuite.TestCaseArgs{
					Msgs:      []sdk.Msg{testdata.NewTestMsg(accs[0].Account.GetAddress())},
					GasLimit:  0,
					FeeAmount: nil,
				}
			},
			RunAnte:           true,
			RunPost:           false,
			Simulate:          true,
			ExpPass:           true,
			ExpErr:            nil,
			ExpectConsumedGas: expectedConsumedGas,
			Mock:              false,
		},
		{
			Name: "0 gas given should pass in simulate - fee",
			Malleate: func(s *postsuite.TestSuite) postsuite.TestCaseArgs {
				accs := s.CreateTestAccounts(1)

				return postsuite.TestCaseArgs{
					Msgs:      []sdk.Msg{testdata.NewTestMsg(accs[0].Account.GetAddress())},
					GasLimit:  0,
					FeeAmount: validFee,
				}
			},
			RunAnte:           true,
			RunPost:           false,
			Simulate:          true,
			ExpPass:           true,
			ExpErr:            nil,
			ExpectConsumedGas: expectedConsumedGas,
			Mock:              false,
		},
		{
			Name: "no fee - fail",
			Malleate: func(s *postsuite.TestSuite) postsuite.TestCaseArgs {
				accs := s.CreateTestAccounts(1)

				return postsuite.TestCaseArgs{
					Msgs:      []sdk.Msg{testdata.NewTestMsg(accs[0].Account.GetAddress())},
					GasLimit:  1000000000,
					FeeAmount: nil,
				}
			},
			RunAnte:  true,
			RunPost:  true,
			Simulate: false,
			ExpPass:  false,
			ExpErr:   types.ErrNoFeeCoins,
			Mock:     false,
		},
		{
			Name: "no gas limit - fail",
			Malleate: func(s *postsuite.TestSuite) postsuite.TestCaseArgs {
				accs := s.CreateTestAccounts(1)

				return postsuite.TestCaseArgs{
					Msgs:      []sdk.Msg{testdata.NewTestMsg(accs[0].Account.GetAddress())},
					GasLimit:  0,
					FeeAmount: nil,
				}
			},
			RunAnte:  true,
			RunPost:  true,
			Simulate: false,
			ExpPass:  false,
			ExpErr:   sdkerrors.ErrOutOfGas,
			Mock:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Case %s", tc.Name), func(t *testing.T) {
			s := postsuite.SetupTestSuite(t, tc.Mock)
			s.TxBuilder = s.ClientCtx.TxConfig.NewTxBuilder()
			args := tc.Malleate(s)

			s.RunTestCase(t, tc, args)
		})
	}
}
