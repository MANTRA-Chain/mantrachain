package vault

import (
	"strings"

	"github.com/LimeChain/mantrachain/x/vault/keeper"
	"github.com/LimeChain/mantrachain/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	if strings.TrimSpace(genState.Params.StakingValidatorAddress) == "" {
		panic(sdkerrors.Wrap(types.ErrInvalidStakingValidatorAddress, "staking validator address param should not be empty"))
	}
	if _, err := sdk.AccAddressFromBech32(genState.Params.AdminAccount); err != nil {
		panic(sdkerrors.Wrap(types.ErrInvalidAdminAccount, "admin account param is invalid"))
	}
	for _, elem := range genState.ChainValidatorBridgeList {
		k.SetChainValidatorBridge(ctx, elem.Chain, elem.Validator, elem)
	}
	for _, elem := range genState.EpochList {
		k.SetEpoch(ctx, &elem.StakingChain, &elem.StakingValidator, elem.Id, elem)
	}
	for _, elem := range genState.LastEpochBlockList {
		k.SetLastEpochBlock(ctx, &elem.StakingChain, &elem.StakingValidator, elem)
	}
	for _, elem := range genState.NftStakeList {
		k.SetNftStake(ctx, elem)
	}
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.ChainValidatorBridgeList = k.GetAllChainValidatorBridge(ctx, nil)
	genesis.EpochList = k.GetAllEpoch(ctx, nil, nil)
	genesis.LastEpochBlockList = k.GetAllLastEpochBlock(ctx, nil, nil)
	genesis.NftStakeList = k.GetAllNftStake(ctx, nil, nil)

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
