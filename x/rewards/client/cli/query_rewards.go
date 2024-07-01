package cli

import (
	"context"
	"strconv"

	"github.com/MANTRA-Finance/mantrachain/x/rewards/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

func CmdRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rewards [provider] [pair-id]",
		Short: "Get rewards by provider and pair-id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			provider, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			pairId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			if pairId == 0 {
				return types.ErrInvalidPairId
			}

			params := &types.QueryGetRewardsRequest{
				Provider: provider.String(),
				PairId:   pairId,
			}

			res, err := queryClient.QueryRewards(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
