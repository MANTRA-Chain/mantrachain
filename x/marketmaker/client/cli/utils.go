package cli

import (
	"os"

	"github.com/MANTRA-Finance/mantrachain/x/marketmaker/types"
	"github.com/cosmos/cosmos-sdk/codec"
)

// ParseMarketMakerProposal reads and parses a MarketMakerProposal from a file.
func ParseMarketMakerProposal(cdc codec.JSONCodec, proposalFile string) (types.MarketMakerProposal, error) {
	proposal := types.MarketMakerProposal{}

	contents, err := os.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err = cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
