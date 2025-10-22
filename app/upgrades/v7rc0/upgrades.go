package v7rc0

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/MANTRA-Chain/mantrachain/v7/app/upgrades"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *upgrades.UpgradeKeepers,
	storekeys map[string]*storetypes.KVStoreKey,
) upgradetypes.UpgradeHandler {
	return func(c context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(c)
		ctx.Logger().Info("Starting module migrations...")

		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		ctx.Logger().Info("Migrating x/bank state...")
		if err = migrateBank(ctx, keepers.BankKeeper, *keepers.TokenFactoryKeeper, keepers.AccountKeeper); err != nil {
			return vm, err
		}

		ctx.Logger().Info("Migrating x/precisebank state...")
		if err = migratePreciseBank(ctx, keepers.PreciseBankKeeper, keepers.BankKeeper); err != nil {
			return vm, err
		}

		ctx.Logger().Info("Migrating x/staking state...")
		if err = migrateStaking(ctx, keepers.StakingKeeper); err != nil {
			return vm, err
		}

		ctx.Logger().Info("Migrating x/gov state...")
		if err = migrateGov(ctx, keepers.GovKeeper); err != nil {
			return vm, err
		}

		ctx.Logger().Info("Migrating x/distribution state...")
		if err = migrateDistr(ctx, keepers.DistrKeeper); err != nil {
			return vm, err
		}

		ctx.Logger().Info("Migrating x/mint state...")
		mintParams, err := keepers.MintKeeper.Params.Get(ctx)
		if err != nil {
			return vm, err
		}
		mintParams.MintDenom = AMANTRA
		if err := keepers.MintKeeper.Params.Set(ctx, mintParams); err != nil {
			return vm, err
		}

		ctx.Logger().Info("Migrating x/crisis state...")
		constantFee, err := keepers.CrisisKeeper.ConstantFee.Get(ctx)
		if err != nil {
			return vm, err
		}
		newConstantFee := convertCoinToNewDenom(constantFee)
		if err := keepers.CrisisKeeper.ConstantFee.Set(ctx, newConstantFee); err != nil {
			return vm, err
		}

		ctx.Logger().Info("Migrating x/evm state...")
		evmParams := keepers.EVMKeeper.GetParams(ctx)
		evmParams.EvmDenom = AMANTRA
		evmParams.ExtendedDenomOptions.ExtendedDenom = AMANTRA
		if err := keepers.EVMKeeper.SetParams(ctx, evmParams); err != nil {
			return vm, err
		}

		ctx.Logger().Info("Migrating x/feegrant state...")
		migrateFeeGrant(ctx, keepers.FeeGrantKeeper)

		ctx.Logger().Info("Migrating x/authz state...")
		migrateAuthz(ctx, keepers.AuthzKeeper)

		// --- Post-Migration ---
		ctx.Logger().Info("Finished v7.0.0-rc0 state migrations.")
		ctx.Logger().Info("Assert Invariants...")
		keepers.CrisisKeeper.AssertInvariants(ctx)

		ctx.Logger().Info("Upgrade v7.0.0-rc0 complete")
		return vm, nil
	}
}
