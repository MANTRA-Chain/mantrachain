package cli

import (
	"github.com/MANTRA-Chain/mantrachain/v8/x/sanction/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

// GetSanctioners returns the sanctioner list with roles.
func GetSanctioners() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sanctioners [flags]",
		Short: "Get the sanctioner list with roles",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			res, err := types.NewQueryClient(clientCtx).Sanctioners(cmd.Context(), &types.QuerySanctionersRequest{Pagination: pageReq})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "sanctioners")
	return cmd
}

// GetIsSanctioned returns whether one address is sanctioned.
func GetIsSanctioned() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "is-sanctioned [address] [flags]",
		Short: "Check whether an address is sanctioned",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			res, err := types.NewQueryClient(clientCtx).IsSanctioned(cmd.Context(), &types.QueryIsSanctionedRequest{Address: args[0]})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
