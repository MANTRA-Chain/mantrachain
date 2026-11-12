package keeper

import (
	"github.com/MANTRA-Chain/mantrachain/v8/x/sanction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Migrator handles in-place store migrations for x/sanction.
type Migrator struct {
	keeper Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(k Keeper) Migrator {
	return Migrator{keeper: k}
}

// Migrate1to2 sets the new params to their defaults. Existing blacklist
// entries are left untouched. The sanctioner list starts empty, which keeps
// the module governance-only until governance appoints sanctioners.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	return m.keeper.Params.Set(ctx, types.DefaultParams())
}
