package cli

import (
	"strconv"

	"github.com/LimeChain/mantrachain/x/mdb/types"
	"github.com/LimeChain/mantrachain/x/mdb/utils"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdCreateNftCollection() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-nft-collection [payload-json]",
		Short: "Broadcast message create_nft_collection",
		Long: "Creates a new NFT collection. " +
			"[payload-json] is JSON encoded MsgCreateNftCollectionMetadata.",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argMetadata := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// verification
			signer := clientCtx.GetFromAddress()
			// pubkey
			info, err := clientCtx.Keyring.KeyByAddress(signer)
			if err != nil {
				return err
			}

			pubKey := info.GetPubKey()
			pubKeyHex := utils.GetPubKeyHex(pubKey)
			pubKeyType, err := utils.DerivePubKeyType(pubKey)
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
				pubKeyHex,
				pubKeyType,
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
