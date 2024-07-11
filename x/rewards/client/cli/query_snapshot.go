package cli

import (
	"context"
	"strconv"

	"github.com/MANTRA-Finance/mantrachain/x/rewards/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdListSnapshot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-snapshot",
		Short: "list all snapshot",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllSnapshotRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.SnapshotAll(context.Background(), params)
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
		Use:   "show-snapshot [pair-id] [id]",
		Short: "shows a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			pairId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			if pairId == 0 {
				return types.ErrInvalidPairId
			}

			id, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			params := &types.QueryGetSnapshotRequest{
				PairId: pairId,
				Id:     id,
			}

			res, err := queryClient.Snapshot(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
