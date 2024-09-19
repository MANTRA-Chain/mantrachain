package tax_test

import (
	"testing"

	keepertest "github.com/MANTRA-Chain/mantrachain/testutil/keeper"
	module "github.com/MANTRA-Chain/mantrachain/x/tax/module"
	"github.com/MANTRA-Chain/mantrachain/x/tax/types"
	"github.com/stretchr/testify/require"
)

func TestBeginBlocker(t *testing.T) {
	k, ctx, _ := keepertest.TaxKeeper(t)

	tests := []struct {
		name          string
		blockHeight   int64
		mcaTax        string
		mcaAddress    string
		expectedError bool
	}{
		{
			name:          "First block, no allocation",
			blockHeight:   1,
			mcaTax:        "0.1",
			mcaAddress:    "mantra15m77x4pe6w9vtpuqm22qxu0ds7vn4ehzwx8pls",
			expectedError: false,
		},
		{
			name:          "Subsequent block, allocation occurs",
			blockHeight:   2,
			mcaTax:        "0.1",
			mcaAddress:    "mantra15m77x4pe6w9vtpuqm22qxu0ds7vn4ehzwx8pls",
			expectedError: false,
		},
		{
			name:          "Zero MCA tax, no allocation",
			blockHeight:   2,
			mcaTax:        "0",
			mcaAddress:    "mantra15m77x4pe6w9vtpuqm22qxu0ds7vn4ehzwx8pls",
			expectedError: false,
		},
		{
			name:          "Invalid MCA address",
			blockHeight:   2,
			mcaTax:        "0.1",
			mcaAddress:    "invalid_address",
			expectedError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx = ctx.WithBlockHeight(tc.blockHeight)

			params := types.NewParams(tc.mcaTax, tc.mcaAddress)
			err := k.Params.Set(ctx, params)
			require.NoError(t, err)

			err = module.BeginBlocker(ctx, k)

			if tc.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
