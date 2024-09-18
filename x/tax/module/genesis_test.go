package tax_test

import (
	"testing"

	keepertest "github.com/MANTRA-Chain/mantrachain/testutil/keeper"
	"github.com/MANTRA-Chain/mantrachain/testutil/nullify"
	tax "github.com/MANTRA-Chain/mantrachain/x/tax/module"
	"github.com/MANTRA-Chain/mantrachain/x/tax/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx, _ := keepertest.TaxKeeper(t)
	err := tax.InitGenesis(ctx, k, genesisState)
	require.NoError(t, err)
	got, err := tax.ExportGenesis(ctx, k)
	require.NoError(t, err)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}

func TestGenesisState_Validate(t *testing.T) {
	tests := []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
				Params: types.NewParams("0.1", "mantra15m77x4pe6w9vtpuqm22qxu0ds7vn4ehzwx8pls"),
			},
			valid: true,
		},
		{
			desc: "negative proportion is invalid",
			genState: &types.GenesisState{
				Params: types.NewParams("-0.5", types.DefaultMcaAddress),
			},
			valid: false,
		},
		{
			desc: "proportion greater than 1 is invalid",
			genState: &types.GenesisState{
				Params: types.NewParams("1.5", types.DefaultMcaAddress),
			},
			valid: false,
		},
		{
			desc: "invalid bech32 address",
			genState: &types.GenesisState{
				Params: types.NewParams(types.DefaultProportion, "invalid_address"),
			},
			valid: false,
		},
		{
			desc: "empty mca address is invalid",
			genState: &types.GenesisState{
				Params: types.NewParams(types.DefaultProportion, ""),
			},
			valid: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		params  types.Params
		wantErr bool
	}{
		{
			name:    "default params",
			params:  types.DefaultParams(),
			wantErr: false,
		},
		{
			name:    "valid params",
			params:  types.NewParams("0.1", "mantra15m77x4pe6w9vtpuqm22qxu0ds7vn4ehzwx8pls"),
			wantErr: false,
		},
		{
			name:    "invalid proportion",
			params:  types.NewParams("-0.1", types.DefaultMcaAddress),
			wantErr: true,
		},
		{
			name:    "invalid mca address",
			params:  types.NewParams(types.DefaultProportion, "invalid_address"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestParams_ValidateProportion(t *testing.T) {
	tests := []struct {
		name       string
		proportion string
		wantErr    bool
	}{
		{"valid proportion", "0.1", false},
		{"zero proportion", "0", false},
		{"max proportion", "1", false},
		{"negative proportion", "-0.1", true},
		{"proportion greater than 1", "1.1", true},
		{"invalid format", "abc", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := types.ValidateMcaTax(tt.proportion)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestParams_ValidateMcaAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
		wantErr bool
	}{
		{"valid address", "cosmos15m77x4pe6w9vtpuqm22qxu0ds7vn4ehz9dd9u2", false},
		{"empty address", "", true},
		{"invalid bech32", "invalid_address", true},
		{"wrong prefix", "cosmos1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqnrql8a", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := types.ValidateMcaAddress(tt.address)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
