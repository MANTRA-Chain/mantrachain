package keeper

import (
	"github.com/MANTRA-Finance/mantrachain/x/lpfarm/exported"
	v2 "github.com/MANTRA-Finance/mantrachain/x/lpfarm/migrations/v2"
	"github.com/MANTRA-Finance/mantrachain/x/lpfarm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper         Keeper
	legacySubspace exported.Subspace
	guardKeeper    types.GuardKeeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper, legacySubspace exported.Subspace, guardKeeper types.GuardKeeper) Migrator {
	return Migrator{
		keeper:         keeper,
		legacySubspace: legacySubspace,
		guardKeeper:    guardKeeper,
	}
}

// Migrate1to2 migrates the store from consensus version 1 to 2
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	return v2.MigrateStore(ctx, m.keeper.storeService, m.keeper.cdc, m.legacySubspace, m.guardKeeper)
}
