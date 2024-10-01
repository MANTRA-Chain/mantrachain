package ante

import (
	"testing"

	"github.com/osmosis-labs/fee-abstraction/v8/x/feeabs/types"
	feemarketante "github.com/skip-mev/feemarket/x/feemarket/ante"
	feemarkettypes "github.com/skip-mev/feemarket/x/feemarket/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"cosmossdk.io/errors"
	math "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func TestMempoolDecorator(t *testing.T) {
	gasLimit := uint64(200000)
	minGasPrice := sdk.NewDecCoinsFromCoins(sdk.NewInt64Coin("ueve", feemarkettypes.DefaultMinBaseGasPrice.TruncateInt64()))
	validFeeAmount := feemarkettypes.DefaultMinBaseGasPrice.MulInt64(int64(gasLimit))
	validFee := sdk.NewCoins(sdk.NewCoin("ueve", validFeeAmount.TruncateInt()))
	validIbcFee := sdk.NewCoins(sdk.NewCoin("ibcfee", validFeeAmount.TruncateInt()))
	// mockHostZoneConfig is used to mock the host zone config, with ibcfee as the ibc fee denom to be used as alternative fee
	mockHostZoneConfig := types.HostChainFeeAbsConfig{
		IbcDenom:                "ibcfee",
		OsmosisPoolTokenDenomIn: "osmosis",
		PoolId:                  1,
		Status:                  types.HostChainFeeAbsStatus_UPDATED,
	}
	testCases := []struct {
		name      string
		feeAmount sdk.Coins
		malleate  func(*AnteTestSuite)
		expErr    error
	}{
		{
			"empty fee, should fail",
			sdk.Coins{},
			func(suite *AnteTestSuite) {
			},
			errors.Wrapf(feemarkettypes.ErrNoFeeCoins, "%s", "got length 0"),
		},
		{
			"valid native fee, should pass",
			validFee,
			func(suite *AnteTestSuite) {
				suite.bankKeeper.On("SendCoinsFromAccountToModule", mock.Anything, mock.Anything,
					feemarkettypes.FeeCollectorName, mock.Anything).Return(nil).Once()
			},
			nil,
		},
		{
			"valid ibc fee, should pass",
			validIbcFee,
			func(suite *AnteTestSuite) {
				err := suite.feeabsKeeper.SetHostZoneConfig(suite.ctx, mockHostZoneConfig)
				require.NoError(t, err)
				suite.feeabsKeeper.SetTwapRate(suite.ctx, "ibcfee", math.LegacyNewDec(1))
				suite.stakingKeeper.EXPECT().BondDenom(gomock.Any()).Return("ueve", nil).AnyTimes()
				suite.bankKeeper.On("SendCoinsFromAccountToModule", mock.Anything, mock.Anything,
					feemarkettypes.FeeCollectorName, mock.Anything).Return(nil).Once()
			},
			nil,
		},
		{
			"not enough ibc fee, should fail",
			validIbcFee.Sub(sdk.NewCoin("ibcfee", math.NewInt(1))),
			func(suite *AnteTestSuite) {
				err := suite.feeabsKeeper.SetHostZoneConfig(suite.ctx, mockHostZoneConfig)
				require.NoError(t, err)
				suite.feeabsKeeper.SetTwapRate(suite.ctx, "ibcfee", math.LegacyNewDec(1))
				suite.stakingKeeper.EXPECT().BondDenom(gomock.Any()).Return("ueve", nil).AnyTimes()
			},
			sdkerrors.ErrInsufficientFee,
		},
		{
			"fee in unsupported denom, should fail",
			sdk.NewCoins(sdk.NewCoin("unsupported", validFeeAmount.TruncateInt())),
			func(suite *AnteTestSuite) {
				suite.stakingKeeper.EXPECT().BondDenom(gomock.Any()).Return("ueve", nil).AnyTimes()
			},
			ErrDenomNotRegistered("unsupported"),
		},
		{
			"multiple fee denoms, only one supported, should pass",
			sdk.NewCoins(validFee[0], sdk.NewCoin("unsupported", math.NewInt(100))),
			func(suite *AnteTestSuite) {},
			feemarkettypes.ErrTooManyFeeCoins,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			suite := SetupTestSuite(t, true)

			tc.malleate(suite)
			suite.txBuilder.SetGasLimit(gasLimit)
			suite.txBuilder.SetFeeAmount(tc.feeAmount)
			accs := suite.CreateTestAccounts(1)
			require.NoError(t, suite.txBuilder.SetMsgs([]sdk.Msg{testdata.NewTestMsg(accs[0].acc.GetAddress())}...))

			suite.ctx = suite.ctx.WithMinGasPrices(minGasPrice)

			// Construct tx and run through mempool decorator
			tx := suite.txBuilder.GetTx()
			feemarketDecorator := feemarketante.NewFeeMarketCheckDecorator(
				suite.accountKeeper,
				suite.bankKeeper,
				suite.feeGrantKeeper,
				suite.feemarketKeeper,
				nil)
			antehandler := sdk.ChainAnteDecorators(feemarketDecorator)

			// Run the ante handler
			_, err := antehandler(suite.ctx, tx, false)

			if tc.expErr != nil {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
