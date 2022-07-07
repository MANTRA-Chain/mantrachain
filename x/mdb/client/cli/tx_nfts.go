package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/LimeChain/mantrachain/x/mdb/types"
	"github.com/LimeChain/mantrachain/x/mdb/utils"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdMintNfts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint-nfts [payload-json]",
		Short: "Broadcast message mint_nfts",
		Long: "Mints new NFTs. " +
			"[payload-json] is JSON encoded MsgNftsMetadata.",
		Example: fmt.Sprintf(
			"$ %s tx mdb mint-nfts <payload-json> "+
				"--from=<key-name> "+
				"--receiver=<receiver> "+
				"--collection-creator=<collection-creator> "+
				"--collection-id=<collection-id> "+
				"--strict-collection=<true/false> "+
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

			collectionCreator, err := cmd.Flags().GetString(FlagCollectionCreator)
			if err != nil {
				return err
			}

			collectionId, err := cmd.Flags().GetString(FlagCollectionId)
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
			if len(receiverStr) == 0 {
				receiver = signer.String()
			}

			pubKey := info.GetPubKey()
			pubKeyHex := utils.GetPubKeyHex(pubKey)
			pubKeyType, err := utils.DerivePubKeyType(pubKey)
			if err != nil {
				return err
			}

			strictCollectionFlag, err := cmd.Flags().GetString(FlagStrictCollection)
			if err != nil {
				return err
			}

			strictCollection, err := strconv.ParseBool(strictCollectionFlag)
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
				collectionCreator,
				collectionId,
				&nfts,
				receiver,
				pubKeyHex,
				pubKeyType,
				strictCollection,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsMintNFT)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdBurnNfts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn-nfts [payload-json]",
		Short: "Broadcast message burn_nfts",
		Long: "Burns NFTs. " +
			"[payload-json] is JSON encoded MsgNftsIds.",
		Example: fmt.Sprintf(
			"$ %s tx mdb burn-nfts <payload-json> "+
				"--from=<key-name> "+
				"--collection-creator=<collection-creator> "+
				"--collection-id=<collection-id> "+
				"--strict-collection=<true/false> "+
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

			collectionCreator, err := cmd.Flags().GetString(FlagCollectionCreator)
			if err != nil {
				return err
			}

			collectionId, err := cmd.Flags().GetString(FlagCollectionId)
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

			strictCollectionFlag, err := cmd.Flags().GetString(FlagStrictCollection)
			if err != nil {
				return err
			}

			strictCollection, err := strconv.ParseBool(strictCollectionFlag)
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
				collectionCreator,
				collectionId,
				&nfts,
				pubKeyHex,
				pubKeyType,
				strictCollection,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsBurnNFT)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdTransferNfts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-nfts [from] [to] [payload-json]",
		Short: "Broadcast message transfer_nfts",
		Long: "Transfer NFTs to a recipient. " +
			"[from] is NFTs owner. " +
			"[to] is NFTs receiver. " +
			"[payload-json] is JSON encoded MsgNftsIds.",
		Example: fmt.Sprintf(
			"$ %s tx mdb transfer-nfts <owner> <receiver> <payload-json> "+
				"--from=<key-name> "+
				"--collection-creator=<collection-creator> "+
				"--collection-id=<collection-id> "+
				"--strict-collection=<true/false> "+
				"--chain-id=<chain-id> ",
			version.AppName,
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argOwner := args[0]
			argReceiver := args[1]
			argMetadata := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			collectionCreator, err := cmd.Flags().GetString(FlagCollectionCreator)
			if err != nil {
				return err
			}

			collectionId, err := cmd.Flags().GetString(FlagCollectionId)
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

			strictCollectionFlag, err := cmd.Flags().GetString(FlagStrictCollection)
			if err != nil {
				return err
			}

			strictCollection, err := strconv.ParseBool(strictCollectionFlag)
			if err != nil {
				return err
			}

			// Unmarshal payload
			var nftsIds types.MsgNftsIds
			err = clientCtx.Codec.UnmarshalJSON([]byte(argMetadata), &nftsIds)
			if err != nil {
				return err
			}

			msg := types.NewMsgTransferNfts(
				clientCtx.GetFromAddress().String(),
				collectionCreator,
				collectionId,
				&nftsIds,
				argOwner,
				argReceiver,
				pubKeyHex,
				pubKeyType,
				strictCollection,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsTransferNFT)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdApproveNfts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve-nfts [operator] [approved] [payload-json]",
		Short: "Broadcast message approve_nfts",
		Long: "Adds/Removes the address to the approved ntfs lists. " +
			"[operator] is the operator who can transfer the NTFs. " +
			"[approved] is boolean. " +
			"[payload-json] is JSON encoded MsgNftsIds.",
		Example: fmt.Sprintf(
			"$ %s tx mdb approve-nfts <operator> <approved> <payload-json> "+
				"--from=<key-name> "+
				"--collection-creator=<collection-creator> "+
				"--collection-id=<collection-id> "+
				"--strict-collection=<true/false> "+
				"--chain-id=<chain-id> ",
			version.AppName,
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argReceiver := args[0]
			argApproved := args[1]
			argMetadata := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			collectionCreator, err := cmd.Flags().GetString(FlagCollectionCreator)
			if err != nil {
				return err
			}

			collectionId, err := cmd.Flags().GetString(FlagCollectionId)
			if err != nil {
				return err
			}

			approved, err := strconv.ParseBool(argApproved)
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

			strictCollectionFlag, err := cmd.Flags().GetString(FlagStrictCollection)
			if err != nil {
				return err
			}

			strictCollection, err := strconv.ParseBool(strictCollectionFlag)
			if err != nil {
				return err
			}

			// Unmarshal payload
			var nftsIds types.MsgNftsIds
			err = clientCtx.Codec.UnmarshalJSON([]byte(argMetadata), &nftsIds)
			if err != nil {
				return err
			}

			msg := types.NewMsgApproveNfts(
				clientCtx.GetFromAddress().String(),
				argReceiver,
				collectionCreator,
				collectionId,
				&nftsIds,
				approved,
				pubKeyHex,
				pubKeyType,
				strictCollection,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsApproveNFT)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdApproveAllNfts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve-all-nfts [operator] [approved]",
		Short: "Broadcast message approve_all_nfts",
		Long: "Adds/Removes the address to the globally approved owner's list. " +
			"[operator] is the operator who can transfer the NTFs. " +
			"[approved] is boolean. " +
			"[payload-json] is JSON encoded MsgNftsIds.",
		Example: fmt.Sprintf(
			"$ %s tx mdb approve-all-nfts <operator> <approved> "+
				"--from=<key-name> "+
				"--chain-id=<chain-id> ",
			version.AppName,
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argReceiver := args[0]
			argApproved := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			approved, err := strconv.ParseBool(argApproved)
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

			msg := types.NewMsgApproveAllNfts(
				clientCtx.GetFromAddress().String(),
				argReceiver,
				approved,
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

func CmdMintNft() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint-nft [payload-json]",
		Short: "Broadcast message mint_nft",
		Long: "Mint a new NFT. " +
			"[payload-json] is JSON encoded MsgNftsMetadata.",
		Example: fmt.Sprintf(
			"$ %s tx mdb mint-nft <payload-json> "+
				"--from=<key-name> "+
				"--receiver=<receiver> "+
				"--collection-creator=<collection-creator> "+
				"--collection-id=<collection-id> "+
				"--strict-collection=<true/false> "+
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

			collectionCreator, err := cmd.Flags().GetString(FlagCollectionCreator)
			if err != nil {
				return err
			}

			collectionId, err := cmd.Flags().GetString(FlagCollectionId)
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
			if len(receiverStr) == 0 {
				receiver = signer.String()
			}

			pubKey := info.GetPubKey()
			pubKeyHex := utils.GetPubKeyHex(pubKey)
			pubKeyType, err := utils.DerivePubKeyType(pubKey)
			if err != nil {
				return err
			}

			strictCollectionFlag, err := cmd.Flags().GetString(FlagStrictCollection)
			if err != nil {
				return err
			}

			strictCollection, err := strconv.ParseBool(strictCollectionFlag)
			if err != nil {
				return err
			}

			// Unmarshal payload
			var nft types.MsgNftMetadata
			err = clientCtx.Codec.UnmarshalJSON([]byte(argMetadata), &nft)
			if err != nil {
				return err
			}

			msg := types.NewMsgMintNft(
				clientCtx.GetFromAddress().String(),
				collectionCreator,
				collectionId,
				&nft,
				receiver,
				pubKeyHex,
				pubKeyType,
				strictCollection,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsMintNFT)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdBurnNft() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn-nft [nft-id]",
		Short: "Broadcast message burn_nft",
		Long: "Burns NFT. " +
			"[nft-id] is the NFT id.",
		Example: fmt.Sprintf(
			"$ %s tx mdb burn-nft <nft-id> "+
				"--from=<key-name> "+
				"--collection-creator=<collection-creator> "+
				"--collection-id=<collection-id> "+
				"--strict-collection=<true/false> "+
				"--chain-id=<chain-id> ",
			version.AppName,
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argNftId := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			collectionCreator, err := cmd.Flags().GetString(FlagCollectionCreator)
			if err != nil {
				return err
			}

			collectionId, err := cmd.Flags().GetString(FlagCollectionId)
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

			strictCollectionFlag, err := cmd.Flags().GetString(FlagStrictCollection)
			if err != nil {
				return err
			}

			strictCollection, err := strconv.ParseBool(strictCollectionFlag)
			if err != nil {
				return err
			}

			msg := types.NewMsgBurnNft(
				clientCtx.GetFromAddress().String(),
				collectionCreator,
				collectionId,
				argNftId,
				pubKeyHex,
				pubKeyType,
				strictCollection,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsBurnNFT)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdTransferNft() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-nft [from] [to] [nft-id]",
		Short: "Broadcast message transfer_nft",
		Long: "Transfer NFT to a recipient. " +
			"[from] is NFT owner. " +
			"[to] is NFT receiver. " +
			"[nft-id] is the NFT id.",
		Example: fmt.Sprintf(
			"$ %s tx mdb transfer-nft <owner> <receiver> <nft-id> "+
				"--from=<key-name> "+
				"--collection-creator=<collection-creator> "+
				"--collection-id=<collection-id> "+
				"--strict-collection=<true/false> "+
				"--chain-id=<chain-id> ",
			version.AppName,
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argOwner := args[0]
			argReceiver := args[1]
			argNftId := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			collectionCreator, err := cmd.Flags().GetString(FlagCollectionCreator)
			if err != nil {
				return err
			}

			collectionId, err := cmd.Flags().GetString(FlagCollectionId)
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

			strictCollectionFlag, err := cmd.Flags().GetString(FlagStrictCollection)
			if err != nil {
				return err
			}

			strictCollection, err := strconv.ParseBool(strictCollectionFlag)
			if err != nil {
				return err
			}

			msg := types.NewMsgTransferNft(
				clientCtx.GetFromAddress().String(),
				collectionCreator,
				collectionId,
				argNftId,
				argOwner,
				argReceiver,
				pubKeyHex,
				pubKeyType,
				strictCollection,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsTransferNFT)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdApproveNft() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve-nft [operator] [approved] [nft-id]",
		Short: "Broadcast message approve_nft",
		Long: "Adds/Removes the address to the approved ntf list. " +
			"[operator] is the operator who can transfer the NTF. " +
			"[approved] is boolean. " +
			"[nft-id] is the NFT id.",
		Example: fmt.Sprintf(
			"$ %s tx mdb approve-nft <operator> <approved> <nft-id> "+
				"--from=<key-name> "+
				"--collection-creator=<collection-creator> "+
				"--collection-id=<collection-id> "+
				"--strict-collection=<true/false> "+
				"--chain-id=<chain-id> ",
			version.AppName,
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argReceiver := args[0]
			argApproved := args[1]
			argNftId := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			collectionCreator, err := cmd.Flags().GetString(FlagCollectionCreator)
			if err != nil {
				return err
			}

			collectionId, err := cmd.Flags().GetString(FlagCollectionId)
			if err != nil {
				return err
			}

			approved, err := strconv.ParseBool(argApproved)
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

			strictCollectionFlag, err := cmd.Flags().GetString(FlagStrictCollection)
			if err != nil {
				return err
			}

			strictCollection, err := strconv.ParseBool(strictCollectionFlag)
			if err != nil {
				return err
			}

			msg := types.NewMsgApproveNft(
				clientCtx.GetFromAddress().String(),
				argReceiver,
				collectionCreator,
				collectionId,
				argNftId,
				approved,
				pubKeyHex,
				pubKeyType,
				strictCollection,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsApproveNFT)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
