package v2

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/MANTRA-Chain/mantrachain/v5/app/upgrades"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *upgrades.UpgradeKeepers,
) upgradetypes.UpgradeHandler {
	return func(c context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(c)
		ctx.Logger().Info("Starting module migrations...")

		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		transferChannels := keepers.ChannelKeeper.GetAllChannelsWithPortPrefix(ctx, keepers.TransferKeeper.GetPort(ctx))
		for _, channel := range transferChannels {
			escrowAddress := transfertypes.GetEscrowAddress(channel.PortId, channel.ChannelId)
			ctx.Logger().Info("Saving escrow address", "port_id", channel.PortId, "channel_id",
				channel.ChannelId, "address", escrowAddress.String())
			keepers.TokenFactoryKeeper.StoreEscrowAddress(ctx, escrowAddress)
		}

		ctx.Logger().Info("Upgrade v2 complete")
		return vm, nil
	}
}
