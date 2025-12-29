package v7

import (
	"context"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/MANTRA-Chain/mantrachain/v7/app/upgrades"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
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

		// Unblock all module accounts for the duration of the migration.
		// This is temporary and only applies to the bank keeper instance used within this upgrade handler.
		keepers.BankKeeper = keepers.BankKeeper.WithBlockedAddrs(nil)

		migrationCtx := ctx.WithEventManager(sdk.NewEventManager())

		ctx.Logger().Info("Migrating x/auth state...")
		if err = migrateAuth(migrationCtx, keepers.AccountKeeper); err != nil {
			return vm, err
		}

		ctx.Logger().Info("Migrating x/precisebank state...")
		if err = migratePreciseBank(migrationCtx, keepers.PreciseBankKeeper, keepers.BankKeeper, keepers.AccountKeeper); err != nil {
			return vm, err
		}

		ctx.Logger().Info("Migrating x/bank state...")
		if err = migrateBank(migrationCtx, keepers.BankKeeper, *keepers.TokenFactoryKeeper, keepers.AccountKeeper); err != nil {
			return vm, err
		}

		ctx.Logger().Info("Migrating x/staking state...")
		if err = migrateStaking(migrationCtx, keepers.StakingKeeper, storekeys[stakingtypes.ModuleName]); err != nil {
			return vm, err
		}

		ctx.Logger().Info("Migrating x/gov state...")
		if err = migrateGov(migrationCtx, keepers.GovKeeper); err != nil {
			return vm, err
		}

		ctx.Logger().Info("Migrating x/distribution state...")
		if err = migrateDistr(migrationCtx, keepers.DistrKeeper); err != nil {
			return vm, err
		}

		ctx.Logger().Info("Migrating x/mint state...")
		mintParams, err := keepers.MintKeeper.Params.Get(ctx)
		if err != nil {
			return vm, err
		}
		mintParams.MintDenom = AMANTRA
		mintParams.MaxSupply = sdkmath.NewIntWithDecimal(10_000_000_000, 18)
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
		if err := keepers.EVMKeeper.InitEvmCoinInfo(ctx); err != nil {
			return vm, err
		}

		ctx.Logger().Info("Migrating x/feemarket state...")
		feemarketParams := keepers.FeeMarketKeeper.GetParams(ctx)
		feemarketParams.BaseFee = feemarketParams.BaseFee.Mul(ScalingFactor.ToLegacyDec())
		feemarketParams.MinGasPrice = feemarketParams.MinGasPrice.Mul(ScalingFactor.ToLegacyDec())
		if err := keepers.FeeMarketKeeper.SetParams(ctx, feemarketParams); err != nil {
			return vm, err
		}

		ctx.Logger().Info("Migrating x/feegrant state...")
		migrateFeeGrant(migrationCtx, keepers.FeeGrantKeeper)

		ctx.Logger().Info("Migrating x/authz state...")
		migrateAuthz(migrationCtx, keepers.AuthzKeeper)

		ctx.Logger().Info("Migrating wom contract state...")
		if err := migrateWOMs(ctx, keepers.EVMKeeper); err != nil {
			return vm, err
		}

		ctx.Logger().Info("Migrating token factory state...")
		tokenFactoryParams := keepers.TokenFactoryKeeper.GetParams(ctx)
		tokenFactoryParams.DenomCreationFee = convertCoinsToNewDenom(tokenFactoryParams.DenomCreationFee)
		if err := keepers.TokenFactoryKeeper.SetParams(ctx, tokenFactoryParams); err != nil {
			return vm, err
		}

		ctx.Logger().Info("Removing OM from oracle...")
		if keepers.OracleKeeper != nil {
			pairs := keepers.OracleKeeper.GetAllCurrencyPairs(ctx)
			for _, pair := range pairs {
				if pair.Base == "OM" || pair.Quote == "OM" {
					if err := keepers.OracleKeeper.RemoveCurrencyPair(ctx, pair); err != nil {
						ctx.Logger().Error("Failed to remove OM currency pair", "pair", pair.String(), "err", err)
					}
				}
			}
		}

		// --- Post-Migration ---
		ctx.Logger().Info("Finished v7.0.0-rc4 state migrations.")
		ctx.Logger().Info("Assert Invariants...")
		keepers.CrisisKeeper.AssertInvariants(ctx)

		ctx.Logger().Info("Upgrade v7.0.0-rc4 complete")
		return vm, nil
	}
}
