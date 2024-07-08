package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/MANTRA-Finance/mantrachain/x/bridge/types"
)

func CmdListBridged() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-bridged",
		Short: "list all bridged",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllBridgedRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.QueryBridgedAll(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowBridged() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-bridged [ethTxHash]",
		Short: "shows a bridged",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argEthTxHash := args[0]

			params := &types.QueryGetBridgedRequest{
				EthTxHash: argEthTxHash,
			}

			res, err := queryClient.QueryBridged(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}