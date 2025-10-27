package v61

import (
	"cosmossdk.io/store/types"
	"github.com/MANTRA-Chain/mantrachain/v6/app/upgrades"
)

const (
	// UpgradeName defines the on-chain upgrade name.
	UpgradeName = "v6.1.0"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades:        types.StoreUpgrades{},
}
