package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/MANTRA-Finance/mantrachain/x/marketmaker/client/cli"
)

// ProposalHandler is the market maker proposal command handler.
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitMarketMakerProposal)
)
