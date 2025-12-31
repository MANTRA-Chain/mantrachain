package v7providerrc0

import (
	"cosmossdk.io/store/types"
	"github.com/MANTRA-Chain/mantrachain/v7/app/upgrades"
	providertypes "github.com/cosmos/interchain-security/v7/x/ccv/provider/types"
)

const (
	// UpgradeName defines the on-chain upgrade name.
	UpgradeName = "v7.0.0-provider-rc0"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: types.StoreUpgrades{
		Added:   []string{providertypes.ModuleName},
		Deleted: []string{},
	},
}
