package cli

import (
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
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllBridgedRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.BridgedAll(cmd.Context(), params)
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
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			argEthTxHash := args[0]

			params := &types.QueryGetBridgedRequest{
				EthTxHash: argEthTxHash,
			}

			res, err := queryClient.Bridged(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
