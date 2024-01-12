package farming

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/AumegaChain/aumega/x/farming/keeper"
	"github.com/AumegaChain/aumega/x/farming/types"
)

func NewHandler(_ keeper.Keeper) sdk.Handler {
	return func(_ sdk.Context, _ sdk.Msg) (*sdk.Result, error) {
		return nil, types.ErrModuleDisabled
	}
}

// NewPublicPlanProposalHandler creates a governance handler to manage new proposal types.
// It enables PublicPlanProposal to propose a plan creation / modification / deletion.
func NewPublicPlanProposalHandler(_ keeper.Keeper) govv1beta1.Handler {
	return func(_ sdk.Context, _ govv1beta1.Content) error {
		return types.ErrModuleDisabled
	}
}
