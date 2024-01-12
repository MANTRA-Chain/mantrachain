package cli

import (
	"fmt"
	"strconv"

	"github.com/AumegaChain/aumega/x/txfees/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdCreateFeeToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-fee-token [denom] [pair-id]",
		Short: "Create a new fee_token",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get value arguments
			argDenom := args[0]

			pairId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("parse pair id: %w", err)
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateFeeToken(
				clientCtx.GetFromAddress().String(),
				argDenom,
				pairId,
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

func CmdUpdateFeeToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-fee-token [denom] [pair-id]",
		Short: "Update a fee_token",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get value arguments
			argDenom := args[0]

			pairId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("parse pair id: %w", err)
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateFeeToken(
				clientCtx.GetFromAddress().String(),
				argDenom,
				pairId,
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

func CmdDeleteFeeToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-fee-token [denom]",
		Short: "Delete a fee_token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argDenom := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteFeeToken(
				clientCtx.GetFromAddress().String(),
				argDenom,
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
