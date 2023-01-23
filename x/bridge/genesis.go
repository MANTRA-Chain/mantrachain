package bridge

import (
	"github.com/LimeChain/mantrachain/x/bridge/keeper"
	"github.com/LimeChain/mantrachain/x/bridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	if _, err := sdk.AccAddressFromBech32(genState.Params.AdminAccount); err != nil {
		panic(sdkerrors.Wrap(types.ErrInvalidAdminAccount, "admin account param is invalid"))
	}
	for _, elem := range genState.TxHashList {
		k.SetTxHash(ctx, elem)
	}
	for _, elem := range genState.BridgeList {
		k.SetBridge(ctx, elem)
	}
	if genState.Cw20Contract != nil {
		k.SetCw20Contract(ctx, *genState.Cw20Contract)
	}
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.TxHashList = k.GetAllTxHash(ctx)
	genesis.BridgeList = k.GetAllBridge(ctx)
	cw20Contract, found := k.GetCw20Contract(ctx)
	if found {
		genesis.Cw20Contract = &cw20Contract
	}
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
