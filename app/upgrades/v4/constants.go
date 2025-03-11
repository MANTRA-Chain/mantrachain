package v4

import (
	"cosmossdk.io/store/types"
	"github.com/MANTRA-Chain/mantrachain/v4/app/upgrades"
	icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
	icahosttypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/types"
)

const (
	// UpgradeName defines the on-chain upgrade name.
	UpgradeName = "v4"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: types.StoreUpgrades{
		Added:   []string{icacontrollertypes.StoreKey, icahosttypes.StoreKey},
		Deleted: []string{},
	},
}
