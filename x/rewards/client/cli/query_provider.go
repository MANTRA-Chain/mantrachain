package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/AumegaChain/aumega/x/rewards/types"
)

func CmdListProvider() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-provider",
		Short: "list all provider",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllProviderRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.ProviderAll(cmd.Context(), params)
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

func CmdProviderPairs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "provider-pairs [provider]",
		Short: "Shows a provider pairs ids",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			provider, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			params := &types.QueryGetProviderPairsRequest{
				Provider: provider.String(),
			}

			res, err := queryClient.ProviderPairs(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdProvider() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "provider [provider]",
		Short: "Shows a provider",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			provider, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			params := &types.QueryGetProviderRequest{
				Provider: provider.String(),
			}

			res, err := queryClient.Provider(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
