package cli

import (
	"encoding/base64"

	"github.com/MANTRA-Finance/mantrachain/x/guard/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdUpdateAccountPrivileges() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-account-privileges [account] [privileges]",
		Short: "Update a account_privileges",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			indexAcc := args[0]

			argPrivileges, err := base64.StdEncoding.DecodeString(args[1])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateAccountPrivileges(
				clientCtx.GetFromAddress().String(),
				indexAcc,
				argPrivileges,
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
