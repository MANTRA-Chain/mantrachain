package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"mantrachain/x/token/types"
)

var _ = strconv.Itoa(0)

func CmdUpdateGuardSoulBondNftImage() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-guard-soul-bond-nft-image [nft-id] [image-index] [payload-json]",
		Short: "Update guard soul-bond nft image by index",
		Long: "Update guard soul-bond NFT image. " +
			"[payload-json] is JSON encoded MsgNftImageMetadata.",
		Example: fmt.Sprintf(
			"$ %s tx token update-guard-soul-bond-nft-image <nft-id> <image-index> <payload-json> "+
				"--from=<from> "+
				"--owner=<owner> "+
				"--chain-id=<chain-id> ",
			version.AppName,
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argNftId := args[0]
			argIndex := args[1]
			argMetadata := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// verification
			signer := clientCtx.GetFromAddress()

			owner, err := cmd.Flags().GetString(FlagOwner)
			if err != nil {
				return err
			}

			ownerStr := strings.TrimSpace(owner)
			if len(ownerStr) == 0 {
				owner = signer.String()
			}

			// Unmarshal payload
			var image types.MsgNftImageMetadata
			err = clientCtx.Codec.UnmarshalJSON([]byte(argMetadata), &image)
			if err != nil {
				return err
			}

			index, err := strconv.ParseUint(argIndex, 10, 64)
			if err != nil {
				return fmt.Errorf("image index: %w", err)
			}

			msg := types.NewMsgUpdateGuardSoulBondNftImage(
				clientCtx.GetFromAddress().String(),
				owner,
				argNftId,
				index,
				&image,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsUpdateGuardSoulBondNFTImage)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdMintNfts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint-nfts [payload-json]",
		Short: "Broadcast message mint_nfts",
		Long: "Mints new NFTs. " +
			"[payload-json] is JSON encoded MsgNftsMetadata.",
		Example: fmt.Sprintf(
			"$ %s tx token mint-nfts <payload-json> "+
				"--from=<from> "+
				"--receiver=<receiver> "+
				"--collection-creator=<collection-creator> "+
				"--collection-id=<collection-id> "+
				"--did "+
				"--strict "+
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

			receiver, err := cmd.Flags().GetString(FlagReceiver)
			if err != nil {
				return err
			}

			receiverStr := strings.TrimSpace(receiver)
			if len(receiverStr) == 0 {
				receiver = signer.String()
			}

			strict, err := cmd.Flags().GetBool(FlagStrict)
			if err != nil {
				return err
			}

			did, err := cmd.Flags().GetBool(FlagDid)
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
				signer.String(),
				collectionCreator,
				collectionId,
				&nfts,
				receiver,
				strict,
				did,
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
			"$ %s tx token burn-nfts <payload-json> "+
				"--from=<from> "+
				"--collection-creator=<collection-creator> "+
				"--collection-id=<collection-id> "+
				"--strict "+
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

			strict, err := cmd.Flags().GetBool(FlagStrict)
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
				strict,
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
			"$ %s tx token transfer-nfts <owner> <receiver> <payload-json> "+
				"--from=<from> "+
				"--collection-creator=<collection-creator> "+
				"--collection-id=<collection-id> "+
				"--strict "+
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

			strict, err := cmd.Flags().GetBool(FlagStrict)
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
				strict,
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
			"$ %s tx token approve-nfts <operator> <approved> <payload-json> "+
				"--from=<from> "+
				"--collection-creator=<collection-creator> "+
				"--collection-id=<collection-id> "+
				"--strict "+
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

			strict, err := cmd.Flags().GetBool(FlagStrict)
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
				strict,
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
			"$ %s tx token approve-all-nfts <operator> <approved> "+
				"--from=<from> "+
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

			msg := types.NewMsgApproveAllNfts(
				clientCtx.GetFromAddress().String(),
				argReceiver,
				approved,
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
			"$ %s tx token mint-nft <payload-json> "+
				"--from=<from> "+
				"--receiver=<receiver> "+
				"--collection-creator=<collection-creator> "+
				"--collection-id=<collection-id> "+
				"--strict "+
				"--did "+
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

			receiver, err := cmd.Flags().GetString(FlagReceiver)
			if err != nil {
				return err
			}

			receiverStr := strings.TrimSpace(receiver)
			if len(receiverStr) == 0 {
				receiver = signer.String()
			}

			strict, err := cmd.Flags().GetBool(FlagStrict)
			if err != nil {
				return err
			}

			did, err := cmd.Flags().GetBool(FlagDid)
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
				signer.String(),
				collectionCreator,
				collectionId,
				&nft,
				receiver,
				strict,
				did,
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
			"$ %s tx token burn-nft <nft-id> "+
				"--from=<from> "+
				"--collection-creator=<collection-creator> "+
				"--collection-id=<collection-id> "+
				"--strict "+
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

			strict, err := cmd.Flags().GetBool(FlagStrict)
			if err != nil {
				return err
			}

			msg := types.NewMsgBurnNft(
				clientCtx.GetFromAddress().String(),
				collectionCreator,
				collectionId,
				argNftId,
				strict,
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
			"$ %s tx token transfer-nft <owner> <receiver> <nft-id> "+
				"--from=<from> "+
				"--collection-creator=<collection-creator> "+
				"--collection-id=<collection-id> "+
				"--strict "+
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

			strict, err := cmd.Flags().GetBool(FlagStrict)
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
				strict,
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
			"$ %s tx token approve-nft <operator> <approved> <nft-id> "+
				"--from=<from> "+
				"--collection-creator=<collection-creator> "+
				"--collection-id=<collection-id> "+
				"--strict "+
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

			strict, err := cmd.Flags().GetBool(FlagStrict)
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
				strict,
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
