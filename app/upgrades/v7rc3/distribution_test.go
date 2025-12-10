package v7rc3

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalPeriodData(t *testing.T) {
	testCases := []struct {
		name      string
		jsonData  []byte
		valAddr   string
		expPeriod uint64
		expRatio  string
	}{
		{
			name:      "Dryrun data",
			jsonData:  DryrunBeforeUpgrade,
			valAddr:   "mantravaloper143a99rce5u0p5l2tzy62hqjdl2dx6pmgxq4w78",
			expPeriod: 8005,
			expRatio:  "0.098805762397545508",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var dataBeforeUpgrade map[string]Period
			err := json.Unmarshal(tc.jsonData, &dataBeforeUpgrade)
			require.NoError(t, err, "unmarshaling should not fail")
			require.Equal(t, tc.expPeriod, dataBeforeUpgrade[tc.valAddr].Period)
			require.Equal(t, tc.expRatio, dataBeforeUpgrade[tc.valAddr].CumulativeRewardRatio)
		})
	}
}
