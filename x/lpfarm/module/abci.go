package lpfarm

import (
	"time"

	"github.com/MANTRA-Finance/mantrachain/x/lpfarm/keeper"
	"github.com/MANTRA-Finance/mantrachain/x/lpfarm/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	if err := k.TerminateEndedPlans(ctx); err != nil {
		panic(err)
	}
	if err := k.AllocateRewards(ctx); err != nil {
		panic(err)
	}
	k.SetLastBlockTime(ctx, ctx.BlockTime())
}
