package cli

import (
	"fmt"
	// "strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group guard queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryParams())
	cmd.AddCommand(CmdShowAccountPrivileges())
	cmd.AddCommand(CmdListAccountPrivileges())
	cmd.AddCommand(CmdShowGuardTransferCoins())
	cmd.AddCommand(CmdListRequiredPrivileges())
	cmd.AddCommand(CmdShowRequiredPrivileges())
	cmd.AddCommand(CmdListLocked())
	cmd.AddCommand(CmdShowLocked())
	// this line is used by starport scaffolding # 1

	return cmd
}
