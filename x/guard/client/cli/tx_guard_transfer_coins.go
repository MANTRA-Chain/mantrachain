package cli

import (
	"github.com/AumegaChain/aumega/x/guard/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdUpdateGuardTransferCoins() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-guard-transfer-coins [enabled]",
		Short: "Update guard_transfer_coins",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argEnabled, err := cast.ToBoolE(args[0])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateGuardTransferCoins(clientCtx.GetFromAddress().String(), argEnabled)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
