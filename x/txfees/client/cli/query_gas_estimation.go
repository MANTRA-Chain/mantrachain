package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/MANTRA-Finance/mantrachain/x/txfees/types"
)

func CmdGasEstimation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gas-estimation [amount] [denom]",
		Short: "shows a gas estimation",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			coins, err := sdk.ParseCoinsNormalized(args[0])
			if err != nil {
				return err
			}

			if len(coins) != 1 {
				return types.ErrInvalidAmount
			}

			if args[1] == "" {
				return types.ErrInvalidFeeDenom
			}

			params := &types.QueryGetGasEstimationRequest{
				Amount: coins[0],
				Denom:  args[1],
			}

			res, err := queryClient.GasEstimation(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
