package cli

import (
	"fmt"
	"strconv"

	"github.com/LimeChain/mantrachain/x/vault/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdStartEpoch() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-epoch [block-start] [reward] [chain] [validator]",
		Short: "Broadcast message start-epoch",
		Long:  "Starts new epoch. ",
		Example: fmt.Sprintf(
			"$ %s tx vault start-epoch <block-start> <reward> <chain> <validator> "+
				"--from=<from> "+
				"--chain-id=<chain-id> ",
			version.AppName,
		),
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argBlockStart := args[0]
			argReward := args[1]
			argChain := args[2]
			argValidator := args[3]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			blockStart, err := strconv.ParseInt(argBlockStart, 10, 64)

			if err != nil {
				return err
			}

			msg := types.NewMsgStartEpoch(
				clientCtx.GetFromAddress().String(),
				blockStart,
				argReward,
				argChain,
				argValidator,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
