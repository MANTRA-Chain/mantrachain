package cli

import (
	"context"

	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdListRequiredPrivileges() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-required-privileges [kind]",
		Short: "list all required_privileges",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			argKind, err := types.ParseRequiredPrivilegesKind(args[0])
			if err != nil {
				return err
			}

			params := &types.QueryAllRequiredPrivilegesRequest{
				Pagination: pageReq,
				Kind:       argKind.String(),
			}

			res, err := queryClient.RequiredPrivilegesAll(context.Background(), params)
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

func CmdShowRequiredPrivileges() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-required-privileges [index] [kind]",
		Short: "shows a required_privileges",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argIndex := []byte(args[0])
			argKind, err := types.ParseRequiredPrivilegesKind(args[1])
			if err != nil {
				return err
			}

			params := &types.QueryGetRequiredPrivilegesRequest{
				Index: argIndex,
				Kind:  argKind.String(),
			}

			res, err := queryClient.RequiredPrivileges(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
