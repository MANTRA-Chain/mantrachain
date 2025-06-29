package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	reflectionv1 "cosmossdk.io/api/cosmos/reflection/v1"
	"cosmossdk.io/client/v2/autocli"

	// Overriders
	clienthelpers "cosmossdk.io/client/v2/helpers"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/circuit"
	circuitkeeper "cosmossdk.io/x/circuit/keeper"
	circuittypes "cosmossdk.io/x/circuit/types"
	"cosmossdk.io/x/evidence"
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/feegrant"
	feegrantkeeper "cosmossdk.io/x/feegrant/keeper"
	feegrantmodule "cosmossdk.io/x/feegrant/module"
	"cosmossdk.io/x/nft"
	nftkeeper "cosmossdk.io/x/nft/keeper"
	nftmodule "cosmossdk.io/x/nft/module"
	"cosmossdk.io/x/upgrade"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	// Force-load the tracer engines to trigger registration due to Go-Ethereum v1.10.15 changes
	"github.com/MANTRA-Chain/mantrachain/v5/app/ante"
	_ "github.com/MANTRA-Chain/mantrachain/v5/app/params"
	queries "github.com/MANTRA-Chain/mantrachain/v5/app/queries"
	"github.com/MANTRA-Chain/mantrachain/v5/app/upgrades"
	v5 "github.com/MANTRA-Chain/mantrachain/v5/app/upgrades/v5"
	_ "github.com/MANTRA-Chain/mantrachain/v5/client/docs/statik"
	"github.com/MANTRA-Chain/mantrachain/v5/client/docs/swagger"
	sanctionkeeper "github.com/MANTRA-Chain/mantrachain/v5/x/sanction/keeper"
	sanction "github.com/MANTRA-Chain/mantrachain/v5/x/sanction/module"
	sanctiontypes "github.com/MANTRA-Chain/mantrachain/v5/x/sanction/types"
	taxkeeper "github.com/MANTRA-Chain/mantrachain/v5/x/tax/keeper"
	tax "github.com/MANTRA-Chain/mantrachain/v5/x/tax/module"
	taxtypes "github.com/MANTRA-Chain/mantrachain/v5/x/tax/types"
	"github.com/MANTRA-Chain/mantrachain/v5/x/tokenfactory"
	tokenfactorykeeper "github.com/MANTRA-Chain/mantrachain/v5/x/tokenfactory/keeper"
	tokenfactorytypes "github.com/MANTRA-Chain/mantrachain/v5/x/tokenfactory/types"
	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	runtimeservices "github.com/cosmos/cosmos-sdk/runtime/services"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	sigtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	txmodule "github.com/cosmos/cosmos-sdk/x/auth/tx/config"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/group"
	groupkeeper "github.com/cosmos/cosmos-sdk/x/group/keeper"
	groupmodule "github.com/cosmos/cosmos-sdk/x/group/module"
	"github.com/cosmos/cosmos-sdk/x/mint"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	evmosencoding "github.com/cosmos/evm/encoding"
	srvflags "github.com/cosmos/evm/server/flags"
	cosmosevmtypes "github.com/cosmos/evm/types"
	cosmosevmutils "github.com/cosmos/evm/utils"
	"github.com/cosmos/evm/x/erc20"
	erc20keeper "github.com/cosmos/evm/x/erc20/keeper"
	erc20types "github.com/cosmos/evm/x/erc20/types"
	"github.com/cosmos/evm/x/feemarket"
	feemarketkeeper "github.com/cosmos/evm/x/feemarket/keeper"
	feemarkettypes "github.com/cosmos/evm/x/feemarket/types"
	"github.com/cosmos/evm/x/ibc/transfer"
	transferkeeper "github.com/cosmos/evm/x/ibc/transfer/keeper"
	"github.com/cosmos/evm/x/precisebank"
	precisebankkeeper "github.com/cosmos/evm/x/precisebank/keeper"
	precisebanktypes "github.com/cosmos/evm/x/precisebank/types"
	"github.com/cosmos/evm/x/vm"
	evmkeeper "github.com/cosmos/evm/x/vm/keeper"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/cosmos/gogoproto/proto"
	ratelimit "github.com/cosmos/ibc-apps/modules/rate-limiting/v10"
	ratelimitkeeper "github.com/cosmos/ibc-apps/modules/rate-limiting/v10/keeper"
	ratelimittypes "github.com/cosmos/ibc-apps/modules/rate-limiting/v10/types"
	ica "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts"
	icacontroller "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/controller"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/controller/types"
	icahost "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/types"
	ibctransfer "github.com/cosmos/ibc-go/v10/modules/apps/transfer"
	ibctransfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v10/modules/core"
	ibcclienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v10/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v10/modules/core/keeper"
	ibctm "github.com/cosmos/ibc-go/v10/modules/light-clients/07-tendermint"
	corevm "github.com/ethereum/go-ethereum/core/vm"

	// Force-load the tracer engines to trigger registration due to Go-Ethereum v1.10.15 changes
	_ "github.com/ethereum/go-ethereum/eth/tracers/js"
	_ "github.com/ethereum/go-ethereum/eth/tracers/native"
	"github.com/gorilla/mux"
	marketmap "github.com/skip-mev/connect/v2/x/marketmap"
	marketmapkeeper "github.com/skip-mev/connect/v2/x/marketmap/keeper"
	marketmaptypes "github.com/skip-mev/connect/v2/x/marketmap/types"
	oracle "github.com/skip-mev/connect/v2/x/oracle"
	oraclekeeper "github.com/skip-mev/connect/v2/x/oracle/keeper"
	oracletypes "github.com/skip-mev/connect/v2/x/oracle/types"
	"github.com/spf13/cast"
)

func init() {
	// Replace evmos defaults
	// manually update the power reduction by replacing micro (u) -> atto (a) evmos
	sdk.DefaultPowerReduction = cosmosevmtypes.MicroPowerReduction
	stakingtypes.DefaultMinCommissionRate = math.LegacyZeroDec()

	// DefaultNodeHome default home directories for mantrachaind
	var err error
	DefaultNodeHome, err = clienthelpers.GetNodeHomeDirectory(NodeDir)
	if err != nil {
		panic(err)
	}
}

const appName = "App"

const (
	// ContractMemoryLimit is the memory limit of each contract execution (in MiB)
	// constant value so all nodes run with the same limit.
	ContractMemoryLimit = uint32(32)
)

// We pull these out so we can set them with LDFLAGS in the Makefile
var (
	NodeDir = ".mantrachain"
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string
)

// module account permissions
var maccPerms = map[string][]string{
	authtypes.FeeCollectorName:     {authtypes.Burner},
	distrtypes.ModuleName:          nil,
	icatypes.ModuleName:            nil,
	minttypes.ModuleName:           {authtypes.Minter},
	stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
	stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
	govtypes.ModuleName:            {authtypes.Burner},
	nft.ModuleName:                 nil,
	// non sdk modules
	ibctransfertypes.ModuleName:  {authtypes.Minter, authtypes.Burner},
	ratelimittypes.ModuleName:    nil,
	wasmtypes.ModuleName:         {authtypes.Burner},
	tokenfactorytypes.ModuleName: {authtypes.Minter, authtypes.Burner},
	taxtypes.ModuleName:          nil,
	sanctiontypes.ModuleName:     nil,

	// Cosmos EVM modules
	evmtypes.ModuleName:         {authtypes.Minter, authtypes.Burner},
	feemarkettypes.ModuleName:   nil,
	erc20types.ModuleName:       {authtypes.Minter, authtypes.Burner},
	precisebanktypes.ModuleName: {authtypes.Minter, authtypes.Burner},

	oracletypes.ModuleName: nil,
}

var Upgrades = []upgrades.Upgrade{v5.Upgrade}

var (
	_ runtime.AppI            = (*App)(nil)
	_ servertypes.Application = (*App)(nil)
)

// App extended ABCI application
type App struct {
	*baseapp.BaseApp
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry types.InterfaceRegistry

	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	// keepers
	AccountKeeper         authkeeper.AccountKeeper
	BankKeeper            bankkeeper.BaseKeeper
	StakingKeeper         stakingkeeper.Keeper
	SlashingKeeper        slashingkeeper.Keeper
	MintKeeper            mintkeeper.Keeper
	DistrKeeper           distrkeeper.Keeper
	GovKeeper             govkeeper.Keeper
	CrisisKeeper          *crisiskeeper.Keeper //nolint:staticcheck
	UpgradeKeeper         *upgradekeeper.Keeper
	ParamsKeeper          paramskeeper.Keeper //nolint:staticcheck
	AuthzKeeper           authzkeeper.Keeper
	EvidenceKeeper        evidencekeeper.Keeper
	FeeGrantKeeper        feegrantkeeper.Keeper
	GroupKeeper           groupkeeper.Keeper
	NFTKeeper             nftkeeper.Keeper
	ConsensusParamsKeeper consensusparamkeeper.Keeper
	CircuitKeeper         circuitkeeper.Keeper

	// Connect
	OracleKeeper    *oraclekeeper.Keeper
	MarketMapKeeper *marketmapkeeper.Keeper

	// IBC
	IBCKeeper           *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	ICAHostKeeper       icahostkeeper.Keeper
	ICAControllerKeeper icacontrollerkeeper.Keeper
	TransferKeeper      transferkeeper.Keeper
	WasmKeeper          wasmkeeper.Keeper
	RateLimitKeeper     ratelimitkeeper.Keeper

	// MANTRAChain keepers
	TokenFactoryKeeper tokenfactorykeeper.Keeper
	TaxKeeper          taxkeeper.Keeper
	SanctionKeeper     sanctionkeeper.Keeper

	// Cosmos EVM keepers
	FeeMarketKeeper   feemarketkeeper.Keeper
	EVMKeeper         *evmkeeper.Keeper
	Erc20Keeper       erc20keeper.Keeper
	PreciseBankKeeper precisebankkeeper.Keeper

	// the module manager
	ModuleManager      *module.Manager
	BasicModuleManager module.BasicManager

	// simulation manager
	sm *module.SimulationManager

	// module configurator
	configurator module.Configurator
}

// overrideWasmVariables overrides the wasm variables to:
//   - allow for larger wasm files
func overrideWasmVariables() {
	// Override Wasm size limitation from WASMD.
	wasmtypes.MaxWasmSize = 3 * 1024 * 1024
	wasmtypes.MaxProposalWasmSize = wasmtypes.MaxWasmSize
}

// New returns a reference to an initialized App.
func New(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	appOpts servertypes.AppOptions,
	wasmOpts []wasmkeeper.Option,
	evmChainID uint64,
	evmAppOptions EVMOptionsFn,
	baseAppOptions ...func(*baseapp.BaseApp),
) *App {
	overrideWasmVariables()

	encodingConfig := evmosencoding.MakeConfig(evmChainID)
	appCodec := encodingConfig.Codec
	legacyAmino := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := baseapp.NewBaseApp(appName, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)
	bApp.SetTxEncoder(encodingConfig.TxConfig.TxEncoder())

	// initialize the Cosmos EVM application configuration
	if err := evmAppOptions(evmChainID); err != nil {
		panic(err)
	}

	keys := storetypes.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey, crisistypes.StoreKey,
		minttypes.StoreKey, distrtypes.StoreKey, slashingtypes.StoreKey,
		govtypes.StoreKey, paramstypes.StoreKey, consensusparamtypes.StoreKey, upgradetypes.StoreKey, feegrant.StoreKey,
		evidencetypes.StoreKey,
		circuittypes.StoreKey,
		authzkeeper.StoreKey,
		nftkeeper.StoreKey,
		group.StoreKey,
		// non sdk store keys
		ibcexported.StoreKey, ibctransfertypes.StoreKey,
		wasmtypes.StoreKey,
		ratelimittypes.StoreKey,
		tokenfactorytypes.StoreKey, taxtypes.StoreKey, sanctiontypes.StoreKey,
		icacontrollertypes.StoreKey, icahosttypes.StoreKey,
		oracletypes.StoreKey, marketmaptypes.StoreKey,

		// Cosmos EVM store keys
		evmtypes.StoreKey, feemarkettypes.StoreKey, erc20types.StoreKey, precisebanktypes.StoreKey,
	)

	tkeys := storetypes.NewTransientStoreKeys(paramstypes.TStoreKey, evmtypes.TransientKey, feemarkettypes.TransientKey)
	memKeys := storetypes.NewMemoryStoreKeys()

	// register streaming services
	if err := bApp.RegisterStreamingServices(appOpts, keys); err != nil {
		panic(err)
	}

	app := &App{
		BaseApp:           bApp,
		legacyAmino:       legacyAmino,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	app.ParamsKeeper = initParamsKeeper(
		appCodec,
		legacyAmino,
		keys[paramstypes.StoreKey],
		tkeys[paramstypes.TStoreKey],
	)

	// set the BaseApp's parameter store
	app.ConsensusParamsKeeper = consensusparamkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[consensusparamtypes.StoreKey]),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		runtime.EventService{},
	)
	bApp.SetParamStore(&app.ConsensusParamsKeeper.ParamsStore)

	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		maccPerms,
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[banktypes.StoreKey]),
		app.AccountKeeper,
		BlockedAddresses(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		logger,
	)

	app.StakingKeeper = *stakingkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[stakingtypes.StoreKey]),
		app.AccountKeeper,
		&app.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	)

	app.MintKeeper = mintkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[minttypes.StoreKey]),
		app.StakingKeeper,
		app.AccountKeeper,
		&app.BankKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.DistrKeeper = distrkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[distrtypes.StoreKey]),
		app.AccountKeeper,
		&app.BankKeeper,
		app.StakingKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		legacyAmino,
		runtime.NewKVStoreService(keys[slashingtypes.StoreKey]),
		&app.StakingKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	invCheckPeriod := cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod))
	app.CrisisKeeper = crisiskeeper.NewKeeper( //nolint:staticcheck
		appCodec,
		runtime.NewKVStoreService(keys[crisistypes.StoreKey]),
		invCheckPeriod,
		&app.BankKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		app.AccountKeeper.AddressCodec(),
	)

	app.FeeGrantKeeper = feegrantkeeper.NewKeeper(appCodec, runtime.NewKVStoreService(keys[feegrant.StoreKey]), app.AccountKeeper)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.StakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	)

	app.CircuitKeeper = circuitkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[circuittypes.StoreKey]),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		app.AccountKeeper.AddressCodec(),
	)
	app.SetCircuitBreaker(&app.CircuitKeeper)

	app.AuthzKeeper = authzkeeper.NewKeeper(
		runtime.NewKVStoreService(keys[authzkeeper.StoreKey]),
		appCodec,
		app.MsgServiceRouter(),
		app.AccountKeeper,
	)

	groupConfig := group.DefaultConfig()
	/*
		Example of setting group params:
		groupConfig.MaxMetadataLen = 1000
	*/
	app.GroupKeeper = groupkeeper.NewKeeper(
		keys[group.StoreKey],
		// runtime.NewKVStoreService(keys[group.StoreKey]),
		appCodec,
		app.MsgServiceRouter(),
		app.AccountKeeper,
		groupConfig,
	)

	// get skipUpgradeHeights from the app options
	skipUpgradeHeights := map[int64]bool{}
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}
	homePath := cast.ToString(appOpts.Get(flags.FlagHome))
	// set the governance module account as the authority for conducting upgrades
	app.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		runtime.NewKVStoreService(keys[upgradetypes.StoreKey]),
		appCodec,
		homePath,
		app.BaseApp,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[ibcexported.StoreKey]),
		app.GetSubspace(ibcexported.ModuleName),
		app.UpgradeKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	sortedKnownModules := make([]string, 0, len(maccPerms))
	for moduleName := range maccPerms {
		sortedKnownModules = append(sortedKnownModules, moduleName)
	}
	sort.Strings(sortedKnownModules)

	app.TokenFactoryKeeper = tokenfactorykeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[tokenfactorytypes.StoreKey]),
		sortedKnownModules,
		app.AccountKeeper,
		&app.BankKeeper,
		&app.WasmKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.SanctionKeeper = sanctionkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[sanctiontypes.StoreKey]),
		logger,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.BankKeeper.BaseSendKeeper = app.BankKeeper.SetHooks(
		banktypes.NewMultiBankHooks(
			app.TokenFactoryKeeper.Hooks(),
		))

	// Register the proposal types
	// Deprecated: Avoid adding new handlers, instead use the new proposal flow
	// by granting the governance module the right to execute the message.
	// See: https://docs.cosmos.network/main/modules/gov#proposal-messages
	govRouter := govv1beta1.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler)
	govConfig := govtypes.DefaultConfig()
	/*
		Example of setting gov params:
		govConfig.MaxMetadataLen = 10000
	*/
	govKeeper := govkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[govtypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		app.DistrKeeper,
		app.MsgServiceRouter(),
		govConfig,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Set legacy router for backwards compatibility with gov v1beta1
	govKeeper.SetLegacyRouter(govRouter)

	app.GovKeeper = *govKeeper.SetHooks(
		govtypes.NewMultiGovHooks(
		// register the governance hooks
		),
	)

	app.TaxKeeper = taxkeeper.NewKeeper(
		appCodec,
		app.AccountKeeper.AddressCodec(),
		runtime.NewKVStoreService(keys[taxtypes.StoreKey]),
		logger,
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.FeeCollectorName,
	)

	app.MarketMapKeeper = marketmapkeeper.NewKeeper(
		runtime.NewKVStoreService(keys[marketmaptypes.StoreKey]),
		appCodec,
		authtypes.NewModuleAddress(govtypes.ModuleName),
	)
	marketmapModule := marketmap.NewAppModule(appCodec, app.MarketMapKeeper)

	oracleKeeper := oraclekeeper.NewKeeper(runtime.NewKVStoreService(keys[oracletypes.StoreKey]),
		appCodec,
		app.MarketMapKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName))
	app.OracleKeeper = &oracleKeeper
	oracleModule := oracle.NewAppModule(appCodec, *app.OracleKeeper)

	// ICA Host keeper
	app.ICAHostKeeper = icahostkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(app.keys[icahosttypes.StoreKey]),
		app.GetSubspace(icahosttypes.SubModuleName),
		app.IBCKeeper.ChannelKeeper, // ICS4Wrapper
		app.IBCKeeper.ChannelKeeper,
		app.AccountKeeper,
		bApp.MsgServiceRouter(),
		bApp.GRPCQueryRouter(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Create RateLimit keeper
	clientKeeper := app.IBCKeeper.ClientKeeper
	app.RateLimitKeeper = *ratelimitkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[ratelimittypes.StoreKey]),
		app.GetSubspace(ratelimittypes.ModuleName),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		app.BankKeeper,
		app.IBCKeeper.ChannelKeeper,
		clientKeeper,
		app.IBCKeeper.ChannelKeeper,
	)

	// ICA Controller keeper
	app.ICAControllerKeeper = icacontrollerkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(app.keys[icacontrollertypes.StoreKey]),
		app.GetSubspace(icacontrollertypes.SubModuleName),
		app.IBCKeeper.ChannelKeeper, // ICS4Wrapper
		app.IBCKeeper.ChannelKeeper,
		bApp.MsgServiceRouter(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.NFTKeeper = nftkeeper.NewKeeper(
		runtime.NewKVStoreService(keys[nftkeeper.StoreKey]),
		appCodec,
		app.AccountKeeper,
		app.BankKeeper,
	)

	// create evidence keeper with router
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[evidencetypes.StoreKey]),
		&app.StakingKeeper,
		app.SlashingKeeper,
		app.AccountKeeper.AddressCodec(),
		runtime.ProvideCometInfoService(),
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	app.EvidenceKeeper = *evidenceKeeper

	// Cosmos EVM keepers
	app.FeeMarketKeeper = feemarketkeeper.NewKeeper(
		appCodec,
		authtypes.NewModuleAddress(govtypes.ModuleName),
		keys[feemarkettypes.StoreKey],
		tkeys[feemarkettypes.TransientKey],
	)

	// Set up PreciseBank keeper
	//
	// NOTE: PreciseBank is not needed if SDK use 18 decimals for gas coin. Use BankKeeper instead.
	app.PreciseBankKeeper = precisebankkeeper.NewKeeper(
		appCodec,
		keys[precisebanktypes.StoreKey],
		app.BankKeeper,
		app.AccountKeeper,
	)

	// Set up EVM keeper
	tracer := cast.ToString(appOpts.Get(srvflags.EVMTracer))

	app.EVMKeeper = evmkeeper.NewKeeper(
		appCodec, keys[evmtypes.StoreKey], tkeys[evmtypes.TransientKey],
		authtypes.NewModuleAddress(govtypes.ModuleName),
		app.AccountKeeper,
		app.PreciseBankKeeper,
		app.StakingKeeper,
		app.FeeMarketKeeper,
		&app.Erc20Keeper,
		tracer,
	)

	// ERC20 Keeper
	app.Erc20Keeper = erc20keeper.NewKeeper(
		keys[erc20types.StoreKey],
		appCodec,
		authtypes.NewModuleAddress(govtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.EVMKeeper,
		app.StakingKeeper,
		&app.TransferKeeper,
	)

	// instantiate IBC transfer keeper AFTER the ERC-20 keeper to use it in the instantiation
	app.TransferKeeper = transferkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[ibctransfertypes.StoreKey]),
		app.GetSubspace(ibctransfertypes.ModuleName),
		app.RateLimitKeeper, // ISC4 Wrapper: RateLimit IBC middleware
		app.IBCKeeper.ChannelKeeper,
		bApp.MsgServiceRouter(),
		app.AccountKeeper,
		app.BankKeeper,
		app.Erc20Keeper, // Add ERC20 Keeper for ERC20 transfers
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Create IBC modules with ratelimit middleware
	// channel.RecvPacket -> ratelimit.OnRecvPacket -> transfer.OnRecvPacket
	var transferStack porttypes.IBCModule
	transferStack = transfer.NewIBCModule(app.TransferKeeper)
	transferStack = ratelimit.NewIBCMiddleware(app.RateLimitKeeper, transferStack)
	// register escrow address for tokenfactory when channel opens
	transferStack = tokenfactory.NewIBCModule(transferStack, app.TokenFactoryKeeper)

	// Create ICAHost Stack
	var icaHostStack porttypes.IBCModule = icahost.NewIBCModule(app.ICAHostKeeper)

	// Create Interchain Accounts Controller Stack
	var icaControllerStack porttypes.IBCModule = icacontroller.NewIBCMiddleware(app.ICAControllerKeeper)

	// Create fee enabled wasm ibc Stack
	wasmStack := wasm.NewIBCHandler(&app.WasmKeeper, app.IBCKeeper.ChannelKeeper, app.IBCKeeper.ChannelKeeper)

	// Create static IBC router, add app routes, then set and seal it
	ibcRouter := porttypes.NewRouter().
		AddRoute(icahosttypes.SubModuleName, icaHostStack).
		AddRoute(icacontrollertypes.SubModuleName, icaControllerStack).
		AddRoute(ibctransfertypes.ModuleName, transferStack).
		AddRoute(wasmtypes.ModuleName, wasmStack)
	app.IBCKeeper.SetRouter(ibcRouter)

	// TODO: add IBC v2 (Eureka) support in future
	// var transferStackv2 ibcapi.IBCModule
	// transferStackv2 = transferv2.NewIBCModule(app.TransferKeeper)
	// transferStackv2 = ratelimitv2.NewIBCMiddleware(app.RateLimitKeeper, transferStackv2)
	// transferStackv2 = tokenfactory.NewIBCV2Module(transferStackv2, app.TokenFactoryKeeper)

	// ibcRouterV2 := ibcapi.NewRouter()
	// ibcRouterV2.AddRoute(ibctransfertypes.ModuleName, transferStackv2)
	// app.IBCKeeper.SetRouterV2(ibcRouterV2)

	// TODO: Configure EVM precompiles when needed
	// corePrecompiles := NewAvailableStaticPrecompiles(
	// 	app.StakingKeeper,
	// 	app.DistrKeeper,
	// 	app.PreciseBankKeeper,
	// 	app.Erc20Keeper,
	// 	app.TransferKeeper,
	// 	app.IBCKeeper.ChannelKeeper,
	// 	app.EVMKeeper,
	// 	app.GovKeeper,
	// 	app.SlashingKeeper,
	// 	app.EvidenceKeeper,
	// 	appCodec,
	// )
	// app.EVMKeeper.WithStaticPrecompiles(
	// 	corePrecompiles,
	// )

	storeProvider := app.IBCKeeper.ClientKeeper.GetStoreProvider()
	tmLightClientModule := ibctm.NewLightClientModule(appCodec, storeProvider)
	clientKeeper.AddRoute(ibctm.ModuleName, &tmLightClientModule)
	// wasmLightClientModule := ibcwasm.NewLightClientModule(app.WasmKeeper, storeProvider)
	// clientKeeper.AddRoute(ibcwasmtypes.ModuleName, &wasmLightClientModule)

	wasmConfig, err := wasm.ReadNodeConfig(appOpts)
	if err != nil {
		panic(fmt.Sprintf("error while reading wasm config: %s", err))
	}

	wasmDir := filepath.Join(homePath, "wasm")

	// Register custom plugins for the wasm module by appending them to the existing options
	wasmOpts = append(wasmOpts, queries.RegisterCustomPlugins(
		*app.GRPCQueryRouter(),
		app.AppCodec(),
	)...)

	// The last arguments can contain custom message handlers, and custom query handlers,
	// if we want to allow any custom callbacks
	app.WasmKeeper = wasmkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[wasmtypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		distrkeeper.NewQuerier(app.DistrKeeper),
		app.IBCKeeper.ChannelKeeper, // ISC4 Wrapper
		app.IBCKeeper.ChannelKeeper,
		app.TransferKeeper,
		app.MsgServiceRouter(),
		app.GRPCQueryRouter(),
		wasmDir,
		wasmConfig,
		wasmtypes.VMConfig{},
		AllCapabilities(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		nil, // TODO: add ibc v2 wasm keeper when available
		wasmOpts...,
	)

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants)) //nolint:staticcheck

	// optional: enable sign mode textual by overwriting the default tx config (after setting the bank keeper)
	txConfigOpts := authtx.ConfigOptions{
		EnabledSignModes:           append(authtx.DefaultSignModes, sigtypes.SignMode_SIGN_MODE_TEXTUAL),
		TextualCoinMetadataQueryFn: txmodule.NewBankKeeperCoinMetadataQueryFn(app.BankKeeper),
	}
	txConfig, err := authtx.NewTxConfigWithOptions(
		appCodec,
		txConfigOpts,
	)
	if err != nil {
		panic(err)
	}
	app.txConfig = txConfig

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.ModuleManager = module.NewManager(
		genutil.NewAppModule(
			app.AccountKeeper,
			app.StakingKeeper,
			app,
			txConfig,
		),
		auth.NewAppModule(appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts, nil),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, nil),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		gov.NewAppModule(appCodec, &app.GovKeeper, app.AccountKeeper, app.BankKeeper, nil),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper, nil, nil),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, nil, app.interfaceRegistry),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, nil),
		staking.NewAppModule(appCodec, &app.StakingKeeper, app.AccountKeeper, app.BankKeeper, nil),
		upgrade.NewAppModule(app.UpgradeKeeper, app.AccountKeeper.AddressCodec()),
		evidence.NewAppModule(app.EvidenceKeeper),
		params.NewAppModule(app.ParamsKeeper), //nolint:staticcheck
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		groupmodule.NewAppModule(appCodec, app.GroupKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		nftmodule.NewAppModule(appCodec, app.NFTKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		consensus.NewAppModule(appCodec, app.ConsensusParamsKeeper),
		circuit.NewAppModule(appCodec, app.CircuitKeeper),
		// non sdk modules
		wasm.NewAppModule(appCodec, &app.WasmKeeper, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.MsgServiceRouter(), nil),
		ibc.NewAppModule(app.IBCKeeper),
		transfer.NewAppModule(app.TransferKeeper),
		ica.NewAppModule(&app.ICAControllerKeeper, &app.ICAHostKeeper),
		ibctm.NewAppModule(tmLightClientModule),
		ratelimit.NewAppModule(appCodec, app.RateLimitKeeper),
		// connect
		marketmapModule,
		oracleModule,

		// sdk
		crisis.NewAppModule(app.CrisisKeeper, skipGenesisInvariants, nil), //nolint:staticcheck

		// mantrachain modules
		tokenfactory.NewAppModule(appCodec, app.TokenFactoryKeeper),
		tax.NewAppModule(appCodec, app.TaxKeeper),
		sanction.NewAppModule(appCodec, app.SanctionKeeper),

		// Cosmos EVM modules
		vm.NewAppModule(app.EVMKeeper, app.AccountKeeper, app.AccountKeeper.AddressCodec()),
		feemarket.NewAppModule(app.FeeMarketKeeper),
		erc20.NewAppModule(app.Erc20Keeper, app.AccountKeeper),
		precisebank.NewAppModule(app.PreciseBankKeeper, app.BankKeeper, app.AccountKeeper),
	)

	// BasicModuleManager defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration and genesis verification.
	// By default it is composed of all the module from the module manager.
	// Additionally, app module basics can be overwritten by passing them as argument.
	app.BasicModuleManager = module.NewBasicManagerFromManager(
		app.ModuleManager,
		map[string]module.AppModuleBasic{
			genutiltypes.ModuleName: genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
			govtypes.ModuleName: gov.NewAppModuleBasic(
				[]govclient.ProposalHandler{
					paramsclient.ProposalHandler,
				},
			),
			ibctransfertypes.ModuleName: ibctransfer.AppModuleBasic{},
		})
	app.BasicModuleManager.RegisterLegacyAminoCodec(legacyAmino)
	app.BasicModuleManager.RegisterInterfaces(interfaceRegistry)

	// NOTE: upgrade module is required to be prioritized
	app.ModuleManager.SetOrderPreBlockers(
		upgradetypes.ModuleName,
		authtypes.ModuleName,
	)
	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	app.ModuleManager.SetOrderBeginBlockers(
		minttypes.ModuleName,
		// mca tax before distribution
		taxtypes.ModuleName,

		// IBC modules
		ibcexported.ModuleName, ibctransfertypes.ModuleName,

		// Cosmos EVM BeginBlockers
		erc20types.ModuleName, feemarkettypes.ModuleName,
		evmtypes.ModuleName, // NOTE: EVM BeginBlocker must come after FeeMarket BeginBlocker

		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		genutiltypes.ModuleName,
		authz.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		// additional non simd modules
		icatypes.ModuleName,
		ratelimittypes.ModuleName,
		wasmtypes.ModuleName,
		tokenfactorytypes.ModuleName,
		oracletypes.ModuleName,
		marketmaptypes.ModuleName,
		sanctiontypes.ModuleName,
		precisebanktypes.ModuleName,
	)

	app.ModuleManager.SetOrderEndBlockers(
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		genutiltypes.ModuleName,
		// Cosmos EVM EndBlockers
		evmtypes.ModuleName,
		erc20types.ModuleName,
		feemarkettypes.ModuleName,
		feegrant.ModuleName,
		group.ModuleName,
		// burn fees from fee collector
		taxtypes.ModuleName,
		// additional non simd modules
		ibctransfertypes.ModuleName,
		ibcexported.ModuleName,
		icatypes.ModuleName,
		ratelimittypes.ModuleName,
		wasmtypes.ModuleName,
		tokenfactorytypes.ModuleName,
		oracletypes.ModuleName,
		marketmaptypes.ModuleName,
		sanctiontypes.ModuleName,
		precisebanktypes.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: The genutils module must also occur after auth so that it can access the params from auth.
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	// NOTE: wasm module should be at the end as it can call other module functionality direct or via message dispatching during
	// genesis phase. For example bank transfer, auth account check, staking, ...
	genesisModuleOrder := []string{
		// simd modules
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		crisistypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		nft.ModuleName,
		group.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		consensusparamtypes.ModuleName,
		circuittypes.ModuleName,
		ibcexported.ModuleName,

		// Cosmos EVM modules
		//
		// NOTE: feemarket module needs to be initialized before genutil module:
		// gentx transactions use MinGasPriceDecorator.AnteHandle
		evmtypes.ModuleName,
		feemarkettypes.ModuleName,
		erc20types.ModuleName,
		precisebanktypes.ModuleName,

		// additional non simd modules
		ibctransfertypes.ModuleName,
		genutiltypes.ModuleName,
		icatypes.ModuleName,
		ratelimittypes.ModuleName,
		// wasm after ibc transfer

		wasmtypes.ModuleName,
		tokenfactorytypes.ModuleName,
		taxtypes.ModuleName,
		sanctiontypes.ModuleName,

		// market map genesis must be called AFTER all consuming modules (i.e. x/oracle, etc.)
		oracletypes.ModuleName,
		marketmaptypes.ModuleName,
	}
	app.ModuleManager.SetOrderInitGenesis(genesisModuleOrder...)
	app.ModuleManager.SetOrderExportGenesis(genesisModuleOrder...)

	// Uncomment if you want to set a custom migration order here.
	// app.ModuleManager.SetOrderMigrations(custom order)

	app.ModuleManager.RegisterInvariants(app.CrisisKeeper) //nolint:staticcheck
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	err = app.ModuleManager.RegisterServices(app.configurator)
	if err != nil {
		panic(err)
	}

	autocliv1.RegisterQueryServer(app.GRPCQueryRouter(), runtimeservices.NewAutoCLIQueryService(app.ModuleManager.Modules))

	reflectionSvc, err := runtimeservices.NewReflectionService()
	if err != nil {
		panic(err)
	}
	reflectionv1.RegisterReflectionServiceServer(app.GRPCQueryRouter(), reflectionSvc)

	// add test gRPC service for testing gRPC queries in isolation
	// testdata_pulsar.RegisterQueryServer(app.GRPCQueryRouter(), testdata_pulsar.QueryImpl{})

	// create the simulation manager and define the order of the modules for deterministic simulations
	//
	// NOTE: this is not required apps that don't use the simulator for fuzz testing
	// transactions
	overrideModules := map[string]module.AppModuleSimulation{
		authtypes.ModuleName: auth.NewAppModule(app.appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts, nil),
	}
	app.sm = module.NewSimulationManagerFromAppModules(app.ModuleManager.Modules, overrideModules)

	app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	maxGasWanted := cast.ToUint64(appOpts.Get(srvflags.EVMMaxTxGasWanted))

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetPreBlocker(app.PreBlocker)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	// set denom resolver to test variant.
	// app.FeeMarketKeeper.SetDenomResolver(&ante.DenomResolverImpl{
	//	StakingKeeper: &app.StakingKeeper,
	// })

	// must be before Loading version
	// requires the snapshot store to be created and registered as a BaseAppOption
	// see cmd/mantrachaind/root.go: 206 - 214 approx
	if manager := app.SnapshotManager(); manager != nil {
		err := manager.RegisterExtensions(
			wasmkeeper.NewWasmSnapshotter(app.CommitMultiStore(), &app.WasmKeeper),
		)
		if err != nil {
			panic(fmt.Errorf("failed to register snapshot extension: %s", err))
		}
	}

	app.setAnteHandler(txConfig, wasmConfig, keys[wasmtypes.StoreKey], maxGasWanted)
	// app.setPostHandler()

	// oracle initialization
	client, metrics, err := app.initializeOracle(appOpts)
	if err != nil {
		panic(fmt.Errorf("failed to initialize oracle client and metrics: %w", err))
	}

	app.MarketMapKeeper.SetHooks(app.OracleKeeper.Hooks())

	app.initializeABCIExtensions(client, metrics)

	// Register any on-chain upgrades.
	app.setupUpgradeStoreLoaders()
	app.setupUpgradeHandlers()

	// At startup, after all modules have been registered, check that all proto
	// annotations are correct.
	protoFiles, err := proto.MergedRegistry()
	if err != nil {
		panic(err)
	}
	err = msgservice.ValidateProtoAnnotations(protoFiles)
	if err != nil {
		// Once we switch to using protoreflect-based antehandlers, we might
		// want to panic here instead of logging a warning.
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	}

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			panic(fmt.Errorf("error loading last version: %w", err))
		}
		ctx := app.NewUncachedContext(true, tmproto.Header{})

		// Initialize pinned codes in wasmvm as they are not persisted there
		if err := app.WasmKeeper.InitializePinnedCodes(ctx); err != nil {
			panic(fmt.Sprintf("failed initialize pinned codes %s", err))
		}
	}

	return app
}

func (app *App) setAnteHandler(txConfig client.TxConfig, wasmConfig wasmtypes.NodeConfig, txCounterStoreKey *storetypes.KVStoreKey, maxGasWanted uint64) {
	evmHandlerOpts := NewEVMAnteHandlerOptionsFromApp(app, txConfig, maxGasWanted)

	if err := evmHandlerOpts.Validate(); err != nil {
		panic(err)
	}

	handlerOpts := ante.HandlerOptions{
		EvmOptions:            evmHandlerOpts.Options(),
		IBCKeeper:             app.IBCKeeper,
		WasmConfig:            &wasmConfig,
		WasmKeeper:            &app.WasmKeeper,
		TXCounterStoreService: runtime.NewKVStoreService(txCounterStoreKey),
		CircuitKeeper:         &app.CircuitKeeper,
		SanctionKeeper:        &app.SanctionKeeper,
	}

	if err := handlerOpts.Validate(); err != nil {
		panic(err)
	}

	// Set the AnteHandler for the app
	app.SetAnteHandler(ante.NewAnteHandler(handlerOpts))
}

// TODO: Implement post handler
// func (app *App) setPostHandler() {
// }

// Name returns the name of the App
func (app *App) Name() string { return app.BaseApp.Name() }

// PreBlocker application updates every pre block
func (app *App) PreBlocker(ctx sdk.Context, _ *abci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
	return app.ModuleManager.PreBlock(ctx)
}

// BeginBlocker application updates every begin block
func (app *App) BeginBlocker(ctx sdk.Context) (sdk.BeginBlock, error) {
	return app.ModuleManager.BeginBlock(ctx)
}

// EndBlocker application updates every end block
func (app *App) EndBlocker(ctx sdk.Context) (sdk.EndBlock, error) {
	return app.ModuleManager.EndBlock(ctx)
}

func (app *App) Configurator() module.Configurator {
	return app.configurator
}

// InitChainer application update at chain initialization
func (app *App) InitChainer(ctx sdk.Context, req *abci.RequestInitChain) (*abci.ResponseInitChain, error) {
	var genesisState GenesisState
	if err := json.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}
	err := app.UpgradeKeeper.SetModuleVersionMap(ctx, app.ModuleManager.GetVersionMap())
	if err != nil {
		panic(err)
	}
	response, err := app.ModuleManager.InitGenesis(ctx, app.appCodec, genesisState)
	return response, err
}

// LoadHeight loads a particular height
func (app *App) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// LegacyAmino returns legacy amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *App) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

// AppCodec returns app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *App) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns App's InterfaceRegistry
func (app *App) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// TxConfig returns App's TxConfig
func (app *App) TxConfig() client.TxConfig {
	return app.txConfig
}

// AutoCliOpts returns the autocli options for the app.
func (app *App) AutoCliOpts() autocli.AppOptions {
	modules := make(map[string]appmodule.AppModule, 0)
	for _, m := range app.ModuleManager.Modules {
		if moduleWithName, ok := m.(module.HasName); ok {
			moduleName := moduleWithName.Name()
			if appModule, ok := moduleWithName.(appmodule.AppModule); ok {
				modules[moduleName] = appModule
			}
		}
	}

	return autocli.AppOptions{
		Modules:               modules,
		ModuleOptions:         runtimeservices.ExtractAutoCLIOptions(app.ModuleManager.Modules),
		AddressCodec:          authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		ValidatorAddressCodec: authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		ConsensusAddressCodec: authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	}
}

// DefaultGenesis returns a default genesis from the registered AppModuleBasic's.
func (app *App) DefaultGenesis() map[string]json.RawMessage {
	genesis := app.BasicModuleManager.DefaultGenesis(app.appCodec)

	// Add mint denom configuration
	mintGenState := minttypes.DefaultGenesisState()
	mintGenState.Params.MintDenom = sdk.DefaultBondDenom
	genesis[minttypes.ModuleName] = app.appCodec.MustMarshalJSON(mintGenState)

	// Add EVM genesis configuration
	evmGenState := evmtypes.DefaultGenesisState()
	evmGenState.Params.ActiveStaticPrecompiles = evmtypes.AvailableStaticPrecompiles
	genesis[evmtypes.ModuleName] = app.appCodec.MustMarshalJSON(evmGenState)

	// Add ERC20 genesis configuration
	erc20GenState := erc20types.DefaultGenesisState()
	erc20GenState.TokenPairs = ExampleTokenPairs
	erc20GenState.Params.NativePrecompiles = append(erc20GenState.Params.NativePrecompiles, WTokenContractMainnet)
	genesis[erc20types.ModuleName] = app.appCodec.MustMarshalJSON(erc20GenState)

	return genesis
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetKey(storeKey string) *storetypes.KVStoreKey {
	return app.keys[storeKey]
}

// GetStoreKeys returns all the stored store keys.
func (app *App) GetStoreKeys() []storetypes.StoreKey {
	keys := make([]storetypes.StoreKey, 0, len(app.keys))
	for _, key := range app.keys {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Name() < keys[j].Name()
	})
	return keys
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *App) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// SimulationManager implements the SimulationApp interface
func (app *App) SimulationManager() *module.SimulationManager {
	return app.sm
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *App) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register new CometBFT queries routes from grpc-gateway.
	cmtservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register node gRPC service for grpc-gateway.
	nodeservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register grpc-gateway routes for all modules.
	app.BasicModuleManager.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register swagger API from root so that other applications can override easily
	if apiConfig.Swagger {
		RegisterSwaggerAPI(clientCtx, apiSvr.Router)
	}
}

// RegisterSwaggerAPI registers swagger route with API Server.
func RegisterSwaggerAPI(ctx client.Context, rtr *mux.Router) {
	staticServer := http.FileServer(swagger.FS)
	rtr.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticServer))
	rtr.PathPrefix("/swagger/").Handler(staticServer)
	rtr.PathPrefix("/openapi/").Handler(staticServer)
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *App) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.GRPCQueryRouter(), clientCtx, app.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *App) RegisterTendermintService(clientCtx client.Context) {
	cmtApp := server.NewCometABCIWrapper(app)
	cmtservice.RegisterTendermintService(
		clientCtx,
		app.GRPCQueryRouter(),
		app.interfaceRegistry,
		cmtApp.Query,
	)
}

func (app *App) RegisterNodeService(clientCtx client.Context, cfg config.Config) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter(), cfg)
}

// configure store loader that checks if version == upgradeHeight and applies store upgrades
func (app *App) setupUpgradeStoreLoaders() {
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

	for _, upgrade := range Upgrades {
		if upgradeInfo.Name == upgrade.UpgradeName {
			storeUpgrades := upgrade.StoreUpgrades
			app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
		}
	}
}

func (app *App) setupUpgradeHandlers() {
	for _, upgrade := range Upgrades {
		app.UpgradeKeeper.SetUpgradeHandler(
			upgrade.UpgradeName,
			upgrade.CreateUpgradeHandler(
				app.ModuleManager,
				app.configurator,
				&upgrades.UpgradeKeepers{
					ChannelKeeper:      app.IBCKeeper.ChannelKeeper,
					TransferKeeper:     app.TransferKeeper,
					TokenFactoryKeeper: &app.TokenFactoryKeeper,
					SanctionKeeper:     app.SanctionKeeper,
					FeeMarketKeeper:    app.FeeMarketKeeper,
					AccountKeeper:      app.AccountKeeper,
					BankKeeper:         app.BankKeeper,
					EVMKeeper:          *app.EVMKeeper,
				},
				app.keys,
			),
		)
	}
}

// GetMaccPerms returns a copy of the module account permissions
//
// NOTE: This is solely to be used for testing purposes.
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}

	return dupMaccPerms
}

// BlockedAddresses returns all the app's blocked account addresses.
func BlockedAddresses() map[string]bool {
	blockedAddrs := make(map[string]bool)

	// Add after existing code:
	blockedPrecompilesHex := evmtypes.AvailableStaticPrecompiles
	for _, addr := range corevm.PrecompiledAddressesBerlin {
		blockedPrecompilesHex = append(blockedPrecompilesHex, addr.Hex())
	}

	for _, precompile := range blockedPrecompilesHex {
		blockedAddrs[cosmosevmutils.Bech32StringFromHexAddress(precompile)] = true
	}

	return blockedAddrs
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper { //nolint:staticcheck
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey) //nolint:staticcheck
	keyTable := ibcclienttypes.ParamKeyTable()
	paramsKeeper.Subspace(ratelimittypes.ModuleName).WithKeyTable(ratelimittypes.ParamKeyTable())
	paramsKeeper.Subspace(icacontrollertypes.SubModuleName).WithKeyTable(icacontrollertypes.ParamKeyTable())
	paramsKeeper.Subspace(icahosttypes.SubModuleName).WithKeyTable(icahosttypes.ParamKeyTable())
	paramsKeeper.Subspace(ibcexported.ModuleName).WithKeyTable(keyTable)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName).WithKeyTable(ibctransfertypes.ParamKeyTable())

	return paramsKeeper
}
