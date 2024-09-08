package xfeemarket

import (
	"cosmossdk.io/math"
	"github.com/MANTRA-Chain/mantrachain/x/xfeemarket/keeper"
	"github.com/MANTRA-Chain/mantrachain/x/xfeemarket/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) error {
	// this line is used by starport scaffolding # genesis/module/init
	for _, denomMultiplier := range genState.GetDenomMultipliers() {
		multiplier := math.LegacyMustNewDecFromStr(denomMultiplier.Multiplier)
		if err := k.DenomMultipliers.Set(ctx, denomMultiplier.Denom, sdk.DecProto{Dec: multiplier}); err != nil {
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

	denomMultipliers := make([]types.DenomMultiplier, 0)
	iter, err := k.DenomMultipliers.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		denom, err := iter.Key()
		if err != nil {
			return nil, err
		}
		multiplierDec, err := iter.Value()
		if err != nil {
			return nil, err
		}
		denomMultipliers = append(denomMultipliers, types.DenomMultiplier{
			Denom:      denom,
			Multiplier: multiplierDec.Dec.String(),
		})
	}
	genesis.DenomMultipliers = denomMultipliers

	// this line is used by starport scaffolding # genesis/module/export

	return genesis, nil
}
