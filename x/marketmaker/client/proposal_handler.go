package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"mantrachain/x/marketmaker/client/cli"
)

// ProposalHandler is the market maker proposal command handler.
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitMarketMakerProposal)
)
