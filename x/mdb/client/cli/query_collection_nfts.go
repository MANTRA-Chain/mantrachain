package cli

import (
	"strconv"

	"github.com/LimeChain/mantrachain/x/mdb/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdCollectionNfts() *cobra.Command {
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

			params := &types.QueryGetCollectionNftsRequest{
				CollectionCreator: reqCollectionCreator,
				CollectionId:      reqCollectionId,
			}

			res, err := queryClient.CollectionNfts(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
