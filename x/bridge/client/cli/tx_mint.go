package cli

import (
	"fmt"
	"strconv"

	"github.com/LimeChain/mantrachain/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdMint() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint [payload-json]",
		Short: "Broadcast message mint",
		Long: "Mints CR20 Coins. " +
			"[payload-json] is JSON encoded MsgMintListMetadata.",
		Example: fmt.Sprintf(
			"$ %s tx bridge mint <payload-json> "+
				"--from=<from> "+
				"--bridge-creator=<bridge-creator> "+
				"--bridge-id=<bridge-id> "+
				"--chain-id=<chain-id> ",
			version.AppName,
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argMetadata := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			bridgeCreator, err := cmd.Flags().GetString(FlagBridgeCreator)
			if err != nil {
				return err
			}

			bridgeId, err := cmd.Flags().GetString(FlagBridgeId)
			if err != nil {
				return err
			}

			var mint types.MsgMintListMetadata
			err = clientCtx.Codec.UnmarshalJSON([]byte(argMetadata), &mint)
			if err != nil {
				return err
			}

			msg := types.NewMsgMint(
				clientCtx.GetFromAddress().String(),
				bridgeCreator,
				bridgeId,
				&mint,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsMint)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
