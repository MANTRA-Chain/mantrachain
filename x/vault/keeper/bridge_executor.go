package keeper

import (
	bridgetypes "github.com/LimeChain/mantrachain/x/bridge/types"
	"github.com/LimeChain/mantrachain/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BridgeExecutor struct {
	ctx sdk.Context
	bk  types.BridgeKeeper
}

func NewBridgeExecutor(ctx sdk.Context, bk types.BridgeKeeper) *BridgeExecutor {
	return &BridgeExecutor{
		ctx: ctx,
		bk:  bk,
	}
}

func (c *BridgeExecutor) GetBridge(bridgeCreator sdk.AccAddress, bridgeId string) (bridgetypes.Bridge, bool) {
	index := bridgetypes.GetBridgeIndex(bridgeCreator, bridgeId)

	return c.bk.GetBridge(
		c.ctx,
		sdk.AccAddress(bridgeCreator),
		index,
	)
}
