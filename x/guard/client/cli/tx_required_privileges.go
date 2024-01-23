package cli

import (
	"encoding/base64"

	"github.com/MANTRA-Finance/mantrachain/x/guard/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdUpdateRequiredPrivileges() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-required-privileges [index] [privileges] [kind]",
		Short: "Update a required_privileges",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			index := []byte(args[0])

			argPrivileges, err := base64.StdEncoding.DecodeString(args[1])
			if err != nil {
				return err
			}
			argKind := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateRequiredPrivileges(
				clientCtx.GetFromAddress().String(),
				index,
				argPrivileges,
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
