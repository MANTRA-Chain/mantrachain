package txfees

import (
	"strings"

	"github.com/MANTRA-Finance/aumega/x/txfees/keeper"
	"github.com/MANTRA-Finance/aumega/x/txfees/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the feeToken
	for _, elem := range genState.FeeTokenList {
		k.SetFeeToken(ctx, elem)
	}
	if strings.TrimSpace(genState.Params.BaseDenom) == "" {
		panic(errors.Wrap(types.ErrInvalidBaseDenomParam, "base denom param should not be empty"))
	}
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.FeeTokenList = k.GetAllFeeToken(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
