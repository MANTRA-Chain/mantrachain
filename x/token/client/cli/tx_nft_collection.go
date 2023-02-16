package cli

import (
	"fmt"
	"strconv"

	"github.com/LimeChain/mantrachain/x/token/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdCreateNftCollection() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-nft-collection [payload-json]",
		Short: "Broadcast message create_nft_collection",
		Long: "Creates a new NFT collection. " +
			"[payload-json] is JSON encoded MsgCreateNftCollectionMetadata.",
		Example: fmt.Sprintf(
			"$ %s tx token create-nft-collection <payload-json> "+
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
			var metadata types.MsgCreateNftCollectionMetadata
			err = clientCtx.Codec.UnmarshalJSON([]byte(argMetadata), &metadata)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateNftCollection(
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
