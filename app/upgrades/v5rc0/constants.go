package v5rc0

import (
	"cosmossdk.io/store/types"
	"github.com/MANTRA-Chain/mantrachain/v5/app/upgrades"
	erc20types "github.com/cosmos/evm/x/erc20/types"
	precisebanktypes "github.com/cosmos/evm/x/precisebank/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
)

const (
	// UpgradeName defines the on-chain upgrade name.
	UpgradeName = "v5"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: types.StoreUpgrades{
		Added:   []string{evmtypes.StoreKey, erc20types.StoreKey, precisebanktypes.StoreKey},
		Deleted: []string{"capability", "feeibc", "hooks-for-ibc"},
	},
}
