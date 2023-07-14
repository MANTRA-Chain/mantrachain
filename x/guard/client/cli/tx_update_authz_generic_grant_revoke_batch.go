package cli

import (
	"strconv"

	"mantrachain/x/guard/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdUpdateAuthzGenericGrantRevokeBatch() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-authz-generic-grant-revoke-batch [grantee] [payload-json]",
		Short: "Broadcast message update_authz_generic_grant_revoke_batch",
		Long: "Updates authz grant revoke generic messages types. " +
			"[payload-json] is JSON encoded AuthzGrantRevokeMsgsTypes.",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argGrantee := args[0]
			argMsgsTypes := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Unmarshal payload
			var msgsTypes types.AuthzGrantRevokeMsgsTypes
			err = clientCtx.Codec.UnmarshalJSON([]byte(argMsgsTypes), &msgsTypes)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateAuthzGenericGrantRevokeBatch(
				clientCtx.GetFromAddress().String(),
				argGrantee,
				&msgsTypes,
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
