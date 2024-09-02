package types_test

import (
	"testing"

	utils "github.com/MANTRA-Finance/mantrachain/types"
	"github.com/MANTRA-Finance/mantrachain/x/lpfarm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestParams_Validate(t *testing.T) {
	for _, tc := range []struct {
		name        string
		malleate    func(params *types.Params)
		expectedErr string // empty means no error
	}{
		{
			"valid params",
			func(params *types.Params) {},
			"",
		},
		{
			"invalid private plan creation fee",
			func(params *types.Params) {
				params.PrivatePlanCreationFee = sdk.Coins{utils.ParseCoin("0stake")}
			},
			"invalid private plan creation fee: coin 0stake amount is not positive",
		},
		{
			"invalid fee collector",
			func(params *types.Params) {
				params.FeeCollector = invalidAddr
			},
			"invalid fee collector address: invalidaddr",
		},
		{
			"zero max num private plans",
			func(params *types.Params) {
				params.MaxNumPrivatePlans = 0
			},
			"",
		},
		{
			"zero max block duration",
			func(params *types.Params) {
				params.MaxBlockDuration = 0
			},
			"max block duration must be positive",
		},
		{
			"negative max block duration",
			func(params *types.Params) {
				params.MaxBlockDuration = -1
			},
			"max block duration must be positive",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			params := types.DefaultParams()
			tc.malleate(&params)
			err := params.Validate()
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}
