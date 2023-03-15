package cli

import (
	"github.com/LimeChain/mantrachain/x/guard/types"
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
			// Get indexes
			index := []byte(args[0])

			// Get value arguments
			argPrivileges := []byte(args[1])

			argKind := args[1]

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
			argRequiredPrivilegesList := args[0]

			argKind := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Unmarshal payload
			var requiredPrivilegesList types.MsgRequiredPrivilegesList
			err = clientCtx.Codec.UnmarshalJSON([]byte(argRequiredPrivilegesList), &requiredPrivilegesList)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateRequiredPrivilegesBatch(
				clientCtx.GetFromAddress().String(),
				requiredPrivilegesList,
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
			argRequiredPrivilegesListGrouped := args[0]

			argKind := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Unmarshal payload
			var requiredPrivilegesListGrouped types.MsgRequiredPrivilegesListGrouped
			err = clientCtx.Codec.UnmarshalJSON([]byte(argRequiredPrivilegesListGrouped), &requiredPrivilegesListGrouped)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateRequiredPrivilegesGroupedBatch(
				clientCtx.GetFromAddress().String(),
				requiredPrivilegesListGrouped,
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
