package cli

import (
	"context"
	"strconv"
	"strings"

	"github.com/LimeChain/mantrachain/x/mdb/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdGetNftCollection() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nft-collection [creator] [id]",
		Short: "Query a nft collection",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			reqCreator := args[0]
			reqId := args[1]

			if strings.TrimSpace(reqId) == "" {
				return sdkerrors.Wrap(types.ErrInvalidNftCollectionId, "empty nft collection id")
			}

			creator, err := sdk.AccAddressFromBech32(reqCreator)
			if err != nil {
				return err
			}

			params := &types.QueryGetNftCollectionRequest{
				Creator: creator.String(),
				Id:      reqId,
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

func CmdGetNftCollectionSupply() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nft-collection-supply [creator] [id]",
		Short: "Query a nft collection supply",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			reqCreator := args[0]
			reqId := args[1]

			if strings.TrimSpace(reqId) == "" {
				return sdkerrors.Wrap(types.ErrInvalidNftCollectionId, "empty nft collection id")
			}

			creator, err := sdk.AccAddressFromBech32(reqCreator)
			if err != nil {
				return err
			}

			params := &types.QueryGetNftCollectionSupplyRequest{
				Creator: creator.String(),
				Id:      reqId,
			}

			res, err := queryClient.NftCollectionSupply(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdGetNftCollectionsByCreator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nft-collections [creator]",
		Short: "Query nft collections",
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

			params := &types.QueryGetNftCollectionsByCreatorRequest{
				Creator:    creator.String(),
				Pagination: pageReq,
			}

			res, err := queryClient.NftCollectionsByCreator(cmd.Context(), params)
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
		Short: "Query all nft collections",
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
