package cli

import (
	"strconv"

	"github.com/LimeChain/mantrachain/x/mdb/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdNftCollections() *cobra.Command {
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

			params := &types.QueryGetNftCollectionsRequest{
				Creator: reqCreator,
			}

			res, err := queryClient.NftCollections(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
