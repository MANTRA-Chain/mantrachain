package cli

import (
	"fmt"
	"strconv"

	"github.com/LimeChain/mantrachain/x/marketplace/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdRegisterMarketplace() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-marketplace [marketplace]",
		Short: "Broadcast message register-marketplace",
		Long: "Registers a new NFT marketplace. " +
			"[payload-json] is JSON encoded MsgMarketplaceMetadata.",
		Example: fmt.Sprintf(
			"$ %s tx marketplace register-marketplace <payload-json> "+
				"--from=<from> "+
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

			// Unmarshal payload
			var metadata types.MsgMarketplaceMetadata
			err = clientCtx.Codec.UnmarshalJSON([]byte(argMetadata), &metadata)
			if err != nil {
				return err
			}

			msg := types.NewMsgRegisterMarketplace(
				clientCtx.GetFromAddress().String(),
				&metadata,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
