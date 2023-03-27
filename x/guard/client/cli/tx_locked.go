package cli

import (
	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdUpdateLocked() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-locked [index] [locked] [kind]",
		Short: "Update a locked",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			index := []byte(args[0])
			argLocked, err := cast.ToBoolE(args[1])
			if err != nil {
				return err
			}
			argKind := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateLocked(
				clientCtx.GetFromAddress().String(),
				index,
				argLocked,
				argKind,
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
