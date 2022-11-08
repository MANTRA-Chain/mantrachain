package cli

import (
	"context"

	"github.com/LimeChain/mantrachain/x/vault/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdGetChainValidatorBridge() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chain-validator-bridge [chain] [validator]",
		Short: "shows a chain_validator_bridge",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argChain := args[0]
			argValidator := args[1]

			params := &types.QueryGetChainValidatorBridgeRequest{
				Chain:     argChain,
				Validator: argValidator,
			}

			res, err := queryClient.ChainValidatorBridge(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
