package v7

import (
	"cosmossdk.io/store/types"
	"github.com/MANTRA-Chain/mantrachain/v7/app/upgrades"
)

const (
	// UpgradeName defines the on-chain upgrade name.
	UpgradeName = "v7.0.0"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: types.StoreUpgrades{
		Added:   []string{},
		Deleted: []string{},
	},
}
