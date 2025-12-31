package v8providerrc0

import (
	"cosmossdk.io/store/types"
	"github.com/MANTRA-Chain/mantrachain/v8/app/upgrades"
	providertypes "github.com/cosmos/interchain-security/v7/x/ccv/provider/types"
)

const (
	// UpgradeName defines the on-chain upgrade name.
	UpgradeName = "v8.0.0-provider-rc0"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: types.StoreUpgrades{
		Added:   []string{providertypes.ModuleName},
		Deleted: []string{},
	},
}
