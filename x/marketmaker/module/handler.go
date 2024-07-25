package marketmaker

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/MANTRA-Finance/mantrachain/x/marketmaker/keeper"
	"github.com/MANTRA-Finance/mantrachain/x/marketmaker/types"
)

// NewMarketMakerProposalHandler creates a governance handler to manage new proposal types.
// It enables MarketMakerProposal to propose market maker inclusion / exclusion / rejection / distribution.
func NewMarketMakerProposalHandler(k keeper.Keeper) govv1beta1.Handler {
	return func(ctx sdk.Context, content govv1beta1.Content) error {
		switch c := content.(type) {
		case *types.MarketMakerProposal:
			return keeper.HandleMarketMakerProposal(ctx, k, c)

		default:
			return errors.Wrapf(errorstypes.ErrUnknownRequest, "unrecognized market maker proposal content type: %T", c)
		}
	}
}
