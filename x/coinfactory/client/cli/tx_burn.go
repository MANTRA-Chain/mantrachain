package cli

import (
	"strings"

	"github.com/MANTRA-Finance/mantrachain/x/coinfactory/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

func CmdBurn() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn [amount]",
		Short: "Burn tokens from an address. Must have admin authority to do so.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			amount, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			burnFrom, err := cmd.Flags().GetString(FlagBurnFrom)
			if err != nil {
				return err
			}

			burnFromStr := strings.TrimSpace(burnFrom)
			if len(burnFromStr) == 0 {
				burnFrom = clientCtx.GetFromAddress().String()
			}

			msg := types.NewMsgBurnFrom(clientCtx.GetFromAddress().String(), amount, burnFrom)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsBurnFrom)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
