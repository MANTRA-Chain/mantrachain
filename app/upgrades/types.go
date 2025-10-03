package upgrades

import (
	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sanctionkeeper "github.com/MANTRA-Chain/mantrachain/v6/x/sanction/keeper"
	tokenfactorykeeper "github.com/MANTRA-Chain/mantrachain/v6/x/tokenfactory/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	erc20keeper "github.com/cosmos/evm/x/erc20/keeper"
	feemarketkeeper "github.com/cosmos/evm/x/feemarket/keeper"
	transferkeeper "github.com/cosmos/evm/x/ibc/transfer/keeper"
	evmkeeper "github.com/cosmos/evm/x/vm/keeper"
	channelkeeper "github.com/cosmos/ibc-go/v10/modules/core/04-channel/keeper"
)

// Upgrade defines a struct containing necessary fields that a SoftwareUpgradeProposal
// must have written, in order for the state migration to go smoothly.
// An upgrade must implement this struct, and then set it in the app.go.
// The app.go will then define the handler.
type Upgrade struct {
	// Upgrade version name, for the upgrade handler, e.g. `v7`
	UpgradeName string

	// CreateUpgradeHandler defines the function that creates an upgrade handler
	CreateUpgradeHandler func(*module.Manager, module.Configurator, *UpgradeKeepers, map[string]*storetypes.KVStoreKey) upgradetypes.UpgradeHandler

	// Store upgrades, should be used for any new modules introduced, new modules deleted, or store names renamed.
	StoreUpgrades storetypes.StoreUpgrades
}

// Fork defines a struct containing the requisite fields for a non-software upgrade proposal
// Hard Fork at a given height to implement.
// There is one time code that can be added for the start of the Fork, in `BeginForkLogic`.
// Any other change in the code should be height-gated, if the goal is to have old and new binaries
// to be compatible prior to the upgrade height.
type Fork struct {
	// Upgrade version name, for the upgrade handler, e.g. `v7`
	UpgradeName string
	// height the upgrade occurs at
	UpgradeHeight int64

	// Function that runs some custom state transition code at the beginning of a fork.
	BeginForkLogic func(ctx sdk.Context)
}

type UpgradeKeepers struct {
	// keepers
	ChannelKeeper      *channelkeeper.Keeper
	TransferKeeper     transferkeeper.Keeper
	TokenFactoryKeeper *tokenfactorykeeper.Keeper
	// v3
	SanctionKeeper sanctionkeeper.Keeper
	// v5
	FeeMarketKeeper feemarketkeeper.Keeper
	AccountKeeper   authkeeper.AccountKeeper
	BankKeeper      bankkeeper.BaseKeeper
	EVMKeeper       evmkeeper.Keeper
	Erc20Keeper     erc20keeper.Keeper
}
