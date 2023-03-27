package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/MANTRA-Finance/mantrachain/x/lpfarm/client/cli"
)

// ProposalHandler is the public plan command handler.
// Note that rest.ProposalRESTHandler will be deprecated in the future.
var (
	ProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitFarmingPlanProposal)
)
