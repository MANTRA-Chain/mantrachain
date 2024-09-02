package types_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/MANTRA-Finance/mantrachain/x/marketmaker/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"
)

func TestParams(t *testing.T) {
	require.IsType(t, paramstypes.KeyTable{}, types.ParamKeyTable())

	minOpenRatio, err := sdkmath.LegacyNewDecFromStr("0.500000000000000000")
	require.NoError(t, err)
	minOpenDepthRation, err := sdkmath.LegacyNewDecFromStr("0.100000000000000000")
	require.NoError(t, err)

	wantParams := types.Params{
		IncentiveBudgetAddress: "cosmos1ddn66jv0sjpmck0ptegmhmqtn35qsg2vxyk2hn9sqf4qxtzqz3sqanrtcm",
		DepositAmount: sdk.Coins{
			{
				Denom:  "stake",
				Amount: sdkmath.NewInt(1000000000),
			},
		},
		Common: types.Common{
			MinOpenRatio:      minOpenRatio,
			MinOpenDepthRatio: minOpenDepthRation,
			MaxDowntime:       20,
			MaxTotalDowntime:  100,
			MinHours:          16,
			MinDays:           22,
		},
		IncentivePairs: []types.IncentivePair{},
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
			"Default",
			func(params *types.Params) {
			},
			"",
		},
		{
			"EmptyDepositAmount",
			func(params *types.Params) {
				params.DepositAmount = sdk.NewCoins()
			},
			"",
		},
		{
			"NegativeDepositAmount",
			func(params *types.Params) {
				params.DepositAmount = sdk.Coins{
					sdk.Coin{
						Denom:  "stake",
						Amount: sdkmath.NewInt(-1),
					},
				}
			},
			"coin -1stake amount is not positive",
		},
		{
			"EmptyBudgetAddr",
			func(params *types.Params) {
				params.IncentiveBudgetAddress = ""
			},
			"incentive budget address must not be empty",
		},
		{
			"WrongBudgetAddr",
			func(params *types.Params) {
				params.IncentiveBudgetAddress = "addr1"
			},
			"invalid account address: addr1",
		},
		{
			"IncentivePair",
			func(params *types.Params) {
				params.IncentivePairs = []types.IncentivePair{
					{
						PairId:          1,
						IncentiveWeight: sdkmath.LegacyMustNewDecFromStr("0.1"),
					},
				}
				a, _ := json.Marshal(params.IncentivePairs[0])
				fmt.Println(string(a))
			},
			"",
		},
		{
			"DuplicatedIncentivePair",
			func(params *types.Params) {
				params.IncentivePairs = []types.IncentivePair{
					{
						PairId:          1,
						IncentiveWeight: sdkmath.LegacyMustNewDecFromStr("0.1"),
					},
					{
						PairId:          1,
						IncentiveWeight: sdkmath.LegacyMustNewDecFromStr("0.2"),
					},
				}
			},
			"incentive pair id cannot be duplicated: 1",
		},
		{
			"MultipleIncentivePairs",
			func(params *types.Params) {
				params.IncentivePairs = []types.IncentivePair{
					{
						PairId:          1,
						IncentiveWeight: sdkmath.LegacyMustNewDecFromStr("0.1"),
					},
					{
						PairId:          2,
						IncentiveWeight: sdkmath.LegacyMustNewDecFromStr("0.2"),
					},
				}
			},
			"",
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
