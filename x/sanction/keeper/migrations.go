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

// Migrate1to2 sets the new params to their defaults and clears any EIP-7702
// delegation from accounts already on the blacklist. It removes any module
// account from the blacklist, since the new send restriction would stop it
// sending. Other entries are left untouched. The sanctioner list starts
// empty, which keeps the module governance-only until governance appoints
// sanctioners.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	if err := m.keeper.Params.Set(ctx, types.DefaultParams()); err != nil {
		return err
	}
	iter, err := m.keeper.BlacklistAccounts.Iterate(ctx, nil)
	if err != nil {
		return err
	}
	addrs, err := iter.Keys()
	if err != nil {
		return err
	}
	for _, addr := range addrs {
		isModule, err := m.keeper.isModuleAccount(ctx, addr)
		if err != nil {
			return err
		}
		if isModule {
			if err := m.keeper.BlacklistAccounts.Remove(ctx, addr); err != nil {
				return err
			}
			continue
		}
		if err := m.keeper.clearCodeDelegation(ctx, addr); err != nil {
			return err
		}
	}
	return nil
}
