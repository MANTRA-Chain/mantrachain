package cli

import (
	"context"
	"strings"

	"github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/MANTRA-Finance/mantrachain/x/coinfactory/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

func CmdQueryBalance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balance [creator] [subdenom] [address] [flags]",
		Short: "get the address balance for a specific subdenom by creator",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			reqCreator := args[0]
			reqSubDenom := args[1]

			if strings.TrimSpace(reqCreator) == "" {
				return errors.Wrap(types.ErrInvalidCreator, "empty creator")
			}

			if strings.TrimSpace(reqSubDenom) == "" {
				return errors.Wrap(types.ErrInvalidDenom, "empty subdenom")
			}

			addr, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			params := &types.QueryBalanceRequest{
				Creator:  reqCreator,
				Subdenom: reqSubDenom,
				Address:  addr.String(),
			}

			res, err := queryClient.Balance(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryDenomMetadata() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "denom-metadata [creator] [subdenom] [flags]",
		Short: "get the authority metadata for a specific denom",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			reqCreator := args[0]
			reqSubDenom := args[1]

			if strings.TrimSpace(reqCreator) == "" {
				return errors.Wrap(types.ErrInvalidCreator, "empty creator")
			}

			if strings.TrimSpace(reqSubDenom) == "" {
				return errors.Wrap(types.ErrInvalidDenom, "empty subdenom")
			}

			params := &types.QueryDenomMetadataRequest{
				Creator:  reqCreator,
				Subdenom: reqSubDenom,
			}

			res, err := queryClient.DenomMetadata(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQuerySupplyOf() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "supply [creator] [subdenom] [flags]",
		Short: "get the authority metadata for a specific denom",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			reqCreator := args[0]
			reqSubDenom := args[1]

			if strings.TrimSpace(reqCreator) == "" {
				return errors.Wrap(types.ErrInvalidCreator, "empty creator")
			}

			if strings.TrimSpace(reqSubDenom) == "" {
				return errors.Wrap(types.ErrInvalidDenom, "empty subdenom")
			}

			params := &types.QuerySupplyOfRequest{
				Creator:  reqCreator,
				Subdenom: reqSubDenom,
			}

			res, err := queryClient.SupplyOf(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
