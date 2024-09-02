package client

import (
	"github.com/MANTRA-Finance/mantrachain/x/marketmaker/client/cli"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

// ProposalHandler is the market maker proposal command handler.
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitMarketMakerProposal)
)
