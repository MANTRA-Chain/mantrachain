package cli

import (
	"context"

	"github.com/AumegaChain/aumega/x/guard/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdListAccountPrivileges() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-account-privileges",
		Short: "list all account_privileges",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllAccountPrivilegesRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.AccountPrivilegesAll(context.Background(), params)
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

func CmdShowAccountPrivileges() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-account-privileges [account]",
		Short: "shows a account_privileges",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argAccount := args[0]

			params := &types.QueryGetAccountPrivilegesRequest{
				Account: argAccount,
			}

			res, err := queryClient.AccountPrivileges(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
