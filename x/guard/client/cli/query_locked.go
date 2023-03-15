package cli

import (
	"context"

	"github.com/LimeChain/mantrachain/x/guard/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdListLocked() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-locked [kind]",
		Short: "list all locked",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			argKind, err := types.ParseLockedKind(args[0])
			if err != nil {
				return err
			}

			params := &types.QueryAllLockedRequest{
				Pagination: pageReq,
				Kind:       argKind.String(),
			}

			res, err := queryClient.LockedAll(context.Background(), params)
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

func CmdShowLocked() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-locked [index] [kind]",
		Short: "shows a locked",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argIndex := []byte(args[0])
			argKind, err := types.ParseLockedKind(args[1])
			if err != nil {
				return err
			}

			params := &types.QueryGetLockedRequest{
				Index: argIndex,
				Kind:  argKind.String(),
			}

			res, err := queryClient.Locked(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
