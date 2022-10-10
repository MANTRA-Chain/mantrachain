package cli

import (
    "context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
    "github.com/LimeChain/mantrachain/x/bridge/types"
)

func CmdShowCw20Contract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-cw-20-contract",
		Short: "shows Cw20Contract",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
            clientCtx := client.GetClientContextFromCmd(cmd)

            queryClient := types.NewQueryClient(clientCtx)

            params := &types.QueryGetCw20ContractRequest{}

            res, err := queryClient.Cw20Contract(context.Background(), params)
            if err != nil {
                return err
            }

            return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

    return cmd
}
