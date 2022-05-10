package cli

import (
	"context"

	"github.com/LimeChain/mantrachain/x/mns/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdGetDomainName() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "domain-name [domain] [domain-name]",
		Short: "Query a domain-name",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argDomain := args[0]
			argDomainName := args[1]

			params := &types.QueryGetDomainNameRequest{
				Domain:     argDomain,
				DomainName: argDomainName,
			}

			res, err := queryClient.DomainName(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
