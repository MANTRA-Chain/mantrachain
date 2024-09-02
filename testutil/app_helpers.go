package testutil

import (
	"encoding/json"
	"fmt"
	"time"

	coreheader "cosmossdk.io/core/header"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	"github.com/MANTRA-Finance/mantrachain/app"
	utils "github.com/MANTRA-Finance/mantrachain/types"
	coinfactorytypes "github.com/MANTRA-Finance/mantrachain/x/coinfactory/types"
	guardtypes "github.com/MANTRA-Finance/mantrachain/x/guard/types"
	abci "github.com/cometbft/cometbft/abci/types"
	cmtjson "github.com/cometbft/cometbft/libs/json"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmttypes "github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/runtime"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

const (
	DefaultGenTxGas                           = 10000000
	TestAdminAddress                          = "mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka"
	TestAccountPrivilegesGuardNftCollectionId = "nft-guard-collection"
)

// DefaultConsensusParams defines the default CometBFT consensus params used in
// SimApp testing.
var DefaultConsensusParams = &cmtproto.ConsensusParams{
	Block: &cmtproto.BlockParams{
		MaxBytes: 200000,
		MaxGas:   100_000_000,
	},
	Evidence: &cmtproto.EvidenceParams{
		MaxAgeNumBlocks: 302400,
		MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration
		MaxBytes:        10000,
	},
	Validator: &cmtproto.ValidatorParams{
		PubKeyTypes: []string{
			cmttypes.ABCIPubKeyTypeEd25519,
		},
	},
}

func init() {
	InitSDKConfig()
}

func InitSDKConfig() {
	accountAddressPrefix := "mantra"
	accountPubKeyPrefix := accountAddressPrefix + "pub"
	validatorAddressPrefix := accountAddressPrefix + "valoper"
	validatorPubKeyPrefix := accountAddressPrefix + "valoperpub"
	consNodeAddressPrefix := accountAddressPrefix + "valcons"
	consNodePubKeyPrefix := accountAddressPrefix + "valconspub"

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(accountAddressPrefix, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
	config.Seal()
}

func MakeTestEncodingConfig() *codec.LegacyAmino {
	return codec.NewLegacyAmino()
}

// CreateRandomValidatorSet creates a validator set with one random validator
func CreateRandomValidatorSet() (*cmttypes.ValidatorSet, error) {
	privVal := mock.NewPV()
	pubKey, err := privVal.GetPubKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get pub key: %w", err)
	}

	// create validator set with single validator
	validator := cmttypes.NewValidator(pubKey, 1)

	return cmttypes.NewValidatorSet([]*cmttypes.Validator{validator}), nil
}

type GenesisAccount struct {
	authtypes.GenesisAccount
	Coins sdk.Coins
}

// StartupConfig defines the startup configuration new a test application.
//
// ValidatorSet defines a custom validator set to be validating the app.
// BaseAppOption defines the additional operations that must be run on baseapp before app start.
// AtGenesis defines if the app started should already have produced block or not.
type StartupConfig struct {
	ValidatorSet    func() (*cmttypes.ValidatorSet, error)
	BaseAppOption   runtime.BaseAppOption
	AtGenesis       bool
	GenesisAccounts []GenesisAccount
	DB              dbm.DB
}

func DefaultStartUpConfig() StartupConfig {
	priv := secp256k1.GenPrivKey()
	ba := authtypes.NewBaseAccount(priv.PubKey().Address().Bytes(), priv.PubKey(), 0, 0)
	ga := GenesisAccount{ba, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(100000000000000)))}
	return StartupConfig{
		ValidatorSet:    CreateRandomValidatorSet,
		AtGenesis:       false,
		GenesisAccounts: []GenesisAccount{ga},
		DB:              dbm.NewMemDB(),
	}
}

// NextBlock starts a new block.
func NextBlock(app *app.App, ctx sdk.Context, jumpTime time.Duration) (sdk.Context, error) {
	_, err := app.FinalizeBlock(&abci.RequestFinalizeBlock{Height: ctx.BlockHeight(), Time: ctx.BlockTime()})
	if err != nil {
		return sdk.Context{}, err
	}
	_, err = app.Commit()
	if err != nil {
		return sdk.Context{}, err
	}

	newBlockTime := ctx.BlockTime().Add(jumpTime)

	header := ctx.BlockHeader()
	header.Time = newBlockTime
	header.Height++

	newCtx := app.BaseApp.NewUncachedContext(false, header).WithHeaderInfo(coreheader.Info{
		Height: header.Height,
		Time:   header.Time,
	})

	return newCtx, err
}

// SetupWithConfiguration initializes a new app.App. A Nop logger is set in app.App.
// appConfig defines the application configuration (f.e. app_config.go).
// extraOutputs defines the extra outputs to be assigned by the dependency injector (depinject).
func Setup() (*app.App, error) {
	// create the app with depinject
	var (
		startupConfig = DefaultStartUpConfig()
		appConfig     = app.AppConfig()
		app           = &app.App{}
		appBuilder    *runtime.AppBuilder
		appCodec      codec.Codec
	)

	if err := depinject.Inject(
		depinject.Configs(
			appConfig,
			depinject.Supply(log.NewNopLogger()),
		),
		&appBuilder,
		&appCodec,
		&app.GuardKeeper,
		&app.AccountKeeper,
		&app.BankKeeper,
		&app.DistrKeeper,
		&app.StakingKeeper,
		&app.AirdropKeeper,
		&app.BridgeKeeper,
		&app.CoinfactoryKeeper,
		&app.DidKeeper,
		&app.TokenKeeper,
		&app.MarketmakerKeeper,
		&app.LiquidityKeeper,
		&app.LpfarmKeeper,
		&app.TxfeesKeeper,
	); err != nil {
		return nil, fmt.Errorf("failed to inject dependencies: %w", err)
	}

	app.SetAppCodec(appCodec)

	if startupConfig.BaseAppOption != nil {
		app.App = appBuilder.Build(startupConfig.DB, nil, startupConfig.BaseAppOption)
	} else {
		app.App = appBuilder.Build(startupConfig.DB, nil)
	}
	if err := app.Load(true); err != nil {
		return nil, fmt.Errorf("failed to load app: %w", err)
	}

	// create validator set
	valSet, err := startupConfig.ValidatorSet()
	if err != nil {
		return nil, fmt.Errorf("failed to create validator set")
	}

	var (
		balances    []banktypes.Balance
		genAccounts []authtypes.GenesisAccount
	)
	for _, ga := range startupConfig.GenesisAccounts {
		genAccounts = append(genAccounts, ga.GenesisAccount)
		balances = append(balances, banktypes.Balance{Address: ga.GenesisAccount.GetAddress().String(), Coins: ga.Coins})
	}

	genesisState, err := GenesisStateWithValSet(appCodec, app.DefaultGenesis(), valSet, genAccounts, balances...)
	if err != nil {
		return nil, fmt.Errorf("failed to create genesis state: %w", err)
	}

	// init chain must be called to stop deliverState from being nil
	stateBytes, err := cmtjson.MarshalIndent(genesisState, "", " ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal default genesis state: %w", err)
	}

	// init chain will set the validator set and initialize the genesis accounts
	_, err = app.InitChain(&abci.RequestInitChain{
		Validators:      []abci.ValidatorUpdate{},
		ConsensusParams: DefaultConsensusParams,
		AppStateBytes:   stateBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init chain: %w", err)
	}

	// commit genesis changes
	if !startupConfig.AtGenesis {
		_, err = app.FinalizeBlock(&abci.RequestFinalizeBlock{
			Height:             app.LastBlockHeight() + 1,
			NextValidatorsHash: valSet.Hash(),
			Time:               utils.ParseTime("2022-01-01T00:00:00Z"),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to finalize block: %w", err)
		}
	}

	return app, nil
}

// GenesisStateWithValSet returns a new genesis state with the validator set
func GenesisStateWithValSet(
	codec codec.Codec,
	genesisState map[string]json.RawMessage,
	valSet *cmttypes.ValidatorSet,
	genAccs []authtypes.GenesisAccount,
	balances ...banktypes.Balance,
) (map[string]json.RawMessage, error) {
	// set genesis accounts
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = codec.MustMarshalJSON(authGenesis)

	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

	bondAmt := sdk.DefaultPowerReduction

	for _, val := range valSet.Validators {
		pk, err := cryptocodec.FromCmtPubKeyInterface(val.PubKey)
		if err != nil {
			return nil, fmt.Errorf("failed to convert pubkey: %w", err)
		}

		pkAny, err := codectypes.NewAnyWithValue(pk)
		if err != nil {
			return nil, fmt.Errorf("failed to create new any: %w", err)
		}

		validator := stakingtypes.Validator{
			OperatorAddress:   sdk.ValAddress(val.Address).String(),
			ConsensusPubkey:   pkAny,
			Jailed:            false,
			Status:            stakingtypes.Bonded,
			Tokens:            bondAmt,
			DelegatorShares:   sdkmath.LegacyOneDec(),
			Description:       stakingtypes.Description{},
			UnbondingHeight:   int64(0),
			UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingtypes.NewCommission(sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec()),
			MinSelfDelegation: sdkmath.ZeroInt(),
		}
		validators = append(validators, validator)
		delegations = append(delegations, stakingtypes.NewDelegation(genAccs[0].GetAddress().String(), sdk.ValAddress(val.Address).String(), sdkmath.LegacyOneDec()))
	}

	// set validators and delegations
	stakingGenesis := stakingtypes.NewGenesisState(stakingtypes.DefaultParams(), validators, delegations)
	genesisState[stakingtypes.ModuleName] = codec.MustMarshalJSON(stakingGenesis)

	totalSupply := sdk.NewCoins()
	for _, b := range balances {
		// add genesis acc tokens to total supply
		totalSupply = totalSupply.Add(b.Coins...)
	}

	for range delegations {
		// add delegated tokens to total supply
		totalSupply = totalSupply.Add(sdk.NewCoin(sdk.DefaultBondDenom, bondAmt))
	}

	// add bonded amount to bonded pool module account
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, bondAmt)},
	})

	// update total supply
	bankGenesis := banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, balances, totalSupply, []banktypes.Metadata{}, []banktypes.SendEnabled{})
	genesisState[banktypes.ModuleName] = codec.MustMarshalJSON(bankGenesis)

	guardGenesis := guardtypes.NewGenesisState(guardtypes.NewParams(
		TestAdminAddress,
		TestAdminAddress,
		TestAccountPrivilegesGuardNftCollectionId,
		guardtypes.DefaultPrivileges,
		guardtypes.DefaultBaseDenom,
	))
	genesisState[guardtypes.ModuleName] = codec.MustMarshalJSON(guardGenesis)
	coinfactoryGenesis := coinfactorytypes.NewGenesisState(coinfactorytypes.NewParams(
		sdk.Coins{},
		0,
	))
	genesisState[coinfactorytypes.ModuleName] = codec.MustMarshalJSON(coinfactoryGenesis)

	return genesisState, nil
}

// EmptyAppOptions is a stub implementing AppOptions
type EmptyAppOptions struct{}

// Get implements AppOptions
func (ao EmptyAppOptions) Get(o string) interface{} {
	return nil
}

// AppOptionsMap is a stub implementing AppOptions which can get data from a map
type AppOptionsMap map[string]interface{}

func (m AppOptionsMap) Get(key string) interface{} {
	v, ok := m[key]
	if !ok {
		return interface{}(nil)
	}

	return v
}

func NewAppOptionsWithFlagHome(homePath string) servertypes.AppOptions {
	return AppOptionsMap{
		flags.FlagHome: homePath,
	}
}
