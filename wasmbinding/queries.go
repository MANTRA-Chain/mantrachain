package wasmbinding

import (
	mdbkeeper "github.com/LimeChain/mantrachain/x/mdb/keeper"
)

type QueryPlugin struct {
	mdbkeeper *mdbkeeper.Keeper
}

// NewQueryPlugin returns a reference to a new QueryPlugin.
func NewQueryPlugin(mdbkeeper *mdbkeeper.Keeper) *QueryPlugin {
	return &QueryPlugin{
		mdbkeeper: mdbkeeper,
	}
}
