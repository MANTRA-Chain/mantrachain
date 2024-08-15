package v3

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
)

// CreateUpgradeHandler creates an SDK upgrade handler for v3.0.0
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	govkeeper govkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(goCtx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(goCtx)
		// set the min deposit for expedited proposals to 10x the value of the normal min deposit
		govParams, err := govkeeper.Params.Get(ctx)
		if err != nil {
			return nil, err
		}
		expeditedMinDeposit := govParams.GetMinDeposit()
		for i := range expeditedMinDeposit {
			expeditedMinDeposit[i].Amount = expeditedMinDeposit[i].Amount.MulRaw(10)
		}
		govParams.ExpeditedMinDeposit = expeditedMinDeposit
		if err := govkeeper.Params.Set(ctx, govParams); err != nil {
			return nil, err
		}

		return mm.RunMigrations(goCtx, configurator, fromVM)
	}
}
