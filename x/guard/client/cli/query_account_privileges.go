package cli

import (
	"context"

	"github.com/LimeChain/mantrachain/x/guard/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

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

func CmdShowAccountPrivilegesMany() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-account-privileges-many [account ...]",
		Short: "shows many account_privileges",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetAccountPrivilegesManyRequest{
				Accounts: args,
			}

			res, err := queryClient.AccountPrivilegesMany(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
