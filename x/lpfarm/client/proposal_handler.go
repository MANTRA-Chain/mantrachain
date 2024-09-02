package client

import (
	"github.com/MANTRA-Finance/mantrachain/x/lpfarm/client/cli"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

// ProposalHandler is the public plan command handler.
var (
	ProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitFarmingPlanProposal)
)
