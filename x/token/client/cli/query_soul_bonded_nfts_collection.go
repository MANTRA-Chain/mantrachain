package cli

import (
    "context"
	
    "github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
    "github.com/LimeChain/mantrachain/x/token/types"
)

func CmdListSoulBondedNftsCollection() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-soul-bonded-nfts-collection",
		Short: "list all soul_bonded_nfts_collection",
		RunE: func(cmd *cobra.Command, args []string) error {
            clientCtx := client.GetClientContextFromCmd(cmd)

            pageReq, err := client.ReadPageRequest(cmd.Flags())
            if err != nil {
                return err
            }

            queryClient := types.NewQueryClient(clientCtx)

            params := &types.QueryAllSoulBondedNftsCollectionRequest{
                Pagination: pageReq,
            }

            res, err := queryClient.SoulBondedNftsCollectionAll(context.Background(), params)
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

func CmdShowSoulBondedNftsCollection() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-soul-bonded-nfts-collection [index]",
		Short: "shows a soul_bonded_nfts_collection",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
            clientCtx := client.GetClientContextFromCmd(cmd)

            queryClient := types.NewQueryClient(clientCtx)

             argIndex := args[0]
            
            params := &types.QueryGetSoulBondedNftsCollectionRequest{
                Index: argIndex,
                
            }

            res, err := queryClient.SoulBondedNftsCollection(context.Background(), params)
            if err != nil {
                return err
            }

            return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

    return cmd
}
