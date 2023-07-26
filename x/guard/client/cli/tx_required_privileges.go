package cli

import (
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
			argPrivileges := []byte(args[1])
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

func CmdUpdateRequiredPrivilegesBatch() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-required-privileges-batch [payload-json] [kind]",
		Short: "Update required_privileges in a batch",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argRequiredPrivileges := args[0]
			argKind := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Unmarshal payload
			var requiredPrivileges types.MsgRequiredPrivileges
			err = clientCtx.Codec.UnmarshalJSON([]byte(argRequiredPrivileges), &requiredPrivileges)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateRequiredPrivilegesBatch(
				clientCtx.GetFromAddress().String(),
				requiredPrivileges,
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

func CmdUpdateRequiredPrivilegesGroupedBatch() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-required-privileges-grouped-batch [payload-json] [kind]",
		Short: "Update required_privileges_grouped in a batch",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argRequiredPrivilegesGrouped := args[0]
			argKind := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Unmarshal payload
			var requiredPrivilegesGrouped types.MsgRequiredPrivilegesGrouped
			err = clientCtx.Codec.UnmarshalJSON([]byte(argRequiredPrivilegesGrouped), &requiredPrivilegesGrouped)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateRequiredPrivilegesGroupedBatch(
				clientCtx.GetFromAddress().String(),
				requiredPrivilegesGrouped,
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
