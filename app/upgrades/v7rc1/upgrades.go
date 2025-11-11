package v7rc1

import (
	"context"
	"math/big"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/MANTRA-Chain/mantrachain/v7/app/upgrades"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	evmkeeper "github.com/cosmos/evm/x/vm/keeper"
	"github.com/ethereum/go-ethereum/common"
)

const (
	NameSlot   = 1
	SymbolSlot = 2
)

var WOMContractAddress = map[string][]common.Address{
	"mantra-1":        {common.HexToAddress("0xE3047710EF6cB36Bcf1E58145529778eA7Cb5598")},
	"mantra-dukong-1": {common.HexToAddress("0x10d26F0491fA11c5853ED7C1f9817b098317DC46")},
	"mantra-canary-net-1": {
		common.HexToAddress("0x523A024258fc56E4d6d79D4367a98F2548A9f401"),
		common.HexToAddress("0xba44F0669812E24fF5826Ad4302e7a0BAfBa39C5"), // e2e test
	},
	"mantra-dryrun-1": {common.HexToAddress("0x523A024258fc56E4d6d79D4367a98F2548A9f401")},
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *upgrades.UpgradeKeepers,
	storekeys map[string]*storetypes.KVStoreKey,
) upgradetypes.UpgradeHandler {
	return func(c context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(c)
		ctx.Logger().Info("Starting module migrations...")

		migrateWOMs(ctx, keepers.EVMKeeper)

		ctx.Logger().Info("Upgrade v7.0.0-rc1 complete")
		return vm, nil
	}
}

func migrateWOMs(ctx sdk.Context, evmKeeper evmkeeper.Keeper) {
	addresses := WOMContractAddress[ctx.ChainID()]
	for _, addr := range addresses {
		migrateWOM(ctx, evmKeeper, addr)
	}
}

func migrateWOM(ctx sdk.Context, evmKeeper evmkeeper.Keeper, contract common.Address) {
	setStringField(ctx, evmKeeper, contract, NameSlot, "Wrapped MANTRA")
	setStringField(ctx, evmKeeper, contract, SymbolSlot, "wMANTRA")
}

func setStringField(ctx sdk.Context, evmKeeper evmkeeper.Keeper, contract common.Address, slot int, value string) {
	if len(value) >= 32 {
		panic("string length exceeds 31 bytes")
	}
	lengthSlot := common.BigToHash(big.NewInt(int64(slot)))

	var data common.Hash
	copy(data[:], []byte(value))
	data[31] = byte(len(value) * 2) // length * 2 for short string

	evmKeeper.SetState(ctx, contract, lengthSlot, data.Bytes())
}
