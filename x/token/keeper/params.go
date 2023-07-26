package keeper

import (
	"github.com/MANTRA-Finance/mantrachain/x/token/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams get all parameters as types.Params
func (ak Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	ak.paramstore.GetParamSet(ctx, &params)
	return
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
