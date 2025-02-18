package cli

import (
	"fmt"

	"github.com/MANTRA-Chain/mantrachain/v2/x/sanction/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewAddBlacklistAccountCmd(),
		NewRemoveBlacklistAccountCmd(),
	)

	return cmd
}

// NewAddBlacklistAccountCmd broadcast MsgAddBlacklistAccount
func NewAddBlacklistAccountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-blacklist-account [blacklist-account] [flags]",
		Short: "add an account to the blacklist",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf, err := tx.NewFactoryCLI(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			msg := types.NewMsgAddBlacklistAccount(
				clientCtx.GetFromAddress().String(),
				args[0],
			)

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf.WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// NewRemoveBlacklistAccountCmd broadcast MsgRemoveBlacklistAccount
func NewRemoveBlacklistAccountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-blacklist-account [blacklist-account] [flags]",
		Short: "remove an account from the blacklist",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf, err := tx.NewFactoryCLI(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			msg := types.NewMsgRemoveBlacklistAccount(
				clientCtx.GetFromAddress().String(),
				args[0],
			)

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf.WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
