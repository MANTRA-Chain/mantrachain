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

func CmdMint() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint [amount]",
		Short: "Mint a denom to an address. Must have admin authority to do so.",
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

			mintTo, err := cmd.Flags().GetString(FlagMintTo)
			if err != nil {
				return err
			}

			mintToStr := strings.TrimSpace(mintTo)
			if len(mintToStr) == 0 {
				mintTo = clientCtx.GetFromAddress().String()
			}

			msg := types.NewMsgMintTo(clientCtx.GetFromAddress().String(), amount, mintTo)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsMintTo)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
