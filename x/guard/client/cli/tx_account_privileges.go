package cli

import (
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
			// Get indexes
			indexAcc := args[0]

			// Get value arguments
			argPrivileges := []byte(args[1])

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

func CmdUpdateAccountPrivilegesBatch() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-account-privileges-batch [payload-json]",
		Short: "Update account_privileges in a batch",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argAccountsPrivileges := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Unmarshal payload
			var accountsPrivileges types.MsgAccountsPrivileges
			err = clientCtx.Codec.UnmarshalJSON([]byte(argAccountsPrivileges), &accountsPrivileges)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateAccountPrivilegesBatch(
				clientCtx.GetFromAddress().String(),
				accountsPrivileges,
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

func CmdUpdateAccountPrivilegesGroupedBatch() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-account-privileges-grouped-batch [payload-json]",
		Short: "Update account_privileges_grouped in a batch",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argAccountsPrivilegesGrouped := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Unmarshal payload
			var accountsPrivilegesGrouped types.MsgAccountsPrivilegesGrouped
			err = clientCtx.Codec.UnmarshalJSON([]byte(argAccountsPrivilegesGrouped), &accountsPrivilegesGrouped)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateAccountPrivilegesGroupedBatch(
				clientCtx.GetFromAddress().String(),
				accountsPrivilegesGrouped,
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
