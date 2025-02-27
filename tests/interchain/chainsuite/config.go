package chainsuite

import (
	"fmt"
	"strconv"
	"time"

	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/strangelove-ventures/interchaintest/v8/testutil"

	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/types"
)

type ChainScope int

const (
	ChainScopeSuite ChainScope = iota
	ChainScopeTest  ChainScope = iota
)

type SuiteConfig struct {
	ChainSpec      *interchaintest.ChainSpec
	UpgradeOnSetup bool
	CreateRelayer  bool
	Scope          ChainScope
}

const (
	CommitTimeout          = 4 * time.Second
	Uom                    = "uom"
	GovMinDepositAmount    = 1000
	GovDepositAmount       = "5000000" + Uom
	GovDepositPeriod       = 60 * time.Second
	GovVotingPeriod        = 40 * time.Second
	DowntimeJailDuration   = 10 * time.Second
	ProviderSlashingWindow = 10
	GasPrices              = "0.01" + Uom
	ValidatorCount         = 6
	UpgradeDelta           = 12
	ValidatorFunds         = 11_000_000_000
	ChainSpawnWait         = 155 * time.Second
	SlashingWindowConsumer = 20
	BlocksPerDistribution  = 10
)

func (c SuiteConfig) Merge(other SuiteConfig) SuiteConfig {
	if c.ChainSpec == nil {
		c.ChainSpec = other.ChainSpec
	} else if other.ChainSpec != nil {
		c.ChainSpec.ChainConfig = c.ChainSpec.MergeChainSpecConfig(other.ChainSpec.ChainConfig)
		if other.ChainSpec.Name != "" {
			c.ChainSpec.Name = other.ChainSpec.Name
		}
		if other.ChainSpec.ChainName != "" {
			c.ChainSpec.ChainName = other.ChainSpec.ChainName
		}
		if other.ChainSpec.Version != "" {
			c.ChainSpec.Version = other.ChainSpec.Version
		}
		if other.ChainSpec.NoHostMount != nil {
			c.ChainSpec.NoHostMount = other.ChainSpec.NoHostMount
		}
		if other.ChainSpec.NumValidators != nil {
			c.ChainSpec.NumValidators = other.ChainSpec.NumValidators
		}
		if other.ChainSpec.NumFullNodes != nil {
			c.ChainSpec.NumFullNodes = other.ChainSpec.NumFullNodes
		}
	}
	c.UpgradeOnSetup = other.UpgradeOnSetup
	c.CreateRelayer = other.CreateRelayer
	c.Scope = other.Scope
	return c
}

func DefaultGenesisAmounts(denom string) func(i int) (types.Coin, types.Coin) {
	return func(i int) (types.Coin, types.Coin) {
		return types.Coin{
				Denom:  denom,
				Amount: sdkmath.NewInt(ValidatorFunds),
			}, types.Coin{
				Denom: denom,
				Amount: sdkmath.NewInt([ValidatorCount]int64{
					30_000_000,
					29_000_000,
					20_000_000,
					10_000_000,
					7_000_000,
					4_000_000,
				}[i]),
			}
	}
}

func DefaultSuiteConfig(env Environment) SuiteConfig {
	fullNodes := 0
	validators := ValidatorCount
	var repository string
	if env.DockerRegistry == "" {
		repository = env.MantraImageName
	} else {
		repository = fmt.Sprintf("%s/%s", env.DockerRegistry, env.MantraImageName)
	}
	return SuiteConfig{
		ChainSpec: &interchaintest.ChainSpec{
			ChainName:     "mantra",
			Name:          "mantra",
			NumFullNodes:  &fullNodes,
			NumValidators: &validators,
			Version:       env.OldMantraImageVersion,
			ChainConfig: ibc.ChainConfig{
				Denom:         Uom,
				GasPrices:     GasPrices,
				Gas:           "auto",
				GasAdjustment: 2.0,
				ConfigFileOverrides: map[string]any{
					"config/config.toml": DefaultConfigToml(),
				},
				Images: []ibc.DockerImage{{
					Repository: repository,
					UidGid:     "1025:1025", // this is the user in heighliner docker images
				}},
				Type:                 "cosmos",
				Name:                 "mantra",
				ChainID:              "mantra-test-1",
				Bin:                  "mantrachaind",
				Bech32Prefix:         "mantra",
				CoinType:             "118",
				TrustingPeriod:       "48h",
				NoHostMount:          false,
				ModifyGenesis:        cosmos.ModifyGenesis(DefaultGenesis()),
				ModifyGenesisAmounts: DefaultGenesisAmounts(Uom),
			},
		},
	}
}

func DefaultConfigToml() testutil.Toml {
	configToml := make(testutil.Toml)
	consensusToml := make(testutil.Toml)
	consensusToml["timeout_commit"] = CommitTimeout
	configToml["consensus"] = consensusToml
	configToml["block_sync"] = false
	configToml["fast_sync"] = false
	return configToml
}

func DefaultGenesis() []cosmos.GenesisKV {
	return []cosmos.GenesisKV{
		cosmos.NewGenesisKV("app_state.gov.params.voting_period", GovVotingPeriod.String()),
		cosmos.NewGenesisKV("app_state.gov.params.max_deposit_period", GovDepositPeriod.String()),
		cosmos.NewGenesisKV("app_state.gov.params.min_deposit.0.denom", Uom),
		cosmos.NewGenesisKV("app_state.gov.params.min_deposit.0.amount", strconv.Itoa(GovMinDepositAmount)),
		cosmos.NewGenesisKV("app_state.slashing.params.signed_blocks_window", strconv.Itoa(ProviderSlashingWindow)),
		cosmos.NewGenesisKV("app_state.slashing.params.downtime_jail_duration", DowntimeJailDuration.String()),
		cosmos.NewGenesisKV("app_state.wasm.params.code_upload_access.permission", "Nobody"),
		cosmos.NewGenesisKV("app_state.wasm.params.instantiate_default_permission", "AnyOfAddresses"),
	}
}
