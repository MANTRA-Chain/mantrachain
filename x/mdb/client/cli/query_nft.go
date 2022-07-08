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

func CmdGetNft() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nft [collection_creator] [collection_id] [id]",
		Short: "Query a nft",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			collection_creator := args[0]
			collection_id := args[1]
			id := args[2]

			params := &types.QueryGetNftRequest{
				CollectionCreator: collection_creator,
				CollectionId:      collection_id,
				Id:                id,
			}

			res, err := queryClient.Nft(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdGetCollectionNfts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection-nfts [collection-creator] [collection-id]",
		Short: "Query collection-nfts",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqCollectionCreator := args[0]
			reqCollectionId := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryGetCollectionNftsRequest{
				CollectionCreator: reqCollectionCreator,
				CollectionId:      reqCollectionId,
				Pagination:        pageReq,
			}

			res, err := queryClient.CollectionNfts(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "collection-nfts")

	return cmd
}
