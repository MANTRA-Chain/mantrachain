package cli

import (
	"context"
	"strconv"
	"strings"

	"github.com/LimeChain/mantrachain/x/marketplace/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdGetMarketplaceCollection() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "marketplace-collection [marketplace_creator] [marketplace_id] [collection_creator] [collection_id]",
		Short: "Query a marketplace collection",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			reqMarketplaceCreator := args[0]
			reqMarketplaceId := args[1]
			reqCollectionCreator := args[2]
			reqCollectionId := args[3]

			if strings.TrimSpace(reqMarketplaceId) == "" {
				return sdkerrors.Wrap(types.ErrInvalidMarketplaceId, "empty marketplace id")
			}

			if strings.TrimSpace(reqCollectionId) == "" {
				return sdkerrors.Wrap(types.ErrInvalidCollectionId, "empty collection id")
			}

			marketplaceCreator, err := sdk.AccAddressFromBech32(reqMarketplaceCreator)
			if err != nil {
				return err
			}

			collectionCreator, err := sdk.AccAddressFromBech32(reqCollectionCreator)
			if err != nil {
				return err
			}

			params := &types.QueryGetMarketplaceCollectionRequest{
				MarketplaceCreator: marketplaceCreator.String(),
				MarketplaceId:      reqMarketplaceId,
				CollectionCreator:  collectionCreator.String(),
				CollectionId:       reqCollectionId,
			}

			res, err := queryClient.MarketplaceCollection(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdGetAllMarketplaceCollections() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-marketplace-collections [marketplace-creator] [marketplace-id]",
		Short: "Query all marketplace collections",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			reqMarketplaceCreator := args[0]
			reqMarketplaceId := args[1]

			if strings.TrimSpace(reqMarketplaceId) == "" {
				return sdkerrors.Wrap(types.ErrInvalidMarketplaceId, "empty nft marketplace id")
			}

			marketplaceCreator, err := sdk.AccAddressFromBech32(reqMarketplaceCreator)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryGetAllMarketplaceCollectionsRequest{
				MarketplaceCreator: marketplaceCreator.String(),
				MarketplaceId:      reqMarketplaceId,
				Pagination:         pageReq,
			}

			res, err := queryClient.AllMarketplaceCollections(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "all-marketplace-collections")

	return cmd
}
