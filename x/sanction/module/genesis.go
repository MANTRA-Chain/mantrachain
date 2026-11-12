package sanction

import (
	"github.com/MANTRA-Chain/mantrachain/v8/x/sanction/keeper"
	"github.com/MANTRA-Chain/mantrachain/v8/x/sanction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) error {
	for _, addr := range genState.BlacklistAccounts {
		if err := k.CheckNotModuleAccount(ctx, addr); err != nil {
			return err
		}
		if err := k.BlacklistAccounts.Set(ctx, addr); err != nil {
			return err
		}
	}
	if err := k.SetSanctioners(ctx, genState.Sanctioners); err != nil {
		return err
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

	genesis.BlacklistAccounts, err = iterKeys(ctx, k)
	if err != nil {
		return nil, err
	}

	genesis.Sanctioners, err = k.GetSanctioners(ctx)
	if err != nil {
		return nil, err
	}

	return genesis, nil
}

func iterKeys(ctx sdk.Context, k keeper.Keeper) ([]string, error) {
	iter, err := k.BlacklistAccounts.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	keys, err := iter.Keys()
	if err != nil {
		return nil, err
	}
	if keys == nil {
		keys = []string{}
	}
	return keys, nil
}
