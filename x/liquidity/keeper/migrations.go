package keeper

import (
	"github.com/MANTRA-Finance/mantrachain/x/liquidity/exported"
	v4 "github.com/MANTRA-Finance/mantrachain/x/liquidity/migrations/v4"
	"github.com/MANTRA-Finance/mantrachain/x/liquidity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper         *Keeper
	legacySubspace exported.Subspace
	guardKeeper    types.GuardKeeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper *Keeper, legacySubspace exported.Subspace, guardKeeper types.GuardKeeper) Migrator {
	return Migrator{
		keeper:         keeper,
		legacySubspace: legacySubspace,
		guardKeeper:    guardKeeper,
	}
}

// Migrate3to4 migrates the store from consensus version 3 to 4
func (m Migrator) Migrate3to4(ctx sdk.Context) error {
	return v4.MigrateStore(ctx, m.keeper.storeService, m.keeper.cdc, m.legacySubspace, m.guardKeeper)
}
