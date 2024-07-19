package types_test

import (
	"reflect"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/MANTRA-Finance/mantrachain/x/farming/types"
)

func TestParams(t *testing.T) {
	require.IsType(t, paramstypes.KeyTable{}, types.ParamKeyTable())

	wantParams := types.Params{
		PrivatePlanCreationFee: sdk.Coins{
			{
				Denom:  "stake",
				Amount: sdkmath.NewInt(1000000000),
			},
		},
		NextEpochDays:        1,
		FarmingFeeCollector:  "cosmos1h292smhhttwy0rl3qr4p6xsvpvxc4v05s6rxtczwq3cs6qc462mqejwy8x",
		DelayedStakingGasFee: 60000,
		MaxNumPrivatePlans:   10000,
	}

	require.Equal(t, wantParams, types.DefaultParams())
}

func TestParamsValidate(t *testing.T) {
	require.NoError(t, types.DefaultParams().Validate())

	testCases := []struct {
		name        string
		configure   func(*types.Params)
		expectedErr string
	}{
		{
			"EmptyPrivatePlanCreationFee",
			func(params *types.Params) {
				params.PrivatePlanCreationFee = sdk.NewCoins()
			},
			"",
		},
		{
			"ZeroNextEpochDays",
			func(params *types.Params) {
				params.NextEpochDays = uint32(0)
			},
			"next epoch days must be positive: 0",
		},
		{
			"EmptyFarmingFeeCollector",
			func(params *types.Params) {
				params.FarmingFeeCollector = ""
			},
			"farming fee collector address must not be empty",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			params := types.DefaultParams()
			tc.configure(&params)
			err := params.Validate()

			var err2 error
			for _, p := range params.ParamSetPairs() {
				err := p.ValidatorFn(reflect.ValueOf(p.Value).Elem().Interface())
				if err != nil {
					err2 = err
					break
				}
			}
			if tc.expectedErr != "" {
				require.EqualError(t, err, tc.expectedErr)
				require.EqualError(t, err2, tc.expectedErr)
			} else {
				require.Nil(t, err)
				require.Nil(t, err2)
			}
		})
	}
}
