package cli

import (
	"fmt"
	// "strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/LimeChain/mantrachain/x/token/types"
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
	cmd.AddCommand(CmdGetNftOwner())
	cmd.AddCommand(CmdGetNftBalance())
	cmd.AddCommand(CmdGetNftCollectionSupply())
	cmd.AddCommand(CmdGetNftCollectionsByCreator())
	cmd.AddCommand(CmdGetAllCollectionNfts())
	cmd.AddCommand(CmdGetAllNftCollections())
	cmd.AddCommand(CmdGetCollectionNftsByOwner())

	cmd.AddCommand(CmdListSoulBondedNftsCollection())
	cmd.AddCommand(CmdShowSoulBondedNftsCollection())
// this line is used by starport scaffolding # 1

	return cmd
}
