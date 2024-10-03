package tax_test

import (
	"testing"

	_ "github.com/MANTRA-Chain/mantrachain/app/params"
	keepertest "github.com/MANTRA-Chain/mantrachain/testutil/keeper"
	"github.com/MANTRA-Chain/mantrachain/testutil/nullify"
	tax "github.com/MANTRA-Chain/mantrachain/x/tax/module"
	"github.com/MANTRA-Chain/mantrachain/x/tax/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
	}

	k, ctx, _ := keepertest.TaxKeeper(t)
	err := tax.InitGenesis(ctx, k, genesisState)
	require.NoError(t, err)
	got, err := tax.ExportGenesis(ctx, k)
	require.NoError(t, err)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)
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
			desc: "negative MCA tax is invalid",
			genState: &types.GenesisState{
				Params: types.NewParams("-0.5", types.DefaultMcaAddress),
			},
			valid: false,
		},
		{
			desc: "MCA tax greater than 1 is invalid",
			genState: &types.GenesisState{
				Params: types.NewParams("1.5", types.DefaultMcaAddress),
			},
			valid: false,
		},
		{
			desc: "invalid bech32 address",
			genState: &types.GenesisState{
				Params: types.NewParams(types.DefaultMcaTax, "invalid_address"),
			},
			valid: false,
		},
		{
			desc: "empty mca address is invalid",
			genState: &types.GenesisState{
				Params: types.NewParams(types.DefaultMcaTax, ""),
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
			name:    "invalid mca tax",
			params:  types.NewParams("-0.1", types.DefaultMcaAddress),
			wantErr: true,
		},
		{
			name:    "invalid mca address",
			params:  types.NewParams(types.DefaultMcaTax, "invalid_address"),
			wantErr: true,
		},
		{
			name:    "mca tax too high",
			params:  types.NewParams("1.1", "mantra1qypqxpq9qcrsszg2pvxq6rs0zqg3yyc5lzutu9"),
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
		{"valid tax", "0.1", false},
		{"zero tax", "0", false},
		{"max tax", "0.3", false},
		{"negative tax", "-0.1", true},
		{"tax greater than 1", "1.1", true},
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
		{"valid address", "mantra15m77x4pe6w9vtpuqm22qxu0ds7vn4ehzwx8pls", false},
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
