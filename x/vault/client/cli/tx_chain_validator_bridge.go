package cli

import (
	"github.com/LimeChain/mantrachain/x/vault/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdCreateChainValidatorBridge() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-chain-validator-bridge [chain] [validator] [bridge-id]",
		Short: "Create a new chain_validator_bridge",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ardChain := args[0]
			ardValidator := args[1]
			argBridgeId := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateChainValidatorBridge(
				clientCtx.GetFromAddress().String(),
				ardChain,
				ardValidator,
				argBridgeId,
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

func CmdUpdateChainValidatorBridge() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-chain-validator-bridge [chain] [validator] [bridge-id]",
		Short: "Update a chain_validator_bridge",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ardChain := args[0]
			ardValidator := args[1]
			argBridgeId := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateChainValidatorBridge(
				clientCtx.GetFromAddress().String(),
				ardChain,
				ardValidator,
				argBridgeId,
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

func CmdDeleteChainValidatorBridge() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-chain-validator-bridge [chain] [validator]",
		Short: "Delete a chain_validator_bridge",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ardChain := args[0]
			ardValidator := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteChainValidatorBridge(
				clientCtx.GetFromAddress().String(),
				ardChain,
				ardValidator,
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
