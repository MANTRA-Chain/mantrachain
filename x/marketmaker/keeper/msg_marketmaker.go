package keeper

import (
	"context"

	"github.com/MANTRA-Finance/mantrachain/x/marketmaker/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ApplyMarketMaker defines a method for apply to be market maker.
func (k msgServer) ApplyMarketMaker(goCtx context.Context, msg *types.MsgApplyMarketMaker) (*types.MsgApplyMarketMakerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.Keeper.ApplyMarketMaker(ctx, msg.GetAccAddress(), msg.PairIds); err != nil {
		return nil, err
	}

	return &types.MsgApplyMarketMakerResponse{}, nil
}

// ClaimIncentives defines a method for claim all claimable incentives of the market maker.
func (k msgServer) ClaimIncentives(goCtx context.Context, msg *types.MsgClaimIncentives) (*types.MsgClaimIncentivesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.Keeper.ClaimIncentives(ctx, msg.GetAccAddress()); err != nil {
		return nil, err
	}

	return &types.MsgClaimIncentivesResponse{}, nil
}
