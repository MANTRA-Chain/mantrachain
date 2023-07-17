package guard

import (
	"strings"

	"mantrachain/x/guard/keeper"
	"mantrachain/x/guard/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k *keeper.Keeper, genState types.GenesisState) {
	if _, err := sdk.AccAddressFromBech32(genState.Params.AccountPrivilegesTokenCollectionCreator); err != nil {
		panic(err)
	}
	if strings.TrimSpace(genState.Params.AccountPrivilegesTokenCollectionId) == "" {
		panic(errors.Wrap(types.ErrInvalidAccountPrivilegesTokenCollectionIdParam, "account privileges token collection id param should not be empty"))
	}
	// Set all the accountPrivileges
	for _, elem := range genState.AccountPrivilegesList {
		k.SetAccountPrivileges(ctx, elem.Account, elem.Privileges)
	}
	// Set if defined
	if genState.GuardTransferCoins != nil {
		k.SetGuardTransferCoins(ctx)
	}
	// Set all the requiredPrivileges
	for _, elem := range genState.RequiredPrivilegesList {
		kind, err := types.ParseRequiredPrivilegesKind(elem.Kind)
		if err != nil {
			panic("kind is invalid")
		}
		k.SetRequiredPrivileges(ctx, elem.Index, kind, elem.Privileges)
	}
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k *keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.AccountPrivilegesList = k.GetAllAccountPrivileges(ctx)
	// Get all guardTransfer
	if k.HasGuardTransferCoins(ctx) {
		genesis.GuardTransferCoins = types.Placeholder
	}
	genesis.RequiredPrivilegesList = k.GetAllRequiredPrivileges(ctx, types.RequiredPrivilegesCoin)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
