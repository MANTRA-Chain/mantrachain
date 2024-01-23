package marketmaker

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/MANTRA-Finance/mantrachain/x/marketmaker/keeper"
	"github.com/MANTRA-Finance/mantrachain/x/marketmaker/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgApplyMarketMaker:
			res, err := msgServer.ApplyMarketMaker(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgClaimIncentives:
			res, err := msgServer.ClaimIncentives(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}

// NewMarketMakerProposalHandler creates a governance handler to manage new proposal types.
// It enables MarketMakerProposal to propose market maker inclusion / exclusion / rejection / distribution.
func NewMarketMakerProposalHandler(k keeper.Keeper) govv1beta1.Handler {
	return func(ctx sdk.Context, content govv1beta1.Content) error {
		switch c := content.(type) {
		case *types.MarketMakerProposal:
			return keeper.HandleMarketMakerProposal(ctx, k, c)

		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized market maker proposal content type: %T", c)
		}
	}
}
