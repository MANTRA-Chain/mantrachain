package cli

import (
	"context"
	"strconv"

	"github.com/LimeChain/mantrachain/x/vault/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdGetLastEpochs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "last-epochs [staking_chain] [staking_validator]",
		Short: "Query the current and previous epochs",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argStakingChain := args[0]
			argStakingValidator := args[1]

			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetLastEpochsRequest{
				StakingChain:     argStakingChain,
				StakingValidator: argStakingValidator,
			}

			res, err := queryClient.LastEpochs(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdGetLastEpochBlock() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "last-epoch-block [stakingchain] [stakingvalidator]",
		Short: "Query the last epoch block",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argStakingChain := args[0]
			argStakingValidator := args[1]

			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetLastEpochBlockRequest{
				StakingChain:     argStakingChain,
				StakingValidator: argStakingValidator,
			}

			res, err := queryClient.LastEpochBlock(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
