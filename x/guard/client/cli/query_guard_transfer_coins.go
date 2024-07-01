package cli

import (
	"context"

	"github.com/MANTRA-Finance/mantrachain/x/guard/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdShowGuardTransferCoins() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-guard-transfer-coins",
		Short: "shows guard_transfer_coins",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetGuardTransferCoinsRequest{}

			res, err := queryClient.QueryGuardTransferCoins(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
