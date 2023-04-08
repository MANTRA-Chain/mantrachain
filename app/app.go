package app

// DONTCOVER

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cast"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	// IBC modules
	ica "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts"
	icahost "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/types"
	"github.com/cosmos/ibc-go/v4/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v4/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v4/modules/core"
	porttypes "github.com/cosmos/ibc-go/v4/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v4/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v4/modules/core/keeper"

	// budget module
	"github.com/tendermint/budget/x/budget"
	budgetkeeper "github.com/tendermint/budget/x/budget/keeper"
	budgettypes "github.com/tendermint/budget/x/budget/types"

	// core modules
	farmingparams "github.com/MANTRA-Finance/mantrachain/app/params"
	"github.com/MANTRA-Finance/mantrachain/x/liquidfarming"
	liquidfarmingkeeper "github.com/MANTRA-Finance/mantrachain/x/liquidfarming/keeper"
	liquidfarmingtypes "github.com/MANTRA-Finance/mantrachain/x/liquidfarming/types"
	"github.com/MANTRA-Finance/mantrachain/x/liquidity"
	liquiditykeeper "github.com/MANTRA-Finance/mantrachain/x/liquidity/keeper"
	liquiditytypes "github.com/MANTRA-Finance/mantrachain/x/liquidity/types"
	"github.com/MANTRA-Finance/mantrachain/x/lpfarm"
	lpfarmkeeper "github.com/MANTRA-Finance/mantrachain/x/lpfarm/keeper"
	lpfarmtypes "github.com/MANTRA-Finance/mantrachain/x/lpfarm/types"
	"github.com/MANTRA-Finance/mantrachain/x/marketmaker"
	marketmakerkeeper "github.com/MANTRA-Finance/mantrachain/x/marketmaker/keeper"
	marketmakertypes "github.com/MANTRA-Finance/mantrachain/x/marketmaker/types"

	"github.com/MANTRA-Finance/mantrachain/x/coinfactory"
	coinfactorykeeper "github.com/MANTRA-Finance/mantrachain/x/coinfactory/keeper"
	coinfactorytypes "github.com/MANTRA-Finance/mantrachain/x/coinfactory/types"
	"github.com/MANTRA-Finance/mantrachain/x/guard"
	guardkeeper "github.com/MANTRA-Finance/mantrachain/x/guard/keeper"
	guardtypes "github.com/MANTRA-Finance/mantrachain/x/guard/types"
	"github.com/MANTRA-Finance/mantrachain/x/nft"
	nftkeeper "github.com/MANTRA-Finance/mantrachain/x/nft/keeper"
	nfttypes "github.com/MANTRA-Finance/mantrachain/x/nft/types"
	"github.com/MANTRA-Finance/mantrachain/x/token"
	tokenkeeper "github.com/MANTRA-Finance/mantrachain/x/token/keeper"
	tokentypes "github.com/MANTRA-Finance/mantrachain/x/token/types"

	ibcconsumer "github.com/cosmos/interchain-security/x/ccv/consumer"
	ibcconsumerkeeper "github.com/cosmos/interchain-security/x/ccv/consumer/keeper"
	ibcconsumertypes "github.com/cosmos/interchain-security/x/ccv/consumer/types"

	// unnamed import of statik for swagger UI support
	_ "github.com/cosmos/cosmos-sdk/client/docs/statik"
)

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string

	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		feegrantmodule.AppModuleBasic{},
		authzmodule.AppModuleBasic{},
		ibc.AppModuleBasic{},
		transfer.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		vesting.AppModuleBasic{},
		budget.AppModuleBasic{},
		liquidity.AppModuleBasic{},
		liquidfarming.AppModuleBasic{},
		marketmaker.AppModuleBasic{},
		lpfarm.AppModuleBasic{},
		ica.AppModuleBasic{},
		nft.AppModuleBasic{},
		token.AppModuleBasic{},
		guard.AppModuleBasic{},
		coinfactory.AppModuleBasic{},
		ibcconsumer.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:                    nil,
		budgettypes.ModuleName:                        nil,
		liquiditytypes.ModuleName:                     {authtypes.Minter, authtypes.Burner},
		liquidfarmingtypes.ModuleName:                 {authtypes.Minter, authtypes.Burner},
		ibctransfertypes.ModuleName:                   {authtypes.Minter, authtypes.Burner},
		marketmakertypes.ModuleName:                   nil,
		lpfarmtypes.ModuleName:                        nil,
		icatypes.ModuleName:                           nil,
		nfttypes.ModuleName:                           nil,
		coinfactorytypes.ModuleName:                   {authtypes.Minter, authtypes.Burner},
		tokentypes.ModuleName:                         nil,
		guardtypes.ModuleName:                         nil,
		ibcconsumertypes.ConsumerRedistributeName:     nil,
		ibcconsumertypes.ConsumerToSendToProviderName: nil,
	}
)

// Verify app interface at compile time
var (
	_ servertypes.Application = (*App)(nil)
	_ servertypes.Application = (*App)(nil)
)

// App extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type App struct {
	*baseapp.BaseApp

	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*sdk.KVStoreKey
	tkeys   map[string]*sdk.TransientStoreKey
	memKeys map[string]*sdk.MemoryStoreKey

	// keepers
	AccountKeeper       authkeeper.AccountKeeper
	BankKeeper          bankkeeper.Keeper
	CapabilityKeeper    *capabilitykeeper.Keeper
	SlashingKeeper      slashingkeeper.Keeper
	CrisisKeeper        crisiskeeper.Keeper
	UpgradeKeeper       upgradekeeper.Keeper
	ParamsKeeper        paramskeeper.Keeper
	IBCKeeper           *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	EvidenceKeeper      evidencekeeper.Keeper
	TransferKeeper      ibctransferkeeper.Keeper
	FeeGrantKeeper      feegrantkeeper.Keeper
	AuthzKeeper         authzkeeper.Keeper
	BudgetKeeper        budgetkeeper.Keeper
	LiquidityKeeper     liquiditykeeper.Keeper
	LiquidFarmingKeeper liquidfarmingkeeper.Keeper
	MarketMakerKeeper   marketmakerkeeper.Keeper
	LPFarmKeeper        lpfarmkeeper.Keeper
	ICAHostKeeper       icahostkeeper.Keeper

	TokenKeeper       tokenkeeper.Keeper
	GuardKeeper       guardkeeper.Keeper
	CoinFactoryKeeper coinfactorykeeper.Keeper
	NFTKeeper         nftkeeper.Keeper
	ConsumerKeeper    ibcconsumerkeeper.Keeper

	// scoped keepers
	ScopedIBCKeeper         capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper    capabilitykeeper.ScopedKeeper
	ScopedICAHostKeeper     capabilitykeeper.ScopedKeeper
	ScopedIBCConsumerKeeper capabilitykeeper.ScopedKeeper

	// IBC app modules
	transferModule transfer.AppModule
	icaModule      ica.AppModule

	// the module manager
	mm *module.Manager

	// module configurator
	configurator module.Configurator
}

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, AppUserHomeDir)
}

// NewApp returns a reference to an initialized App.
func NewApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	encodingConfig farmingparams.EncodingConfig,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *App {
	appCodec := encodingConfig.Marshaler
	legacyAmino := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := baseapp.NewBaseApp(AppName, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey,
		banktypes.StoreKey,
		slashingtypes.StoreKey,
		paramstypes.StoreKey,
		ibchost.StoreKey,
		upgradetypes.StoreKey,
		evidencetypes.StoreKey,
		ibctransfertypes.StoreKey,
		capabilitytypes.StoreKey,
		feegrant.StoreKey,
		authzkeeper.StoreKey,
		budgettypes.StoreKey,
		liquiditytypes.StoreKey,
		liquidfarmingtypes.StoreKey,
		marketmakertypes.StoreKey,
		lpfarmtypes.StoreKey,
		icahosttypes.StoreKey,
		tokentypes.StoreKey,
		guardtypes.StoreKey,
		coinfactorytypes.StoreKey,
		nftkeeper.StoreKey,
		ibcconsumertypes.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	app := &App{
		BaseApp:           bApp,
		legacyAmino:       legacyAmino,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
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
	bApp.SetParamStore(
		app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramskeeper.ConsensusParamsKeyTable()),
	)

	app.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, keys[capabilitytypes.StoreKey], memKeys[capabilitytypes.MemStoreKey])
	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	scopedICAHostKeeper := app.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)
	scopedIBCConsumerKeeper := app.CapabilityKeeper.ScopeToModule(ibcconsumertypes.ModuleName)
	app.CapabilityKeeper.Seal()

	// add keepers
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		keys[authtypes.StoreKey],
		app.GetSubspace(authtypes.ModuleName),
		authtypes.ProtoBaseAccount,
		maccPerms,
	)
	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		keys[banktypes.StoreKey],
		app.AccountKeeper,
		app.GetSubspace(banktypes.ModuleName),
		app.BlockedModuleAccountAddrs(),
		&app.GuardKeeper,
	)
	app.AuthzKeeper = authzkeeper.NewKeeper(
		keys[authzkeeper.StoreKey],
		appCodec,
		app.BaseApp.MsgServiceRouter(),
	)
	app.FeeGrantKeeper = feegrantkeeper.NewKeeper(
		appCodec,
		keys[feegrant.StoreKey],
		app.AccountKeeper,
	)
	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		keys[slashingtypes.StoreKey],
		&app.ConsumerKeeper,
		app.GetSubspace(slashingtypes.ModuleName),
	)
	app.CrisisKeeper = crisiskeeper.NewKeeper(
		app.GetSubspace(crisistypes.ModuleName),
		invCheckPeriod,
		app.BankKeeper,
		authtypes.FeeCollectorName,
	)
	app.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		keys[upgradetypes.StoreKey],
		appCodec,
		homePath,
		app.BaseApp,
	)

	app.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		keys[ibchost.StoreKey],
		app.GetSubspace(ibchost.ModuleName),
		&app.ConsumerKeeper,
		app.UpgradeKeeper,
		scopedIBCKeeper,
	)

	// Create CCV consumer and modules
	app.ConsumerKeeper = ibcconsumerkeeper.NewKeeper(
		appCodec,
		keys[ibcconsumertypes.StoreKey],
		app.GetSubspace(ibcconsumertypes.ModuleName),
		scopedIBCConsumerKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.IBCKeeper.ConnectionKeeper,
		app.IBCKeeper.ClientKeeper,
		app.SlashingKeeper,
		app.BankKeeper,
		app.AccountKeeper,
		&app.TransferKeeper,
		app.IBCKeeper,
		authtypes.FeeCollectorName,
	)

	// register slashing module Slashing hooks to the consumer keeper
	app.ConsumerKeeper = *app.ConsumerKeeper.SetHooks(app.SlashingKeeper.Hooks())
	consumerModule := ibcconsumer.NewAppModule(app.ConsumerKeeper)

	app.BudgetKeeper = budgetkeeper.NewKeeper(
		appCodec,
		keys[budgettypes.StoreKey],
		app.GetSubspace(budgettypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.ModuleAccountAddrs(),
	)
	app.LiquidityKeeper = liquiditykeeper.NewKeeper(
		appCodec,
		keys[liquiditytypes.StoreKey],
		app.GetSubspace(liquiditytypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		&app.GuardKeeper,
	)
	app.MarketMakerKeeper = marketmakerkeeper.NewKeeper(
		appCodec,
		keys[marketmakertypes.StoreKey],
		app.GetSubspace(marketmakertypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
	)
	app.LPFarmKeeper = lpfarmkeeper.NewKeeper(
		appCodec,
		keys[lpfarmtypes.StoreKey],
		app.GetSubspace(lpfarmtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.LiquidityKeeper,
	)
	app.LiquidFarmingKeeper = liquidfarmingkeeper.NewKeeper(
		appCodec,
		keys[liquidfarmingtypes.StoreKey],
		app.GetSubspace(liquidfarmingtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.LPFarmKeeper,
		app.LiquidityKeeper,
		&app.GuardKeeper,
	)

	app.NFTKeeper = nftkeeper.NewKeeper(keys[nftkeeper.StoreKey], appCodec, app.AccountKeeper, app.BankKeeper)

	app.TokenKeeper = *tokenkeeper.NewKeeper(
		appCodec,
		keys[tokentypes.StoreKey],
		keys[tokentypes.MemStoreKey],
		app.GetSubspace(tokentypes.ModuleName),
		app.NFTKeeper,
	)

	app.CoinFactoryKeeper = *coinfactorykeeper.NewKeeper(
		appCodec,
		keys[coinfactorytypes.StoreKey],
		app.GetSubspace(coinfactorytypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
	)

	app.GuardKeeper = guardkeeper.NewKeeper(
		appCodec,
		keys[guardtypes.StoreKey],
		keys[guardtypes.MemStoreKey],
		app.GetSubspace(guardtypes.ModuleName),
		app.ModuleAccountAddrs(),
		app.MsgServiceRouter(),
		app.AccountKeeper,
		app.BankKeeper,
		app.AuthzKeeper,
		app.TokenKeeper,
		app.NFTKeeper,
		app.CoinFactoryKeeper,
	)

	app.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		keys[ibctransfertypes.StoreKey],
		app.GetSubspace(ibctransfertypes.ModuleName),
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		scopedTransferKeeper,
	)
	app.transferModule = transfer.NewAppModule(app.TransferKeeper)
	transferIBCModule := transfer.NewIBCModule(app.TransferKeeper)

	app.ICAHostKeeper = icahostkeeper.NewKeeper(
		appCodec,
		keys[icahosttypes.StoreKey],
		app.GetSubspace(icahosttypes.SubModuleName),
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		scopedICAHostKeeper,
		app.MsgServiceRouter(),
	)
	app.icaModule = ica.NewAppModule(nil, &app.ICAHostKeeper)
	icaHostIBCModule := icahost.NewIBCModule(app.ICAHostKeeper)

	// create evidence keeper with router
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec,
		keys[evidencetypes.StoreKey],
		&app.ConsumerKeeper,
		app.SlashingKeeper,
	)
	app.EvidenceKeeper = *evidenceKeeper

	// create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()

	ibcRouter.
		AddRoute(ibctransfertypes.ModuleName, transferIBCModule).
		AddRoute(icahosttypes.SubModuleName, icaHostIBCModule).
		AddRoute(ibcconsumertypes.ModuleName, consumerModule)
	app.IBCKeeper.SetRouter(ibcRouter)

	/****  Module Options ****/
	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		auth.NewAppModule(appCodec, app.AccountKeeper, nil),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, app.GuardKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		crisis.NewAppModule(&app.CrisisKeeper, skipGenesisInvariants),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.ConsumerKeeper),
		budget.NewAppModule(appCodec, app.BudgetKeeper, app.AccountKeeper, app.BankKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		ibc.NewAppModule(app.IBCKeeper),
		params.NewAppModule(app.ParamsKeeper),
		liquidity.NewAppModule(appCodec, app.LiquidityKeeper, app.AccountKeeper, app.BankKeeper),
		liquidfarming.NewAppModule(appCodec, app.LiquidFarmingKeeper, app.AccountKeeper, app.BankKeeper),
		marketmaker.NewAppModule(appCodec, app.MarketMakerKeeper, app.AccountKeeper, app.BankKeeper),
		lpfarm.NewAppModule(appCodec, app.LPFarmKeeper, app.AccountKeeper, app.BankKeeper, app.LiquidityKeeper),
		app.transferModule,
		consumerModule,
		app.icaModule,
		token.NewAppModule(appCodec, app.TokenKeeper, app.AccountKeeper, app.BankKeeper),
		guard.NewAppModule(appCodec, app.GuardKeeper, app.AccountKeeper, app.NFTKeeper, app.CoinFactoryKeeper),
		coinfactory.NewAppModule(appCodec, app.CoinFactoryKeeper, app.AccountKeeper, app.BankKeeper),
		nft.NewAppModule(appCodec, app.NFTKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	app.mm.SetOrderBeginBlockers(
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		budgettypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		liquiditytypes.ModuleName,
		liquidfarmingtypes.ModuleName,
		ibchost.ModuleName,
		lpfarmtypes.ModuleName,

		// empty logic modules
		authtypes.ModuleName,
		banktypes.ModuleName,
		crisistypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,
		ibctransfertypes.ModuleName,
		marketmakertypes.ModuleName,
		icatypes.ModuleName,
		guardtypes.ModuleName,
		nfttypes.ModuleName,
		tokentypes.ModuleName,
		coinfactorytypes.ModuleName,
		ibcconsumertypes.ModuleName,
	)
	app.mm.SetOrderEndBlockers(
		// EndBlocker of crisis module called AssertInvariants
		crisistypes.ModuleName,
		liquiditytypes.ModuleName,
		liquidfarmingtypes.ModuleName,

		// empty logic modules
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		ibchost.ModuleName,
		ibctransfertypes.ModuleName,
		budgettypes.ModuleName,
		marketmakertypes.ModuleName,
		lpfarmtypes.ModuleName,
		icatypes.ModuleName,
		guardtypes.ModuleName,
		nfttypes.ModuleName,
		tokentypes.ModuleName,
		coinfactorytypes.ModuleName,
		ibcconsumertypes.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	app.mm.SetOrderInitGenesis(
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		slashingtypes.ModuleName,
		ibchost.ModuleName,
		evidencetypes.ModuleName,
		ibctransfertypes.ModuleName,
		feegrant.ModuleName,
		authz.ModuleName,
		budgettypes.ModuleName,
		liquiditytypes.ModuleName,
		liquidfarmingtypes.ModuleName,
		marketmakertypes.ModuleName,
		lpfarmtypes.ModuleName,

		// empty logic modules
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		icatypes.ModuleName,

		guardtypes.ModuleName,
		nfttypes.ModuleName,
		tokentypes.ModuleName,
		coinfactorytypes.ModuleName,
		ibcconsumertypes.ModuleName,

		// InitGenesis of crisis module called AssertInvariants
		crisistypes.ModuleName,
	)

	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)

	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.mm.RegisterServices(app.configurator)

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	anteHandler, err := NewAnteHandler(
		HandlerOptions{
			HandlerOptions: ante.HandlerOptions{
				AccountKeeper:   app.AccountKeeper,
				BankKeeper:      app.BankKeeper,
				FeegrantKeeper:  app.FeeGrantKeeper,
				SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
				SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
			},
			IBCKeeper:   app.IBCKeeper,
			TokenKeeper: app.TokenKeeper,
			GuardKeeper: app.GuardKeeper,
		},
	)
	if err != nil {
		panic(fmt.Errorf("failed to create AnteHandler: %s", err))
	}

	app.SetAnteHandler(anteHandler)
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper
	app.ScopedICAHostKeeper = scopedICAHostKeeper

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(fmt.Sprintf("failed to load latest version: %s", err))
		}
	}

	return app
}

// Name returns the name of the App.
func (app *App) Name() string { return app.BaseApp.Name() }

// BeginBlocker application updates every begin block.
func (app *App) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block.
func (app *App) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// InitChainer application update at chain initialization
func (app *App) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}

	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())
	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// LoadHeight loads a particular height.
func (app *App) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}
	return dupMaccPerms
}

// BlockedModuleAccountAddrs returns all the app's blocked module account
// addresses.
func (app *App) BlockedModuleAccountAddrs() map[string]bool {
	modAccAddrs := app.ModuleAccountAddrs()
	for acc := range GetMaccPerms() {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	// allow the following addresses to receive funds
	delete(modAccAddrs, authtypes.NewModuleAddress(
		ibcconsumertypes.ConsumerToSendToProviderName).String())

	return modAccAddrs
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *App) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// LegacyAmino returns App's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *App) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

// AppCodec returns App's app codec.
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

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetKey(storeKey string) *sdk.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetTKey(storeKey string) *sdk.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *App) GetMemKey(storeKey string) *sdk.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *App) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	rpc.RegisterRoutes(clientCtx, apiSvr.Router)
	// Register legacy tx routes.
	authrest.RegisterTxRoutes(clientCtx, apiSvr.Router)
	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register legacy and grpc-gateway routes for all modules.
	ModuleBasics.RegisterRESTRoutes(clientCtx, apiSvr.Router)
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register swagger API from root so that other applications can override easily
	if apiConfig.Swagger {
		RegisterSwaggerAPI(apiSvr.Router)
	}
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *App) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *App) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.interfaceRegistry)
}

// RegisterSwaggerAPI registers swagger route with API Server
func RegisterSwaggerAPI(rtr *mux.Router) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	staticServer := http.FileServer(statikFS)
	rtr.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", staticServer))
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey sdk.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibchost.ModuleName)
	paramsKeeper.Subspace(budgettypes.ModuleName)
	paramsKeeper.Subspace(liquiditytypes.ModuleName)
	paramsKeeper.Subspace(liquidfarmingtypes.ModuleName)
	paramsKeeper.Subspace(marketmakertypes.ModuleName)
	paramsKeeper.Subspace(lpfarmtypes.ModuleName)
	paramsKeeper.Subspace(icahosttypes.SubModuleName)
	paramsKeeper.Subspace(tokentypes.ModuleName)
	paramsKeeper.Subspace(guardtypes.ModuleName)
	paramsKeeper.Subspace(coinfactorytypes.ModuleName)
	paramsKeeper.Subspace(ibcconsumertypes.ModuleName)

	return paramsKeeper
}
