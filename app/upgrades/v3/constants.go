package v2

import (
	"cosmossdk.io/store/types"
	"github.com/MANTRA-Chain/mantrachain/v3/app/upgrades"
	sanctiontypes "github.com/MANTRA-Chain/mantrachain/v3/x/sanction/types"
)

const (
	// UpgradeName defines the on-chain upgrade name.
	UpgradeName = "v3"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: types.StoreUpgrades{
		Added:   []string{sanctiontypes.StoreKey},
		Deleted: []string{},
	},
}
