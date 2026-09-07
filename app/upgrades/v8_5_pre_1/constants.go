package v8_5_pre_1

import (
	"cosmossdk.io/store/types"
	"github.com/MANTRA-Chain/mantrachain/v8/app/upgrades"
)

const (
	// UpgradeName defines the on-chain upgrade name.
	UpgradeName = "v8.5.0-pre.1"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: types.StoreUpgrades{
		Added:   []string{},
		Deleted: []string{},
	},
}
