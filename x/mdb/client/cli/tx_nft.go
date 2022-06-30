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

func CmdMintNfts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint-nfts [collection_creator] [collection_id] [payload-json]",
		Short: "Broadcast message mint_nft_collection",
		Long: "Mints a new NFT. " +
			"[payload-json] is JSON encoded MsgMintNft.",
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argCollectionCreator := args[0]
			argCollectionId := args[1]
			argMetadata := args[2]

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
			var nfts types.MsgMintNftsMetadata
			err = clientCtx.Codec.UnmarshalJSON([]byte(argMetadata), &nfts)
			if err != nil {
				return err
			}

			msg := types.NewMsgMintNfts(
				clientCtx.GetFromAddress().String(),
				argCollectionCreator,
				argCollectionId,
				&nfts,
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
