package v5

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/MANTRA-Chain/mantrachain/v5/app/upgrades"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	feemarkettypes "github.com/cosmos/evm/x/feemarket/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
)

// skip-mev feemarket store prefixes to delete from feemarket store
const (
	skipFeeMarketPrefixParams = iota + 1
	skipFeeMarketPrefixState
	skipFeeMarketPrefixEnableHeight = 3
)

var (
	skipKeyParams        = []byte{skipFeeMarketPrefixParams}
	skipKeyState         = []byte{skipFeeMarketPrefixState}
	skipKeyEnabledHeight = []byte{skipFeeMarketPrefixEnableHeight}
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

		// delete skip-mev feemarket state
		feemarketstore := runtime.NewKVStoreService(storekeys[feemarkettypes.StoreKey]).OpenKVStore(c)
		if err := feemarketstore.Delete(skipKeyParams); err != nil {
			return vm, err
		}
		if err := feemarketstore.Delete(skipKeyState); err != nil {
			return vm, err
		}
		if err := feemarketstore.Delete(skipKeyEnabledHeight); err != nil {
			return vm, err
		}

		// delete feemarket from the version map so it can be reinitialized to evm feemarket
		delete(vm, feemarkettypes.StoreKey)

		// burn excess fees stuck in the feemarket fee collector account that are stake tokens
		legacyFeeMarketFeeCollectorAddr := authtypes.NewModuleAddress("feemarket-fee-collector")
		balances := keepers.BankKeeper.GetAllBalances(ctx, legacyFeeMarketFeeCollectorAddr)
		// return an error if the legacy feemarket fee collector account has more than one balance
		if balances.Len() > 1 {
			return vm, errorsmod.Wrapf(
				errortypes.ErrLogic,
				"cannot run upgrade v5, legacy feemarket fee collector account has more than one balance: %s",
				balances.String(),
			)
		}
		if !balances.IsZero() {
			if balances[0].Denom != evmtypes.GetEVMCoinDenom() {
				return vm, errorsmod.Wrapf(
					errortypes.ErrLogic,
					"cannot run upgrade v5, legacy feemarket fee collector account %s has non-stake balance: %s",
					legacyFeeMarketFeeCollectorAddr.String(),
					balances[0].String(),
				)
			}
			// send the balance to the gov module account and burn it
			err := keepers.BankKeeper.SendCoinsFromAccountToModule(ctx, legacyFeeMarketFeeCollectorAddr, govtypes.ModuleName, balances)
			if err != nil {
				return vm, err
			}
			if err := keepers.BankKeeper.BurnCoins(ctx, govtypes.ModuleName, balances); err != nil {
				return vm, err
			}
		}

		// remove the legacy feemarket fee collector account
		legacyFeeMarketFeeCollectorAcc := keepers.AccountKeeper.NewAccountWithAddress(ctx, legacyFeeMarketFeeCollectorAddr)
		keepers.AccountKeeper.RemoveAccount(ctx, legacyFeeMarketFeeCollectorAcc)

		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		// Reset the feemarket module with new params
		feemarketKeeper := keepers.FeeMarketKeeper
		params := feemarkettypes.DefaultParams()
		params.MinGasPrice = sdkmath.LegacyMustNewDecFromStr("0.01")
		params.BaseFee = params.MinGasPrice
		if err := feemarketKeeper.SetParams(ctx, params); err != nil {
			return vm, err
		}

		// set the evm/vm params
		evmParams := evmtypes.DefaultParams()
		evmParams.EvmDenom = evmtypes.GetEVMCoinDenom()
		// enable AllowUnprotectedTxs see adr-006
		evmParams.AllowUnprotectedTxs = true
		if err := keepers.EVMKeeper.SetParams(ctx, evmParams); err != nil {
			return vm, err
		}

		// add burner authority to the fee collector
		macc := authtypes.NewEmptyModuleAccount(authtypes.FeeCollectorName, authtypes.Burner)
		maccI := (keepers.AccountKeeper.NewAccount(ctx, macc)).(sdk.ModuleAccountI) // set the account number
		keepers.AccountKeeper.SetModuleAccount(ctx, maccI)

		// add erc20 address for all existing tokenfactory tokens
		iter := keepers.TokenFactoryKeeper.GetAllDenomsIterator(ctx)
		defer iter.Close()
		for ; iter.Valid(); iter.Next() {
			denom := string(iter.Value())
			if err := keepers.TokenFactoryKeeper.UpdateDenomWithERC20(ctx, denom); err != nil {
				return vm, err
			}
		}

		ctx.Logger().Info("Upgrade v5 complete")
		return vm, nil
	}
}
