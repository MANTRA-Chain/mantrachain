package keeper

import (
	"github.com/MANTRA-Finance/mantrachain/x/token/exported"
	v3 "github.com/MANTRA-Finance/mantrachain/x/token/migrations/v3"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper         Keeper
	legacySubspace exported.Subspace
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper, legacySubspace exported.Subspace) Migrator {
	return Migrator{
		keeper:         keeper,
		legacySubspace: legacySubspace,
	}
}

// Migrate1to2 migrates the store from consensus version 1 to 2
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	return v3.MigrateStore(ctx, m.keeper.storeService, m.keeper.cdc, m.legacySubspace)
}
