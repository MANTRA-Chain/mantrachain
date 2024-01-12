package cli

import (
	"context"
	"strconv"
	"strings"

	"github.com/AumegaChain/aumega/x/token/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
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
				return errors.Wrap(types.ErrInvalidNftCollectionId, "empty nft collection id")
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

func CmdShowNftCollectionOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-nft-collection-owner [index]",
		Short: "shows a nft_collection_owner",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argIndex := []byte(args[0])

			params := &types.QueryGetNftCollectionOwnerRequest{
				Index: argIndex,
			}

			res, err := queryClient.NftCollectionOwner(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowOpenedNftsCollection() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-opened-nfts-collection [index]",
		Short: "shows a opened_nfts_collection",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argIndex := []byte(args[0])

			params := &types.QueryGetOpenedNftsCollectionRequest{
				Index: argIndex,
			}

			res, err := queryClient.OpenedNftsCollection(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowRestrictedNftsCollection() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-restricted-nfts-collection [index]",
		Short: "shows a restricted_nfts_collection",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argIndex := []byte(args[0])

			params := &types.QueryGetRestrictedNftsCollectionRequest{
				Index: argIndex,
			}

			res, err := queryClient.RestrictedNftsCollection(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

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

			argIndex := []byte(args[0])

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
