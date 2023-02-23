package cli

import (
	"context"
	"strings"

	"cosmossdk.io/errors"

	"github.com/LimeChain/mantrachain/x/coinfactory/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdQueryDenomAuthorityMetadata() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "denom-authority-metadata [denom] [flags]",
		Short: "get the authority metadata for a specific denom",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			reqDenom := args[0]

			if strings.TrimSpace(reqDenom) == "" {
				return errors.Wrap(types.ErrInvalidDenom, "empty denom")
			}

			params := &types.QueryDenomAuthorityMetadataRequest{
				Denom: reqDenom,
			}

			res, err := queryClient.DenomAuthorityMetadata(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryDenomsFromCreator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "denoms-from-creator [creator] [flags]",
		Short: "returns a list of all tokens created by a specific creator address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			reqCreator := args[0]

			if strings.TrimSpace(reqCreator) == "" {
				return errors.Wrap(types.ErrInvalidCreator, "empty creator")
			}

			params := &types.QueryDenomsFromCreatorRequest{
				Creator: reqCreator,
			}

			res, err := queryClient.DenomsFromCreator(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
