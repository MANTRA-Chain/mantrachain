package app

import (
	ibcratelimittypes "github.com/osmosis-labs/osmosis/v26/x/ibc-rate-limit/types"
)

func (app *App) registerIBCRateLimit() {
	// Ensure the subspace is properly initialized
	app.ParamsKeeper.Subspace(ibcratelimittypes.ModuleName)
}
