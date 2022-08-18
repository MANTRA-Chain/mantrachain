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

func CmdGetMarketplace() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "marketplace [creator] [id]",
		Short: "Query a marketplace",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			reqCreator := args[0]
			reqId := args[1]

			if strings.TrimSpace(reqId) == "" {
				return sdkerrors.Wrap(types.ErrInvalidMarketplaceId, "empty marketplace id")
			}

			creator, err := sdk.AccAddressFromBech32(reqCreator)
			if err != nil {
				return err
			}

			params := &types.QueryGetMarketplaceRequest{
				Creator: creator.String(),
				Id:      reqId,
			}

			res, err := queryClient.Marketplace(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdGetMarketplacesByCreator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "marketplaces [creator]",
		Short: "Query marketplaces",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			reqCreator := args[0]

			creator, err := sdk.AccAddressFromBech32(reqCreator)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryGetMarketplacesByCreatorRequest{
				Creator:    creator.String(),
				Pagination: pageReq,
			}

			res, err := queryClient.MarketplacesByCreator(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "marketplaces")

	return cmd
}

func CmdGetAllMarketplaces() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-marketplaces",
		Short: "Query all marketplaces",
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

			params := &types.QueryGetAllMarketplacesRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.AllMarketplaces(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "all-marketplaces")

	return cmd
}

func CmdGetAllMarketplaceColelctions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-marketplace-collections [creator] [id]",
		Short: "Query all marketplace collections",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			reqCreator := args[0]
			reqId := args[1]

			if strings.TrimSpace(reqId) == "" {
				return sdkerrors.Wrap(types.ErrInvalidMarketplaceId, "empty marketplace id")
			}

			creator, err := sdk.AccAddressFromBech32(reqCreator)
			if err != nil {
				return err
			}

			params := &types.QueryGetAllMarketplaceCollectionsRequest{
				MarketplaceCreator: creator.String(),
				MarketplaceId:      reqId,
			}

			res, err := queryClient.AllMarketplaceCollections(context.Background(), params)
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
