package cli

import (
	"context"

	"github.com/LimeChain/mantrachain/x/guard/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdShowGuardTransfer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-guard-transfer",
		Short: "shows guard_transfer",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetGuardTransferRequest{}

			res, err := queryClient.GuardTransfer(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
