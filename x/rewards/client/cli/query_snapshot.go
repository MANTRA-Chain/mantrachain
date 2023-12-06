package cli

import (
    "strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
    "github.com/MANTRA-Finance/mantrachain/x/rewards/types"
)

func CmdListSnapshot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-snapshot",
		Short: "list all snapshot",
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

            params := &types.QueryAllSnapshotRequest{
                Pagination: pageReq,
            }

            res, err := queryClient.SnapshotAll(cmd.Context(), params)
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

func CmdShowSnapshot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-snapshot [id]",
		Short: "shows a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
            clientCtx, err := client.GetClientQueryContext(cmd)
            if err != nil {
                return err
            }

            queryClient := types.NewQueryClient(clientCtx)

            id, err := strconv.ParseUint(args[0], 10, 64)
            if err != nil {
                return err
            }

            params := &types.QueryGetSnapshotRequest{
                Id: id,
            }

            res, err := queryClient.Snapshot(cmd.Context(), params)
            if err != nil {
                return err
            }

            return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

    return cmd
}
