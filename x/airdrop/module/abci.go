package airdrop

import (
	"time"

	"github.com/MANTRA-Finance/mantrachain/x/airdrop/keeper"
	"github.com/MANTRA-Finance/mantrachain/x/airdrop/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)
	if err := k.TerminateEndedCampaigns(ctx); err != nil {
		panic(err)
	}
}
