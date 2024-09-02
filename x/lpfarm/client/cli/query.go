package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/MANTRA-Finance/mantrachain/x/lpfarm/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewQueryPlansCmd(),
		NewQueryPlanCmd(),
		NewQueryFarmCmd(),
		NewQueryPositionsCmd(),
		NewQueryPositionCmd(),
		NewQueryHistoricalRewardsCmd(),
		NewQueryTotalRewardsCmd(),
		NewQueryRewardsCmd(),
	)

	return cmd
}

// NewQueryPlansCmd implements the plans query cmd.
func NewQueryPlansCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plans",
		Args:  cobra.NoArgs,
		Short: "Query all plans",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all plans.

Example:
$ %s query %s plans
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			res, err := queryClient.Plans(context.Background(), &types.QueryPlansRequest{
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "plans")
	return cmd
}

// NewQueryPlanCmd implements the plan query cmd.
func NewQueryPlanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plan [plan-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query a specific plan",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query a specific plan.

Example:
$ %s query %s plan 1
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			planId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid plan id: %w", err)
			}
			res, err := queryClient.Plan(context.Background(), &types.QueryPlanRequest{
				PlanId: planId,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// NewQueryFarmCmd implements the farm query cmd.
func NewQueryFarmCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "farm [denom]",
		Args:  cobra.ExactArgs(1),
		Short: "Query a specific farm for the denom",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query a specific farm for the denom.

Example:
$ %s query %s farm pool1
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Farm(context.Background(), &types.QueryFarmRequest{
				Denom: args[0],
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// NewQueryPositionsCmd implements the positions query cmd.
func NewQueryPositionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "positions [farmer]",
		Args:  cobra.ExactArgs(1),
		Short: "Query all the positions managed by the farmer",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all the positions managed by the farmer.

Example:
$ %s query %s positions cosmos1...
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			res, err := queryClient.Positions(context.Background(), &types.QueryPositionsRequest{
				Farmer:     args[0],
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "positions")
	return cmd
}

// NewQueryPositionCmd implements the position query cmd.
func NewQueryPositionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "position [farmer] [denom]",
		Args:  cobra.ExactArgs(2),
		Short: "Query a specific position managed by the farmer",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query a specific position managed by the farmer.

Example:
$ %s query %s position cosmos1... pool1
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Position(context.Background(), &types.QueryPositionRequest{
				Farmer: args[0],
				Denom:  args[1],
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// NewQueryHistoricalRewardsCmd implements the historical rewards query cmd.
func NewQueryHistoricalRewardsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "historical-rewards [denom]",
		Args:  cobra.ExactArgs(1),
		Short: "Query all historical rewards for the denom",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all historical rewards for the denom.

Example:
$ %s query %s historical-rewards pool1
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			res, err := queryClient.HistoricalRewards(context.Background(), &types.QueryHistoricalRewardsRequest{
				Denom:      args[0],
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "historical-rewards")
	return cmd
}

// NewQueryTotalRewardsCmd implements the total rewards query cmd.
func NewQueryTotalRewardsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total-rewards [farmer]",
		Args:  cobra.ExactArgs(1),
		Short: "Query total rewards accumulated in all farming assets of the farmer",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query total rewards accumulated in all farming assets of the farmer.

Example:
$ %s query %s all-rewards cosmos1...
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.TotalRewards(context.Background(), &types.QueryTotalRewardsRequest{
				Farmer: args[0],
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// NewQueryRewardsCmd implements the rewards query cmd.
func NewQueryRewardsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rewards [farmer] [denom]",
		Args:  cobra.ExactArgs(2),
		Short: "Query rewards accumulated in a farming asset of the farmer",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query rewards accumulated in a farming asset of the farmer.

Example:
$ %s query %s rewards cosmos1... pool1
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Rewards(context.Background(), &types.QueryRewardsRequest{
				Farmer: args[0],
				Denom:  args[1],
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
