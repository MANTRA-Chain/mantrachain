package cli

import (
	"context"
	"strconv"

	"github.com/LimeChain/mantrachain/x/mdb/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdGetNftCollection() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nft-collection [creator] [id]",
		Short: "Query a nftCollection",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			creator := args[0]
			id := args[1]

			params := &types.QueryGetNftCollectionRequest{
				Creator: creator,
				Id:      id,
			}

			res, err := queryClient.NftCollection(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdGetNftCollections() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nft-collections [creator]",
		Short: "Query nft-collections",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqCreator := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryGetNftCollectionsRequest{
				Creator:    reqCreator,
				Pagination: pageReq,
			}

			res, err := queryClient.NftCollections(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "nft-collections")

	return cmd
}

func CmdGetAllNftCollections() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-nft-collections",
		Short: "Query all-nft-collections",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryGetAllNftCollectionsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.AllNftCollections(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "all-nft-collections")

	return cmd
}
