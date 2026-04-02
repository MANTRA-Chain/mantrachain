package v8rc3

import (
	"cosmossdk.io/store/types"
	"github.com/MANTRA-Chain/mantrachain/v8/app/upgrades"
)

const (
	// UpgradeName defines the on-chain upgrade name.
	UpgradeName = "v8.0.0-rc3"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: types.StoreUpgrades{
		Added:   []string{},
		Deleted: []string{"oracle", "marketmap"},
	},
}
