package cli

import (
	"context"

	"github.com/MANTRA-Finance/mantrachain/x/txfees/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdListFeeToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-fee-token",
		Short: "list all fee_token",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllFeeTokenRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.FeeTokenAll(context.Background(), params)
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

func CmdShowFeeToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-fee-token [denom]",
		Short: "shows a fee_token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argDenom := args[0]

			params := &types.QueryGetFeeTokenRequest{
				Denom: argDenom,
			}

			res, err := queryClient.FeeToken(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
