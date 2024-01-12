package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/AumegaChain/aumega/x/lpfarm/client/cli"
)

// ProposalHandler is the public plan command handler.
var (
	ProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitFarmingPlanProposal)
)
