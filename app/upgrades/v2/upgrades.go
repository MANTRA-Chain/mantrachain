package v2

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	cbtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	consensuskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// NOTE: we're only keeping this logic for the upgrade tests
// This is not the original upgrade logic.
// Look into the previous version if want to know what the upgrade logic was
// CreateUpgradeHandler creates an SDK upgrade handler for v15.0.0
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	cdc codec.Codec,
	baseAppLegacySS paramstypes.Subspace,
	ps baseapp.ParamStore,
	upgradeKey *storetypes.KVStoreKey,
	authority string,
) upgradetypes.UpgradeHandler {
	return func(goCtx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(goCtx)

		baseapp.MigrateParams(ctx, baseAppLegacySS, ps)
		consensusParams := baseapp.GetConsensusParams(ctx, baseAppLegacySS)
		// make sure the consensus params are set
		if consensusParams.Block == nil || consensusParams.Evidence == nil || consensusParams.Validator == nil {
			defaultParams := cbtypes.DefaultConsensusParams().ToProto()
			ps.Set(ctx, defaultParams)
		}

		storesvc := runtime.NewKVStoreService(upgradeKey)
		consensuskeeper := consensuskeeper.NewKeeper(
			cdc,
			storesvc,
			authority,
			runtime.EventService{},
		)

		params, err := consensuskeeper.ParamsStore.Get(ctx)
		if err != nil {
			return nil, err
		}

		err = ps.Set(ctx, params)
		if err != nil {
			return nil, err
		}

		return mm.RunMigrations(goCtx, configurator, fromVM)
	}
}
