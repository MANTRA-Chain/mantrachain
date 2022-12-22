package cli

import (
	"context"
	"strconv"
	"strings"

	"github.com/LimeChain/mantrachain/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdGetBridge() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bridge [bridge_creator] [bridge_id]",
		Short: "Query a bridge",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			reqBridgeCreator := args[0]
			reqBridgeId := args[1]

			if strings.TrimSpace(reqBridgeId) == "" {
				return sdkerrors.Wrap(types.ErrInvalidBridgeId, "empty bridge id")
			}

			bridgeCreator, err := sdk.AccAddressFromBech32(reqBridgeCreator)
			if err != nil {
				return err
			}

			params := &types.QueryGetBridgeRequest{
				BridgeCreator: bridgeCreator.String(),
				BridgeId:      reqBridgeId,
			}

			res, err := queryClient.Bridge(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
