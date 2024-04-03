package airdrop

import (
	"github.com/MANTRA-Finance/mantrachain/x/airdrop/keeper"
	"github.com/MANTRA-Finance/mantrachain/x/airdrop/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the campaign
	for _, elem := range genState.CampaignList {
		k.SetCampaign(ctx, elem)
	}
	// Set all the campaign
	for _, elem := range genState.ClaimedList {
		k.SetClaimed(ctx, elem)
	}
	k.SetLastCampaignId(ctx, genState.LastCampaignId)
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	lastCampaignId, _ := k.GetLastCampaignId(ctx)
	genesis.LastCampaignId = lastCampaignId
	genesis.CampaignList = k.GetAllCampaign(ctx)
	genesis.ClaimedList = k.GetAllClaimed(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}