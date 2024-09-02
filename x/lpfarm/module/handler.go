package lpfarm

import (
	"cosmossdk.io/errors"
	"github.com/MANTRA-Finance/mantrachain/x/lpfarm/keeper"
	"github.com/MANTRA-Finance/mantrachain/x/lpfarm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

func NewFarmingPlanProposalHandler(k keeper.Keeper) govv1beta1.Handler {
	return func(ctx sdk.Context, content govv1beta1.Content) error {
		switch c := content.(type) {
		case *types.FarmingPlanProposal:
			return keeper.HandleFarmingPlanProposal(ctx, k, c)
		default:
			return errors.Wrapf(errorstypes.ErrUnknownRequest, "unrecognized lpfarm proposal content type: %T", c)
		}
	}
}
