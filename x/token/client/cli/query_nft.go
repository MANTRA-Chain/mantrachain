package cli

import (
	"context"
	"strconv"
	"strings"

	"github.com/LimeChain/mantrachain/x/token/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

			reqCollectionCreator := args[0]
			reqCollectionId := args[1]
			reqId := args[2]

			if strings.TrimSpace(reqCollectionId) == "" {
				return sdkerrors.Wrap(types.ErrInvalidNftCollectionId, "empty nft collection id")
			}

			if strings.TrimSpace(reqId) == "" {
				return sdkerrors.Wrap(types.ErrInvalidNftCollectionId, "empty nft id")
			}

			collectionCreator, err := sdk.AccAddressFromBech32(reqCollectionCreator)
			if err != nil {
				return err
			}

			params := &types.QueryGetNftRequest{
				CollectionCreator: collectionCreator.String(),
				CollectionId:      reqCollectionId,
				Id:                reqId,
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

func CmdGetAllCollectionNfts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-collection-nfts [collection-creator] [collection-id]",
		Short: "Query all collection nfts",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			reqCollectionCreator := args[0]
			reqCollectionId := args[1]

			if strings.TrimSpace(reqCollectionId) == "" {
				return sdkerrors.Wrap(types.ErrInvalidNftCollectionId, "empty nft collection id")
			}

			collectionCreator, err := sdk.AccAddressFromBech32(reqCollectionCreator)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryGetAllCollectionNftsRequest{
				CollectionCreator: collectionCreator.String(),
				CollectionId:      reqCollectionId,
				Pagination:        pageReq,
			}

			res, err := queryClient.AllCollectionNfts(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "all-collection-nfts")

	return cmd
}

func CmdGetNftOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nft-owner [collection_creator] [collection_id] [id]",
		Short: "Query a nft owner",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			reqCollectionCreator := args[0]
			reqCollectionId := args[1]
			reqId := args[2]

			if strings.TrimSpace(reqCollectionId) == "" {
				return sdkerrors.Wrap(types.ErrInvalidNftCollectionId, "empty nft collection id")
			}

			if strings.TrimSpace(reqId) == "" {
				return sdkerrors.Wrap(types.ErrInvalidNftCollectionId, "empty nft id")
			}

			collectionCreator, err := sdk.AccAddressFromBech32(reqCollectionCreator)
			if err != nil {
				return err
			}

			params := &types.QueryGetNftOwnerRequest{
				CollectionCreator: collectionCreator.String(),
				CollectionId:      reqCollectionId,
				Id:                reqId,
			}

			res, err := queryClient.NftOwner(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdGetIsApprovedForAllNfts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "is-approved-for-all-nfts [owner] [operator]",
		Short: "Query a operator is approved for all nfts",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			reqOwner := args[0]
			reqOperator := args[1]

			owner, err := sdk.AccAddressFromBech32(reqOwner)
			if err != nil {
				return err
			}

			operator, err := sdk.AccAddressFromBech32(reqOperator)
			if err != nil {
				return err
			}

			params := &types.QueryGetIsApprovedForAllNftsRequest{
				Owner:    owner.String(),
				Operator: operator.String(),
			}

			res, err := queryClient.IsApprovedForAllNfts(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdGetNftApproved() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nft-approvals [collection_creator] [collection_id] [id]",
		Short: "Query a nft approvals",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			reqCollectionCreator := args[0]
			reqCollectionId := args[1]
			reqId := args[2]

			if strings.TrimSpace(reqCollectionId) == "" {
				return sdkerrors.Wrap(types.ErrInvalidNftCollectionId, "empty nft collection id")
			}

			if strings.TrimSpace(reqId) == "" {
				return sdkerrors.Wrap(types.ErrInvalidNftCollectionId, "empty nft id")
			}

			collectionCreator, err := sdk.AccAddressFromBech32(reqCollectionCreator)
			if err != nil {
				return err
			}

			params := &types.QueryGetNftApprovedRequest{
				CollectionCreator: collectionCreator.String(),
				CollectionId:      reqCollectionId,
				Id:                reqId,
			}

			res, err := queryClient.NftApproved(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdGetNftBalance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nft-balance [collection_creator] [collection_id] [owner]",
		Short: "Query a nft collection balance by owner",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			reqCollectionCreator := args[0]
			reqCollectionId := args[1]
			reqOwner := args[2]

			if strings.TrimSpace(reqCollectionId) == "" {
				return sdkerrors.Wrap(types.ErrInvalidNftCollectionId, "empty nft collection id")
			}

			collectionCreator, err := sdk.AccAddressFromBech32(reqCollectionCreator)
			if err != nil {
				return err
			}

			owner, err := sdk.AccAddressFromBech32(reqOwner)
			if err != nil {
				return err
			}

			params := &types.QueryGetNftBalanceRequest{
				Owner:             owner.String(),
				CollectionCreator: collectionCreator.String(),
				CollectionId:      reqCollectionId,
			}

			res, err := queryClient.NftBalance(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdGetCollectionNftsByOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection-nfts [collection_creator] [collection_id] [owner]",
		Short: "Query collection nfts by owner",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			reqCollectionCreator := args[0]
			reqCollectionId := args[1]
			reqOwner := args[2]

			if strings.TrimSpace(reqCollectionId) == "" {
				return sdkerrors.Wrap(types.ErrInvalidNftCollectionId, "empty nft collection id")
			}

			collectionCreator, err := sdk.AccAddressFromBech32(reqCollectionCreator)
			if err != nil {
				return err
			}

			owner, err := sdk.AccAddressFromBech32(reqOwner)
			if err != nil {
				return err
			}

			params := &types.QueryGetCollectionNftsByOwnerRequest{
				Owner:             owner.String(),
				CollectionCreator: collectionCreator.String(),
				CollectionId:      reqCollectionId,
			}

			res, err := queryClient.CollectionNftsByOwner(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "collection-nfts-by-owner")

	return cmd
}