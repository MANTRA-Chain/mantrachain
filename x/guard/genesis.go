package guard

import (
	"strings"

	"github.com/LimeChain/mantrachain/x/guard/keeper"
	"github.com/LimeChain/mantrachain/x/guard/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	if _, err := sdk.AccAddressFromBech32(genState.Params.TokenCollectionCreator); err != nil {
		panic(sdkerrors.Wrap(types.ErrInvalidTokenCollectionCreatorParam, "token collection creator param is invalid"))
	}
	if strings.TrimSpace(genState.Params.TokenCollectionId) == "" {
		panic(sdkerrors.Wrap(types.ErrInvalidTokenCollectionIdParam, "token collection id param should not be empty"))
	}
	// Set all the accPerm
	for _, elem := range genState.AccPermList {
		k.SetAccPerm(ctx, elem)
	}
	// Set if defined
	if genState.GuardTransfer == nil {
		panic(sdkerrors.Wrap(types.ErrInvalidGuardTransfer, "guard transfer is invalid"))
	}
	k.SetGuardTransfer(ctx, *genState.GuardTransfer)
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.AccPermList = k.GetAllAccPerm(ctx)
	// Get all guardTransfer
	guardTransfer, found := k.GetGuardTransfer(ctx)
	if found {
		genesis.GuardTransfer = &guardTransfer
	}
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
