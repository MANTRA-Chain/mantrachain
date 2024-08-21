package app

import (
	"time"

	runtimev1alpha1 "cosmossdk.io/api/cosmos/app/runtime/v1alpha1"
	appv1alpha1 "cosmossdk.io/api/cosmos/app/v1alpha1"
	authmodulev1 "cosmossdk.io/api/cosmos/auth/module/v1"
	authzmodulev1 "cosmossdk.io/api/cosmos/authz/module/v1"
	bankmodulev1 "cosmossdk.io/api/cosmos/bank/module/v1"
	circuitmodulev1 "cosmossdk.io/api/cosmos/circuit/module/v1"
	consensusmodulev1 "cosmossdk.io/api/cosmos/consensus/module/v1"
	crisismodulev1 "cosmossdk.io/api/cosmos/crisis/module/v1"
	distrmodulev1 "cosmossdk.io/api/cosmos/distribution/module/v1"
	evidencemodulev1 "cosmossdk.io/api/cosmos/evidence/module/v1"
	feegrantmodulev1 "cosmossdk.io/api/cosmos/feegrant/module/v1"
	genutilmodulev1 "cosmossdk.io/api/cosmos/genutil/module/v1"
	govmodulev1 "cosmossdk.io/api/cosmos/gov/module/v1"
	groupmodulev1 "cosmossdk.io/api/cosmos/group/module/v1"
	mintmodulev1 "cosmossdk.io/api/cosmos/mint/module/v1"
	nftmodulev1 "cosmossdk.io/api/cosmos/nft/module/v1"
	paramsmodulev1 "cosmossdk.io/api/cosmos/params/module/v1"
	slashingmodulev1 "cosmossdk.io/api/cosmos/slashing/module/v1"
	stakingmodulev1 "cosmossdk.io/api/cosmos/staking/module/v1"
	txconfigv1 "cosmossdk.io/api/cosmos/tx/config/v1"
	upgrademodulev1 "cosmossdk.io/api/cosmos/upgrade/module/v1"
	vestingmodulev1 "cosmossdk.io/api/cosmos/vesting/module/v1"
	"cosmossdk.io/core/appconfig"
	circuittypes "cosmossdk.io/x/circuit/types"
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/feegrant"
	"cosmossdk.io/x/nft"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/group"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	ibcfeetypes "github.com/cosmos/ibc-go/v8/modules/apps/29-fee/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	marketmapmodulev1 "github.com/skip-mev/slinky/api/slinky/marketmap/module/v1"
	oraclemodulev1 "github.com/skip-mev/slinky/api/slinky/oracle/module/v1"
	_ "github.com/skip-mev/slinky/x/marketmap" // import for side-effects
	marketmaptypes "github.com/skip-mev/slinky/x/marketmap/types"
	_ "github.com/skip-mev/slinky/x/oracle" // import for side-effects
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	"google.golang.org/protobuf/types/known/durationpb"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	airdropmodulev1 "github.com/MANTRA-Finance/mantrachain/api/mantrachain/airdrop/module/v1"
	bridgemodulev1 "github.com/MANTRA-Finance/mantrachain/api/mantrachain/bridge/module/v1"
	coinfactorymodulev1 "github.com/MANTRA-Finance/mantrachain/api/mantrachain/coinfactory/module/v1"
	didmodulev1 "github.com/MANTRA-Finance/mantrachain/api/mantrachain/did/module/v1"
	guardmodulev1 "github.com/MANTRA-Finance/mantrachain/api/mantrachain/guard/module/v1"
	liquiditymodulev1 "github.com/MANTRA-Finance/mantrachain/api/mantrachain/liquidity/module/v1"
	lpfarmmodulev1 "github.com/MANTRA-Finance/mantrachain/api/mantrachain/lpfarm/module/v1"
	marketmakermodulev1 "github.com/MANTRA-Finance/mantrachain/api/mantrachain/marketmaker/module/v1"
	tokenmodulev1 "github.com/MANTRA-Finance/mantrachain/api/mantrachain/token/module/v1"
	txfeesmodulev1 "github.com/MANTRA-Finance/mantrachain/api/mantrachain/txfees/module/v1"
	_ "github.com/MANTRA-Finance/mantrachain/x/airdrop/module" // import for side-effects
	airdropmoduletypes "github.com/MANTRA-Finance/mantrachain/x/airdrop/types"
	_ "github.com/MANTRA-Finance/mantrachain/x/bridge/module" // import for side-effects
	bridgemoduletypes "github.com/MANTRA-Finance/mantrachain/x/bridge/types"
	_ "github.com/MANTRA-Finance/mantrachain/x/coinfactory/module" // import for side-effects
	coinfactorymoduletypes "github.com/MANTRA-Finance/mantrachain/x/coinfactory/types"
	_ "github.com/MANTRA-Finance/mantrachain/x/did/module" // import for side-effects
	didmoduletypes "github.com/MANTRA-Finance/mantrachain/x/did/types"
	_ "github.com/MANTRA-Finance/mantrachain/x/guard/module" // import for side-effects
	guardmoduletypes "github.com/MANTRA-Finance/mantrachain/x/guard/types"
	_ "github.com/MANTRA-Finance/mantrachain/x/liquidity/module" // import for side-effects
	liquiditymoduletypes "github.com/MANTRA-Finance/mantrachain/x/liquidity/types"
	_ "github.com/MANTRA-Finance/mantrachain/x/lpfarm/module" // import for side-effects
	lpfarmmoduletypes "github.com/MANTRA-Finance/mantrachain/x/lpfarm/types"
	_ "github.com/MANTRA-Finance/mantrachain/x/marketmaker/module" // import for side-effects
	marketmakermoduletypes "github.com/MANTRA-Finance/mantrachain/x/marketmaker/types"
	_ "github.com/MANTRA-Finance/mantrachain/x/token/module" // import for side-effects
	tokenmoduletypes "github.com/MANTRA-Finance/mantrachain/x/token/types"
	_ "github.com/MANTRA-Finance/mantrachain/x/txfees/module" // import for side-effects
	txfeesmoduletypes "github.com/MANTRA-Finance/mantrachain/x/txfees/types"
	// this line is used by starport scaffolding # stargate/app/moduleImport
)

var (
	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: The genutils module must also occur after auth so that it can access the params from auth.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	genesisModuleOrder = []string{
		// cosmos-sdk/ibc modules
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		crisistypes.ModuleName,
		ibcexported.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		ibctransfertypes.ModuleName,
		icatypes.ModuleName,
		ibcfeetypes.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		circuittypes.ModuleName,
		nft.ModuleName,
		group.ModuleName,
		consensustypes.ModuleName,
		circuittypes.ModuleName,
		wasmtypes.ModuleName,
		// slinky modules
		oracletypes.ModuleName,
		// market map genesis must be called AFTER all consuming modules (i.e. x/oracle, etc.)
		marketmaptypes.ModuleName,
		// chain modules
		bridgemoduletypes.ModuleName,
		airdropmoduletypes.ModuleName,
		coinfactorymoduletypes.ModuleName,
		didmoduletypes.ModuleName,
		tokenmoduletypes.ModuleName,
		guardmoduletypes.ModuleName,
		marketmakermoduletypes.ModuleName,
		liquiditymoduletypes.ModuleName,
		lpfarmmoduletypes.ModuleName,
		txfeesmoduletypes.ModuleName,
		// this line is used by starport scaffolding # stargate/app/initGenesis
	}

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	// NOTE: capability module's beginblocker must come before any modules using capabilities (e.g. IBC)
	beginBlockers = []string{
		// cosmos sdk modules
		minttypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		authz.ModuleName,
		genutiltypes.ModuleName,
		// ibc modules
		capabilitytypes.ModuleName,
		ibcexported.ModuleName,
		ibctransfertypes.ModuleName,
		icatypes.ModuleName,
		ibcfeetypes.ModuleName,
		wasmtypes.ModuleName,
		// slinky modules
		oracletypes.ModuleName,
		marketmaptypes.ModuleName,
		// chain modules
		bridgemoduletypes.ModuleName,
		airdropmoduletypes.ModuleName,
		coinfactorymoduletypes.ModuleName,
		didmoduletypes.ModuleName,
		tokenmoduletypes.ModuleName,
		guardmoduletypes.ModuleName,
		marketmakermoduletypes.ModuleName,
		liquiditymoduletypes.ModuleName,
		lpfarmmoduletypes.ModuleName,
		txfeesmoduletypes.ModuleName,
		// this line is used by starport scaffolding # stargate/app/beginBlockers
	}

	endBlockers = []string{
		// cosmos sdk modules
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		feegrant.ModuleName,
		group.ModuleName,
		genutiltypes.ModuleName,
		// ibc modules
		ibcexported.ModuleName,
		ibctransfertypes.ModuleName,
		capabilitytypes.ModuleName,
		icatypes.ModuleName,
		ibcfeetypes.ModuleName,
		wasmtypes.ModuleName,
		// slinky modules
		oracletypes.ModuleName,
		marketmaptypes.ModuleName,
		// chain modules
		bridgemoduletypes.ModuleName,
		airdropmoduletypes.ModuleName,
		coinfactorymoduletypes.ModuleName,
		didmoduletypes.ModuleName,
		tokenmoduletypes.ModuleName,
		guardmoduletypes.ModuleName,
		marketmakermoduletypes.ModuleName,
		liquiditymoduletypes.ModuleName,
		lpfarmmoduletypes.ModuleName,
		txfeesmoduletypes.ModuleName,
		// this line is used by starport scaffolding # stargate/app/endBlockers
	}

	preBlockers = []string{
		upgradetypes.ModuleName,
		// this line is used by starport scaffolding # stargate/app/preBlockers
	}

	// module account permissions
	moduleAccPerms = []*authmodulev1.ModuleAccountPermission{
		{Account: authtypes.FeeCollectorName},
		{Account: distrtypes.ModuleName},
		{Account: minttypes.ModuleName, Permissions: []string{authtypes.Minter}},
		{Account: stakingtypes.BondedPoolName, Permissions: []string{authtypes.Burner, stakingtypes.ModuleName}},
		{Account: stakingtypes.NotBondedPoolName, Permissions: []string{authtypes.Burner, stakingtypes.ModuleName}},
		{Account: govtypes.ModuleName, Permissions: []string{authtypes.Burner}},
		{Account: nft.ModuleName},
		{Account: ibctransfertypes.ModuleName, Permissions: []string{authtypes.Minter, authtypes.Burner}},
		{Account: ibcfeetypes.ModuleName},
		{Account: icatypes.ModuleName},
		{Account: oracletypes.ModuleName, Permissions: []string{}},
		// this line is used by starport scaffolding # stargate/app/maccPerms

		// Mantrachain
		{Account: airdropmoduletypes.ModuleName},
		{Account: bridgemoduletypes.ModuleName},
		{Account: coinfactorymoduletypes.ModuleName, Permissions: []string{authtypes.Minter, authtypes.Burner}},
		{Account: didmoduletypes.ModuleName},
		{Account: guardmoduletypes.ModuleName},
		{Account: liquiditymoduletypes.ModuleName, Permissions: []string{authtypes.Minter, authtypes.Burner}},
		{Account: lpfarmmoduletypes.ModuleName},
		{Account: marketmakermoduletypes.ModuleName},
		{Account: tokenmoduletypes.ModuleName},
		{Account: txfeesmoduletypes.ModuleName},
	}

	// guard whitelisted transfer acc addresses
	guardWhitelistedTransferAccAddrs = []string{
		ibctransfertypes.ModuleName,
		ibcfeetypes.ModuleName,
		icatypes.ModuleName,
		airdropmoduletypes.ModuleName,
		bridgemoduletypes.ModuleName,
		coinfactorymoduletypes.ModuleName,
		liquiditymoduletypes.ModuleName,
		lpfarmmoduletypes.ModuleName,
		marketmakermoduletypes.ModuleName,
		tokenmoduletypes.ModuleName,
		txfeesmoduletypes.ModuleName,
	}

	// blocked account addresses
	blockAccAddrs = []string{
		authtypes.FeeCollectorName,
		distrtypes.ModuleName,
		minttypes.ModuleName,
		stakingtypes.BondedPoolName,
		stakingtypes.NotBondedPoolName,
		nft.ModuleName,

		//Mantrachain,
		didmoduletypes.ModuleName,
		guardmoduletypes.ModuleName,

		// We allow the following module accounts to receive funds:
		// govtypes.ModuleName
	}

	// appConfig application configuration (used by depinject)
	appConfig = appconfig.Compose(&appv1alpha1.Config{
		Modules: []*appv1alpha1.ModuleConfig{
			{
				Name: runtime.ModuleName,
				Config: appconfig.WrapAny(&runtimev1alpha1.Module{
					AppName:       Name,
					PreBlockers:   preBlockers,
					BeginBlockers: beginBlockers,
					EndBlockers:   endBlockers,
					InitGenesis:   genesisModuleOrder,
					OverrideStoreKeys: []*runtimev1alpha1.StoreKeyConfig{
						{
							ModuleName: authtypes.ModuleName,
							KvStoreKey: "acc",
						},
					},
					// When ExportGenesis is not specified, the export genesis module order
					// is equal to the init genesis order
					// ExportGenesis: genesisModuleOrder,
					// Uncomment if you want to set a custom migration order here.
					// OrderMigrations: nil,
				}),
			},
			{
				Name: authtypes.ModuleName,
				Config: appconfig.WrapAny(&authmodulev1.Module{
					Bech32Prefix:             AccountAddressPrefix,
					ModuleAccountPermissions: moduleAccPerms,
					// By default modules authority is the governance module. This is configurable with the following:
					// Authority: "group", // A custom module authority can be set using a module name
					// Authority: "cosmos1cwwv22j5ca08ggdv9c2uky355k908694z577tv", // or a specific address
				}),
			},
			{
				Name:   nft.ModuleName,
				Config: appconfig.WrapAny(&nftmodulev1.Module{}),
			},
			{
				Name:   vestingtypes.ModuleName,
				Config: appconfig.WrapAny(&vestingmodulev1.Module{}),
			},
			{
				Name: banktypes.ModuleName,
				Config: appconfig.WrapAny(&bankmodulev1.Module{
					BlockedModuleAccountsOverride: blockAccAddrs,
				}),
			},
			{
				Name: stakingtypes.ModuleName,
				Config: appconfig.WrapAny(&stakingmodulev1.Module{
					// NOTE: specifying a prefix is only necessary when using bech32 addresses
					// If not specfied, the auth Bech32Prefix appended with "valoper" and "valcons" is used by default
					Bech32PrefixValidator: AccountAddressPrefix + "valoper",
					Bech32PrefixConsensus: AccountAddressPrefix + "valcons",
				}),
			},
			{
				Name:   slashingtypes.ModuleName,
				Config: appconfig.WrapAny(&slashingmodulev1.Module{}),
			},
			{
				Name:   paramstypes.ModuleName,
				Config: appconfig.WrapAny(&paramsmodulev1.Module{}),
			},
			{
				Name:   "tx",
				Config: appconfig.WrapAny(&txconfigv1.Config{}),
			},
			{
				Name:   oracletypes.ModuleName,
				Config: appconfig.WrapAny(&oraclemodulev1.Module{}),
			},
			{
				Name:   marketmaptypes.ModuleName,
				Config: appconfig.WrapAny(&marketmapmodulev1.Module{}),
			},
			{
				Name:   genutiltypes.ModuleName,
				Config: appconfig.WrapAny(&genutilmodulev1.Module{}),
			},
			{
				Name:   authz.ModuleName,
				Config: appconfig.WrapAny(&authzmodulev1.Module{}),
			},
			{
				Name:   upgradetypes.ModuleName,
				Config: appconfig.WrapAny(&upgrademodulev1.Module{}),
			},
			{
				Name:   distrtypes.ModuleName,
				Config: appconfig.WrapAny(&distrmodulev1.Module{}),
			},
			{
				Name:   evidencetypes.ModuleName,
				Config: appconfig.WrapAny(&evidencemodulev1.Module{}),
			},
			{
				Name:   minttypes.ModuleName,
				Config: appconfig.WrapAny(&mintmodulev1.Module{}),
			},
			{
				Name: group.ModuleName,
				Config: appconfig.WrapAny(&groupmodulev1.Module{
					MaxExecutionPeriod: durationpb.New(time.Second * 1209600),
					MaxMetadataLen:     255,
				}),
			},
			{
				Name:   feegrant.ModuleName,
				Config: appconfig.WrapAny(&feegrantmodulev1.Module{}),
			},
			{
				Name:   govtypes.ModuleName,
				Config: appconfig.WrapAny(&govmodulev1.Module{}),
			},
			{
				Name:   crisistypes.ModuleName,
				Config: appconfig.WrapAny(&crisismodulev1.Module{}),
			},
			{
				Name:   consensustypes.ModuleName,
				Config: appconfig.WrapAny(&consensusmodulev1.Module{}),
			},
			{
				Name:   circuittypes.ModuleName,
				Config: appconfig.WrapAny(&circuitmodulev1.Module{}),
			},
			{
				Name:   airdropmoduletypes.ModuleName,
				Config: appconfig.WrapAny(&airdropmodulev1.Module{}),
			},
			{
				Name:   bridgemoduletypes.ModuleName,
				Config: appconfig.WrapAny(&bridgemodulev1.Module{}),
			},
			{
				Name:   coinfactorymoduletypes.ModuleName,
				Config: appconfig.WrapAny(&coinfactorymodulev1.Module{}),
			},
			{
				Name:   didmoduletypes.ModuleName,
				Config: appconfig.WrapAny(&didmodulev1.Module{}),
			},
			{
				Name: guardmoduletypes.ModuleName,
				Config: appconfig.WrapAny(&guardmodulev1.Module{
					WhitelistedTransferAccAddrs: guardWhitelistedTransferAccAddrs,
				}),
			},
			{
				Name:   tokenmoduletypes.ModuleName,
				Config: appconfig.WrapAny(&tokenmodulev1.Module{}),
			},
			{
				Name:   marketmakermoduletypes.ModuleName,
				Config: appconfig.WrapAny(&marketmakermodulev1.Module{}),
			},
			{
				Name:   liquiditymoduletypes.ModuleName,
				Config: appconfig.WrapAny(&liquiditymodulev1.Module{}),
			},
			{
				Name:   lpfarmmoduletypes.ModuleName,
				Config: appconfig.WrapAny(&lpfarmmodulev1.Module{}),
			},
			{
				Name:   txfeesmoduletypes.ModuleName,
				Config: appconfig.WrapAny(&txfeesmodulev1.Module{}),
			},
			// this line is used by starport scaffolding # stargate/app/moduleConfig
		},
	})
)
