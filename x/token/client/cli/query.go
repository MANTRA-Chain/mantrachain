package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/MANTRA-Finance/aumega/x/token/types"
	"github.com/cosmos/cosmos-sdk/client"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group token queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdGetNftCollection())
	cmd.AddCommand(CmdGetNft())
	cmd.AddCommand(CmdGetNftApproved())
	cmd.AddCommand(CmdGetIsApprovedForAllNfts())
	cmd.AddCommand(CmdGetNftCollectionsByCreator())
	cmd.AddCommand(CmdGetAllCollectionNfts())
	cmd.AddCommand(CmdGetAllNftCollections())

	cmd.AddCommand(CmdShowSoulBondedNftsCollection())
	cmd.AddCommand(CmdShowRestrictedNftsCollection())
	cmd.AddCommand(CmdShowOpenedNftsCollection())
	cmd.AddCommand(CmdShowNftCollectionOwner())
	// this line is used by starport scaffolding # 1

	return cmd
}
