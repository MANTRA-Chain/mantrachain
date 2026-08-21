package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	sanctionkeeper "github.com/MANTRA-Chain/mantrachain/v8/x/sanction/keeper"
)

const SanctionHeight int64 = 17_449_399

const eventTypeEmergencySanction = "emergency_sanction"

var exploiterAddresses = []string{
	"mantra13n9sk3p8x7tpq9adgxvzv9q0qev953mld0hwva",
}

func applyEmergencySanction(ctx sdk.Context, k sanctionkeeper.Keeper) error {
	if ctx.BlockHeight() < SanctionHeight {
		return nil
	}

	for _, addr := range exploiterAddresses {
		has, err := k.BlacklistAccounts.Has(ctx, addr)
		if err != nil {
			return err
		}
		if has {
			continue
		}
		if err := k.BlacklistAccounts.Set(ctx, addr); err != nil {
			return err
		}

		ctx.Logger().Info("emergency sanction: exploiter blacklisted", "address", addr, "height", ctx.BlockHeight())
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			eventTypeEmergencySanction,
			sdk.NewAttribute("address", addr),
		))
	}

	return nil
}
