package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	keepertest "github.com/MANTRA-Chain/mantrachain/testutil/keeper"
	"github.com/MANTRA-Chain/mantrachain/x/tax/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestAllocateMcaTax(t *testing.T) {
	k, ctx, _ := keepertest.TaxKeeper(t)

	tests := []struct {
		name           string
		mcaTax         math.LegacyDec
		mcaAddress     string
		initialBalance sdk.Coins
		expectedError  bool
	}{
		{
			name:           "Successful allocation",
			mcaTax:         math.LegacyNewDecWithPrec(1, 1), // 0.1
			mcaAddress:     "mantra15m77x4pe6w9vtpuqm22qxu0ds7vn4ehzwx8pls",
			initialBalance: sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(1000))),
			expectedError:  false,
		},
		{
			name:           "Zero MCA tax",
			mcaTax:         math.LegacyZeroDec(),
			mcaAddress:     "mantra15m77x4pe6w9vtpuqm22qxu0ds7vn4ehzwx8pls",
			initialBalance: sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(1000))),
			expectedError:  false,
		},
		{
			name:           "Invalid MCA address",
			mcaTax:         math.LegacyNewDecWithPrec(1, 1),
			mcaAddress:     "invalid_address",
			initialBalance: sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(1000))),
			expectedError:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx = ctx.WithBlockHeight(2) // Ensure we're not in the first block

			// Set up initial balance
			feeCollectorAddr := k.GetFeeCollectorAddress()
			err := k.BankKeeper.MintCoins(ctx, types.ModuleName, tc.initialBalance)
			require.NoError(t, err)
			err = k.BankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.GetFeeCollectorName(), tc.initialBalance)
			require.NoError(t, err)

			mcaAddress, _ := sdk.AccAddressFromBech32(tc.mcaAddress)
			err = k.AllocateMcaTax(ctx, tc.mcaTax, mcaAddress)

			if tc.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				// Check if the MCA address received the correct amount
				expectedAllocation := tc.initialBalance.AmountOf("stake").ToDec().Mul(tc.mcaTax).TruncateInt()
				mcaBalance := k.BankKeeper.GetBalance(ctx, mcaAddress, "stake")
				require.Equal(t, expectedAllocation, mcaBalance.Amount)

				// Check if the fee collector's balance was reduced correctly
				feeCollectorBalance := k.BankKeeper.GetBalance(ctx, feeCollectorAddr, "stake")
				expectedFeeCollectorBalance := tc.initialBalance.AmountOf("stake").Sub(expectedAllocation)
				require.Equal(t, expectedFeeCollectorBalance, feeCollectorBalance.Amount)
			}
		})
	}
}
