package cli

import (
	"fmt"

	"github.com/MANTRA-Finance/mantrachain/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	// Group bridge queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdListBridged())
	cmd.AddCommand(CmdShowBridged())
	// this line is used by starport scaffolding # 1

	return cmd
}
