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
)

// CreateUpgradeHandler creates an SDK upgrade handler for v2.0.0
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	cdc codec.Codec,
	ps baseapp.ParamStore,
	upgradeKey *storetypes.KVStoreKey,
	authority string,
) upgradetypes.UpgradeHandler {
	return func(goCtx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(goCtx)
		// Fix the consensus params: https://github.com/cosmos/cosmos-sdk/issues/18733#issuecomment-1854611049
		defaultParams := cbtypes.DefaultConsensusParams().ToProto()

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

		if params.Block == nil {
			params.Block = defaultParams.Block
		}

		if params.Evidence == nil {
			params.Evidence = defaultParams.Evidence
		}

		if params.Validator == nil {
			params.Validator = defaultParams.Validator
		}

		if params.Version == nil {
			params.Version = defaultParams.Version
		}

		if params.Abci == nil {
			params.Abci = defaultParams.Abci
		}

		err = ps.Set(ctx, params)
		if err != nil {
			return nil, err
		}

		return mm.RunMigrations(goCtx, configurator, fromVM)
	}
}
