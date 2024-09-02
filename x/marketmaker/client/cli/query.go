package cli

// DONTCOVER
// client is excluded from test coverage in MVP version

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/MANTRA-Finance/mantrachain/x/marketmaker/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns a root CLI command handler for all x/marketmaker query commands.
func GetQueryCmd() *cobra.Command {
	mmQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the marketmaker module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	mmQueryCmd.AddCommand(
		GetQueryMarketMakersCmd(),
		GetCmdQueryIncentive(),
	)
	return mmQueryCmd
}

// GetQueryMarketMakersCmd implements the market maker query command.
func GetQueryMarketMakersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "marketmakers",
		Args:  cobra.MaximumNArgs(0),
		Short: "Query details of the market makers",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of the market makers

Example:
$ %s query %s marketmakers --pair-id=1
$ %s query %s marketmakers --address=...
$ %s query %s marketmakers --eligible=true...
`,
				version.AppName, types.ModuleName,
				version.AppName, types.ModuleName,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pairIdStr, _ := cmd.Flags().GetString(FlagPairId)
			mmAddr, _ := cmd.Flags().GetString(FlagAddress)
			eligibleStr, _ := cmd.Flags().GetString(FlagEligible)

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			req := &types.QueryMarketMakersRequest{
				Pagination: pageReq,
			}

			switch {
			case pairIdStr != "":
				pairId, err := strconv.ParseUint(pairIdStr, 10, 64)
				if err != nil {
					return fmt.Errorf("parse pair id: %w", err)
				}
				req.PairId = pairId
			case mmAddr != "":
				req.Address = mmAddr
			case eligibleStr != "":
				if _, err := strconv.ParseBool(eligibleStr); err != nil {
					return fmt.Errorf("parse eligible flag: %w", err)
				}
				req.Eligible = eligibleStr
			}

			res, err := queryClient.MarketMakers(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().AddFlagSet(flagSetMarketMakers())
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "marketmakers")

	return cmd
}

// GetCmdQueryIncentive implements the query market maker claimable incentive command.
func GetCmdQueryIncentive() *cobra.Command {
	bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()

	cmd := &cobra.Command{
		Use:   "incentive [mm-address]",
		Args:  cobra.ExactArgs(1),
		Short: "Query claimable incentive of a market maker",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query claimable incentive of a market maker.

Example:
$ %s query %s incentive %s1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
`,
				version.AppName, types.ModuleName, bech32PrefixAccAddr,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			mmAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			resp, err := queryClient.Incentive(context.Background(), &types.QueryIncentiveRequest{
				Address: mmAddr.String(),
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
