package v6providerrc0

import (
	"context"
	"strings"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/MANTRA-Chain/mantrachain/v6/app/upgrades"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	erc20types "github.com/cosmos/evm/x/erc20/types"
	providerkeeper "github.com/cosmos/interchain-security/v7/x/ccv/provider/keeper"
	providertypes "github.com/cosmos/interchain-security/v7/x/ccv/provider/types"
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

		// The provider module is new, so we need to set its version in the version map
		// to avoid running InitGenesis which would return validator updates and cause a conflict.
		vm[providertypes.ModuleName] = 7

		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		// init genesis for provider module
		providerGenesis := providertypes.DefaultGenesisState()
		providerGenesis.Params.ConsumerRewardDenomRegistrationFee = sdk.NewCoin("uom", math.NewInt(1000000))
		providerGenesis.Params.BlocksPerEpoch = 10
		providerGenesis.Params.NumberOfEpochsToStartReceivingRewards = 2
		keepers.ProviderKeeper.InitGenesis(ctx, providerGenesis)

		ctx.Logger().Info("Initializing ConsensusParam Version...")
		err = InitializeConsensusParamVersion(ctx, keepers.ConsensusParamsKeeper)
		if err != nil {
			// don't hard fail here, as this is not critical for the upgrade to succeed
			ctx.Logger().Error("Error initializing ConsensusParam Version:", "message", err.Error())
		}

		ctx.Logger().Info("Initializing MaxProviderConsensusValidators parameter...")
		maxProviderValidators, err := keepers.StakingKeeper.MaxValidators(ctx)
		if err != nil {
			return vm, errorsmod.Wrapf(err, "getting MaxValidators during migration")
		}
		InitializeMaxProviderConsensusParam(ctx, keepers.ProviderKeeper, int64(maxProviderValidators))

		ctx.Logger().Info("Setting MaxValidators parameter...")
		err = SetMaxValidators(ctx, *keepers.StakingKeeper, maxProviderValidators)
		if err != nil {
			return vm, errorsmod.Wrapf(err, "setting MaxValidators during migration")
		}

		ctx.Logger().Info("Initializing LastProviderConsensusValidatorSet...")
		err = InitializeLastProviderConsensusValidatorSet(ctx, keepers.ProviderKeeper, *keepers.StakingKeeper, int(maxProviderValidators))
		if err != nil {
			return vm, errorsmod.Wrapf(err, "initializing LastProviderConsensusValSet during migration")
		}

		// update contract owner for all existing tokenfactory token_pairs
		pairs := keepers.Erc20Keeper.GetTokenPairs(ctx)
		for _, pair := range pairs {
			if strings.HasPrefix(pair.Denom, "factory/") {
				pair.ContractOwner = erc20types.OWNER_MODULE
				keepers.Erc20Keeper.SetTokenPair(ctx, pair)
			}
		}

		disableList := []string{
			"wasm/cosmos.evm.erc20.v1.MsgRegisterERC20",
			"wasm/cosmos.authz.v1beta1.MsgExec",
		}
		for _, msg := range disableList {
			if err := keepers.CircuitKeeper.DisableList.Set(ctx, msg); err != nil {
				return vm, err
			}
		}

		ctx.Logger().Info("Upgrade v6.0.0-rc0 complete")
		return vm, nil
	}
}

// InitializeConsensusParamVersion initializes the consumer params that were missed in a consensus keeper migration.
// Some fields were set to nil values instead of zero values, which causes a panic during Txs to modify the params.
// Context:
// - https://github.com/cosmos/cosmos-sdk/issues/21483
// - https://github.com/cosmos/cosmos-sdk/pull/21484
func InitializeConsensusParamVersion(ctx sdk.Context, consensusKeeper consensusparamkeeper.Keeper) error {
	params, err := consensusKeeper.ParamsStore.Get(ctx)
	if err != nil {
		return err
	}
	params.Version = &cmtproto.VersionParams{}
	return consensusKeeper.ParamsStore.Set(ctx, params)
}

// InitializeMaxProviderConsensusParam initializes the MaxProviderConsensusValidators parameter.
// It is set to the current number of validators participating in consensus on the MANTRA Chain.
// This parameter will be used to govern the number of validators participating in consensus on the MANTRA Chain,
// and takes over this role from the MaxValidators parameter in the staking module.
func InitializeMaxProviderConsensusParam(ctx sdk.Context, providerKeeper providerkeeper.Keeper, maxValidators int64) {
	params := providerKeeper.GetParams(ctx)
	params.MaxProviderConsensusValidators = maxValidators
	providerKeeper.SetParams(ctx, params)
}

// SetMaxValidators sets the MaxValidators parameter in the staking module to the same as MaxProviderConsensusParam
func SetMaxValidators(ctx sdk.Context, stakingKeeper stakingkeeper.Keeper, maxProviderValidators uint32) error {
	params, err := stakingKeeper.GetParams(ctx)
	if err != nil {
		return err
	}

	params.MaxValidators = maxProviderValidators

	return stakingKeeper.SetParams(ctx, params)
}

// InitializeLastProviderConsensusValidatorSet initializes the last provider consensus validator set
// by setting it to the first 180 validators from the current validator set of the staking module.
func InitializeLastProviderConsensusValidatorSet(
	ctx sdk.Context, providerKeeper providerkeeper.Keeper, stakingKeeper stakingkeeper.Keeper, maxProviderValidators int,
) error {
	vals, err := stakingKeeper.GetBondedValidatorsByPower(ctx)
	if err != nil {
		return err
	}

	// cut the validator set to the current maxProviderValidators
	if len(vals) > maxProviderValidators {
		vals = vals[:maxProviderValidators]
	}

	// create consensus validators for the staking validators
	lastValidators := []providertypes.ConsensusValidator{}
	for _, val := range vals {
		consensusVal, err := providerKeeper.CreateProviderConsensusValidator(ctx, val)
		if err != nil {
			return err
		}

		lastValidators = append(lastValidators, consensusVal)
	}

	return providerKeeper.SetLastProviderConsensusValSet(ctx, lastValidators)
}
