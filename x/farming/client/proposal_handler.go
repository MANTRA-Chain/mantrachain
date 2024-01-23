package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/MANTRA-Finance/mantrachain/x/farming/client/cli"
)

// ProposalHandler is the public plan command handler.
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitPublicPlanProposal)
)
