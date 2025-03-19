package e2e

import (
	"fmt"
	"os"

	"cosmossdk.io/log"
	evidencetypes "cosmossdk.io/x/evidence/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/MANTRA-Chain/mantrachain/v4/app"
	"github.com/MANTRA-Chain/mantrachain/v4/app/params"
	sanctiontypes "github.com/MANTRA-Chain/mantrachain/v4/x/sanction/types"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distribtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govv1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	paramsproptypes "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ratelimittypes "github.com/cosmos/ibc-apps/modules/rate-limiting/v8/types"
)

const (
	keyringPassphrase = "testpassphrase"
	keyringAppName    = "testnet"
)

var (
	encodingConfig params.EncodingConfig
	cdc            codec.Codec
	txConfig       client.TxConfig
)

func init() {
	encodingConfig = params.MakeEncodingConfig()
	banktypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	authtypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	authvesting.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	stakingtypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	evidencetypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	cryptocodec.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	govv1types.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	govv1beta1types.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	paramsproptypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	paramsproptypes.RegisterLegacyAminoCodec(encodingConfig.Amino)

	upgradetypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	distribtypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ratelimittypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	sanctiontypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	cdc = encodingConfig.Codec
	txConfig = encodingConfig.TxConfig
}

type chain struct {
	dataDir    string
	id         string
	validators []*validator
	accounts   []*account //nolint:unused
	// initial accounts in genesis
	genesisAccounts        []*account
	genesisVestingAccounts map[string]sdk.AccAddress
}

func newChain() (*chain, error) {
	tmpDir, err := os.MkdirTemp("", "app-e2e-testnet-")
	if err != nil {
		return nil, err
	}

	return &chain{
		id:      "chain-" + tmrand.Str(6),
		dataDir: tmpDir,
	}, nil
}

func (c *chain) configDir() string {
	return fmt.Sprintf("%s/%s", c.dataDir, c.id)
}

func (c *chain) createAndInitValidators(count int) error {
	var emptyWasmOpts []wasmkeeper.Option
	tempApplication := app.New(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		nil,
		true,
		app.EmptyAppOptions{},
		emptyWasmOpts,
		app.NoOpEVMOptions,
	)
	defer func() {
		if err := tempApplication.Close(); err != nil {
			panic(err)
		}
	}()

	genesisState := tempApplication.DefaultGenesis()

	for i := 0; i < count; i++ {
		node := c.createValidator(i)

		// generate genesis files
		if err := node.init(genesisState); err != nil {
			return err
		}

		c.validators = append(c.validators, node)

		// create keys
		if err := node.createKey("val"); err != nil {
			return err
		}
		if err := node.createNodeKey(); err != nil {
			return err
		}
		if err := node.createConsensusKey(); err != nil {
			return err
		}
	}

	return nil
}

func (c *chain) createAndInitValidatorsWithMnemonics(count int, mnemonics []string) error { //nolint:unused // this is called during e2e tests
	var emptyWasmOpts []wasmkeeper.Option
	tempApplication := app.New(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		nil,
		true,
		app.EmptyAppOptions{},
		emptyWasmOpts,
		app.NoOpEVMOptions,
	)
	defer func() {
		if err := tempApplication.Close(); err != nil {
			panic(err)
		}
	}()

	genesisState := tempApplication.DefaultGenesis()

	for i := 0; i < count; i++ {
		// create node
		node := c.createValidator(i)

		// generate genesis files
		if err := node.init(genesisState); err != nil {
			return err
		}

		c.validators = append(c.validators, node)

		// create keys
		if err := node.createKeyFromMnemonic("val", mnemonics[i]); err != nil {
			return err
		}
		if err := node.createNodeKey(); err != nil {
			return err
		}
		if err := node.createConsensusKey(); err != nil {
			return err
		}
	}

	return nil
}

func (c *chain) createValidator(index int) *validator {
	return &validator{
		chain:   c,
		index:   index,
		moniker: fmt.Sprintf("%s-app-%d", c.id, index),
	}
}
