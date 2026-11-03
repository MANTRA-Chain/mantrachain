package upgrades_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MANTRA-Chain/mantrachain/v8/app/upgrades"
)

func TestHandlerNames(t *testing.T) {
	testCases := []struct {
		name     string
		upgrades []upgrades.Upgrade
		expNames []string
	}{
		{
			name:     "no upgrades",
			upgrades: nil,
			expNames: []string{},
		},
		{
			name: "names are sorted",
			upgrades: []upgrades.Upgrade{
				{UpgradeName: "v10"},
				{UpgradeName: "v8"},
				{UpgradeName: "v9"},
			},
			expNames: []string{"v10", "v8", "v9"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expNames, upgrades.HandlerNames(tc.upgrades))
		})
	}
}
