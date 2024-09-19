package tax

import (
	"github.com/MANTRA-Chain/mantrachain/x/tax/keeper"
	"github.com/MANTRA-Chain/mantrachain/x/tax/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker sets the proposer for determining distribution during endblock
// and distribute rewards for the previous block.
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyBeginBlocker)

	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}
	// if the mca tax is zero, no need to continue
	if params.McaTax.IsZero() {
		return nil
	}

	// only allocate rewards if the block height is greater than 1
	if ctx.BlockHeight() > 1 {
		McaAddress, err := sdk.AccAddressFromBech32(params.McaAddress)
		if err != nil {
			return err
		}
		if err := k.AllocateMcaTax(ctx, params.McaTax, McaAddress); err != nil {
			return err
		}
	}

	return nil
}
