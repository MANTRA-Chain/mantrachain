package cli_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/testutil"
	"github.com/stretchr/testify/require"

	mantratestutil "mantrachain/testutil"
	"mantrachain/x/lpfarm/client/cli"
)

func TestParseFarmingPlanProposal(t *testing.T) {
	okJSON := testutil.WriteToNewTempFile(t, `
{
  "title": "Farming Plan Proposal",
  "description": "Let's start farming",
  "create_plan_requests": [
    {
      "description": "New Farming Plan",
      "farming_pool_address": "mantra1t3g4vylrgun8k4wm5dlw8hmcn5x0p6jvknh550",
      "reward_allocations": [
        {
          "pair_id": "1",
          "rewards_per_day": [
            {
              "denom": "stake",
              "amount": "100000000"
            }
          ]
        },
        {
          "denom": "pool2",
          "rewards_per_day": [
            {
              "denom": "stake",
              "amount": "200000000"
            }
          ]
        }
      ],
      "start_time": "2022-01-01T00:00:00Z",
      "end_time": "2023-01-01T00:00:00Z"
    }
  ],
  "terminate_plan_requests": [
    {
      "plan_id": "1"
    },
    {
      "plan_id": "2"
    }
  ]
}
`)

	encodingConfig := mantratestutil.MakeTestEncodingConfig()

	plan, err := cli.ParseFarmingPlanProposal(encodingConfig.Marshaler, okJSON.Name())
	require.NoError(t, err)
	require.NotEmpty(t, plan.String())
}
