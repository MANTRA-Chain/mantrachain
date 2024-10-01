package v1

// TODO: Add any additional upgrades here

/*
func CreateUpgradeHandler(mm upgrades.ModuleManager,
	configurator module.Configurator,
	keepers *upgrades.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		sdkCtx.Logger().Info("Starting module migrations...")

		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		err = ConfigureFeeMarketModule(sdkCtx, keepers)
		if err != nil {
			return vm, err
		}
		return vm, nil
	}
}

*/
