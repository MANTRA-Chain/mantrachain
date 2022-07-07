package cli

import (
	"strconv"
	"strings"

	"github.com/LimeChain/mantrachain/x/mdb/types"
	"github.com/LimeChain/mantrachain/x/mdb/utils"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdMintNfts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint-nfts [collection_creator] [collection_id] [payload-json] --from [creator] --receiver [receiver] --strict [true/false]",
		Short: "Broadcast message mint_nfts",
		Long: "Mints new NFTS. " +
			"[payload-json] is JSON encoded MsgMintNfts.",
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

			receiver, err := cmd.Flags().GetString(FlagReceiver)
			if err != nil {
				return err
			}

			receiverStr := strings.TrimSpace(receiver)
			if len(receiverStr) > 0 {
				if _, err = sdk.AccAddressFromBech32(receiverStr); err != nil {
					return err
				}
			} else {
				receiver = signer.String()
			}

			pubKey := info.GetPubKey()
			pubKeyHex := utils.GetPubKeyHex(pubKey)
			pubKeyType, err := utils.DerivePubKeyType(pubKey)
			if err != nil {
				return err
			}

			strictFlag, err := cmd.Flags().GetString(FlagStrict)
			if err != nil {
				return err
			}

			strict, err := strconv.ParseBool(strictFlag)
			if err != nil {
				return err
			}

			// Unmarshal payload
			var nfts types.MsgNftsMetadata
			err = clientCtx.Codec.UnmarshalJSON([]byte(argMetadata), &nfts)
			if err != nil {
				return err
			}

			msg := types.NewMsgMintNfts(
				clientCtx.GetFromAddress().String(),
				argCollectionCreator,
				argCollectionId,
				&nfts,
				receiver,
				pubKeyHex,
				pubKeyType,
				strict,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsNFT)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdBurnNfts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn-nfts [collection_creator] [collection_id] [payload-json] --from [creator]",
		Short: "Broadcast message burn_nfts",
		Long: "Burns NFTS. " +
			"[payload-json] is JSON encoded MsgBurnNfts.",
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

			strictFlag, err := cmd.Flags().GetString(FlagStrict)
			if err != nil {
				return err
			}

			strict, err := strconv.ParseBool(strictFlag)
			if err != nil {
				return err
			}

			// Unmarshal payload
			var nfts types.MsgNftsIds
			err = clientCtx.Codec.UnmarshalJSON([]byte(argMetadata), &nfts)
			if err != nil {
				return err
			}

			msg := types.NewMsgBurnNfts(
				clientCtx.GetFromAddress().String(),
				argCollectionCreator,
				argCollectionId,
				&nfts,
				pubKeyHex,
				pubKeyType,
				strict,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsNFT)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
