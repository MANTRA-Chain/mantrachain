package wasmbinding

import (
	tokenkeeper "github.com/LimeChain/mantrachain/x/token/keeper"
)

type QueryPlugin struct {
	tokenkeeper *tokenkeeper.Keeper
}

// NewQueryPlugin returns a reference to a new QueryPlugin.
func NewQueryPlugin(tokenkeeper *tokenkeeper.Keeper) *QueryPlugin {
	return &QueryPlugin{
		tokenkeeper: tokenkeeper,
	}
}
