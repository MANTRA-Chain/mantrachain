package slinky_test

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/icza/dyno"
	marketmaptypes "github.com/skip-mev/connect/v2/x/marketmap/types"
	"github.com/skip-mev/connect/v2/x/oracle"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/stretchr/testify/suite"

	integration "github.com/skip-mev/connect/tests/integration/v2"
	marketmapmodule "github.com/skip-mev/connect/v2/x/marketmap"
)

var (
	image = ibc.DockerImage{
		Repository: "mantra-chain/mantrachain",
		Version:    "local",
		UidGid:     "1025:1025",
	}

	numValidators = 4
	numFullNodes  = 0
	noHostMount   = false
	gasAdjustment = 1.5

	oracleImage = ibc.DockerImage{
		Repository: "ghcr.io/skip-mev/connect-sidecar",
		Version:    "v2.0.1",
		UidGid:     "1000:1000",
	}
	encodingConfig = testutil.MakeTestEncodingConfig(
		bank.AppModuleBasic{},
		oracle.AppModuleBasic{},
		gov.AppModuleBasic{},
		auth.AppModuleBasic{},
		marketmapmodule.AppModuleBasic{},
	)

	//nolint:unused // may need to be updated to reflect mantra genesis
	defaultGenesis = marketmaptypes.DefaultGenesisState()
	//nolint:unused // may need to be updated to reflect mantra genesis
	govAddr = authtypes.NewModuleAddress(govtypes.ModuleName).String()
	//nolint:unused // may need to be updated to reflect mantra genesis
	params = marketmaptypes.Params{
		MarketAuthorities: []string{govAddr},
		Admin:             govAddr,
	}
	defaultGenesisKV = []cosmos.GenesisKV{
		{
			Key:   "consensus.params.abci.vote_extensions_enable_height",
			Value: "2",
		},
		{
			Key:   "consensus.params.block.max_gas",
			Value: "1000000000",
		},
		{
			Key:   "app_state.feemarket.params.enabled",
			Value: false,
		},
	}

	denom = "uom"
	spec  = &interchaintest.ChainSpec{
		ChainName:     "slinky",
		Name:          "slinky",
		NumValidators: &numValidators,
		NumFullNodes:  &numFullNodes,
		Version:       "local",
		NoHostMount:   &noHostMount,
		ChainConfig: ibc.ChainConfig{
			EncodingConfig: &encodingConfig,
			Images: []ibc.DockerImage{
				image,
			},
			Type:           "cosmos",
			Name:           "slinky",
			Denom:          denom,
			ChainID:        "chain-id-0",
			Bin:            "mantrachaind",
			Bech32Prefix:   "mantra",
			CoinType:       "118",
			GasAdjustment:  gasAdjustment,
			GasPrices:      fmt.Sprintf("0%s", denom),
			TrustingPeriod: "48h",
			NoHostMount:    noHostMount,
			ModifyGenesis:  ModifyGenesis(defaultGenesisKV),
		},
	}
)

func TestSlinkyOracleIntegration(t *testing.T) {
	baseSuite := integration.NewConnectIntegrationSuite(
		spec,
		oracleImage,
		integration.WithDenom(denom),
	)

	suite.Run(t, integration.NewSlinkyOracleIntegrationSuite(baseSuite))
}

func ModifyGenesis(genesisKV []cosmos.GenesisKV) func(ibc.ChainConfig, []byte) ([]byte, error) {
	return func(chainConfig ibc.ChainConfig, genbz []byte) ([]byte, error) {
		g := make(map[string]interface{})
		if err := json.Unmarshal(genbz, &g); err != nil {
			return nil, fmt.Errorf("failed to unmarshal genesis file: %w", err)
		}

		for idx, values := range genesisKV {
			splitPath := strings.Split(values.Key, ".")

			path := make([]interface{}, len(splitPath))
			for i, component := range splitPath {
				if v, err := strconv.Atoi(component); err == nil {
					path[i] = v
				} else {
					path[i] = component
				}
			}

			if err := dyno.Set(g, values.Value, path...); err != nil {
				return nil, fmt.Errorf("failed to set key '%s' as '%+v' (index:%d) in genesis json: %w", values.Key, values.Value, idx, err)
			}
		}

		out, err := json.Marshal(g)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal genesis bytes to json: %w", err)
		}

		return out, nil
	}
}
