package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/LimeChain/mantrachain/x/marketmaker/client/cli"
)

// ProposalHandler is the market maker proposal command handler.
// Note that rest.ProposalRESTHandler will be deprecated in the future.
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitMarketMakerProposal)
)
