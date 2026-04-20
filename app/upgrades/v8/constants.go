package v8

import (
	"cosmossdk.io/store/types"
	"github.com/MANTRA-Chain/mantrachain/v8/app/upgrades"
	precisebanktypes "github.com/cosmos/evm/x/precisebank/types"
	providertypes "github.com/cosmos/interchain-security/v7/x/ccv/provider/types"
	marketmaptypes "github.com/skip-mev/connect/v2/x/marketmap/types"
	oracletypes "github.com/skip-mev/connect/v2/x/oracle/types"
)

const (
	// UpgradeName defines the on-chain upgrade name.
	UpgradeName = "v8.0.0"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: types.StoreUpgrades{
		Added:   []string{providertypes.ModuleName},
		Deleted: []string{precisebanktypes.ModuleName, oracletypes.ModuleName, marketmaptypes.ModuleName},
	},
}
