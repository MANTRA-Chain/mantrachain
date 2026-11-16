package v8_9

import (
	"cosmossdk.io/store/types"
	"github.com/MANTRA-Chain/mantrachain/v8/app/upgrades"
)

const (
	// UpgradeName defines the on-chain upgrade name.
	UpgradeName = "v8.9.0"

	// The store keys, legacy params subspaces and module name of the removed
	// interchain accounts module. They are the values of ibc-go's
	// icacontrollertypes.StoreKey/SubModuleName, icahosttypes.StoreKey/SubModuleName
	// and icatypes.ModuleName, copied so the app no longer imports the module.
	ICAControllerStoreKey = "icacontroller"
	ICAHostStoreKey       = "icahost"
	ICAModuleName         = "interchainaccounts"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: types.StoreUpgrades{
		Added: []string{},
		// The interchain accounts module is removed.
		Deleted: []string{ICAControllerStoreKey, ICAHostStoreKey},
	},
}
