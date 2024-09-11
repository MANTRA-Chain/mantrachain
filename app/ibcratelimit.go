package app

import (
	"cosmossdk.io/core/appmodule"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	ibcratelimit "github.com/osmosis-labs/osmosis/v26/x/ibc-rate-limit"
	ibcratelimitmodule "github.com/osmosis-labs/osmosis/v26/x/ibc-rate-limit/ibcratelimitmodule"
	ibcratelimittypes "github.com/osmosis-labs/osmosis/v26/x/ibc-rate-limit/types"
)

func (app *App) registerIBCRateLimit() {
	// Ensure the subspace is properly initialized
	app.ParamsKeeper.Subspace(ibcratelimittypes.ModuleName)

	// ChannelKeeper wrapper for rate limiting SendPacket(). The wasmKeeper needs to be added after it's created
	rateLimitingICS4Wrapper := ibcratelimit.NewICS4Middleware(
		app.HooksICS4Wrapper,
		&app.AccountKeeper,
		// wasm keeper we set later.
		nil,
		app.BankKeeper,
		app.GetSubspace(ibcratelimittypes.ModuleName),
	)
	app.RateLimitingICS4Wrapper = &rateLimitingICS4Wrapper

}

// RegisterTokenFactory registers the TokenFactory module with the given interface registry.
func RegisterIBCRateLimit(registry codectypes.InterfaceRegistry) map[string]appmodule.AppModule {
	modules := map[string]appmodule.AppModule{
		ibcratelimittypes.ModuleName: ibcratelimitmodule.AppModule{},
	}

	for name, m := range modules {
		module.CoreAppModuleBasicAdaptor(name, m).RegisterInterfaces(registry)
	}

	return modules
}
