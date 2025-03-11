package sanction

import (
	"github.com/MANTRA-Chain/mantrachain/v4/x/sanction/keeper"
	"github.com/MANTRA-Chain/mantrachain/v4/x/sanction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) error {
	// this line is used by starport scaffolding # genesis/module/init
	for _, blacklist_account := range genState.BlacklistAccounts {
		if err := k.BlacklistAccounts.Set(ctx, blacklist_account); err != nil {
			return err
		}
	}
	return k.Params.Set(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) (*types.GenesisState, error) {
	var err error

	genesis := types.DefaultGenesis()
	genesis.Params, err = k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	blacklistAccounts := make([]string, 0)
	iter, err := k.BlacklistAccounts.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		blacklistAccount, err := iter.Key()
		if err != nil {
			return nil, err
		}
		blacklistAccounts = append(blacklistAccounts, blacklistAccount)
	}
	genesis.BlacklistAccounts = blacklistAccounts

	// this line is used by starport scaffolding # genesis/module/export

	return genesis, nil
}
