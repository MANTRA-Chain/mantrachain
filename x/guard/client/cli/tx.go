package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"mantrachain/x/guard/types"

	"github.com/cosmos/cosmos-sdk/client"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
)

const (
	flagPacketTimeoutTimestamp = "packet-timeout-timestamp"
	listSeparator              = ","
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdUpdateAccountPrivileges())
	cmd.AddCommand(CmdUpdateAccountPrivilegesBatch())
	cmd.AddCommand(CmdUpdateAccountPrivilegesGroupedBatch())
	cmd.AddCommand(CmdUpdateGuardTransferCoins())
	cmd.AddCommand(CmdUpdateRequiredPrivileges())
	cmd.AddCommand(CmdUpdateRequiredPrivilegesBatch())
	cmd.AddCommand(CmdUpdateRequiredPrivilegesGroupedBatch())
	cmd.AddCommand(CmdUpdateLocked())
	cmd.AddCommand(CmdUpdateAuthzGenericGrantRevokeBatch())
	// this line is used by starport scaffolding # 1

	return cmd
}
