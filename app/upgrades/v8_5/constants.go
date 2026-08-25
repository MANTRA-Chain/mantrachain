package v8_5

import (
	"github.com/MANTRA-Chain/mantrachain/v8/app/upgrades"
	"github.com/cosmos/cosmos-sdk/store/v2/types"
)

const (
	// UpgradeName defines the on-chain upgrade name.
	UpgradeName = "v8.5.0"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: types.StoreUpgrades{
		Added:   []string{},
		Deleted: []string{},
	},
}
