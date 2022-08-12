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

func CmdImportCollection() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import-collection [payload-json]",
		Short: "Broadcast message import-collection",
		Long: "Imports a NFT collection. " +
			"[payload-json] is JSON encoded MsgCollectionSettings.",
		Example: fmt.Sprintf(
			"$ %s tx marketplace import-collection <payload-json> "+
				"--from=<from> "+
				"--marketplace-creator=<marketplace-creator> "+
				"--marketplace-id=<marketplace-id> "+
				"--collection-creator=<collection-creator> "+
				"--collection-id=<collection-id> "+
				"--chain-id=<chain-id> ",
			version.AppName,
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argSettings := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			marketplaceCreator, err := cmd.Flags().GetString(FlagMarketplaceCreator)
			if err != nil {
				return err
			}

			marketplaceId, err := cmd.Flags().GetString(FlagMarketplaceId)
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

			// Unmarshal payload
			var settings types.MsgCollectionSettings
			err = clientCtx.Codec.UnmarshalJSON([]byte(argSettings), &settings)
			if err != nil {
				return err
			}

			msg := types.NewMsgImportCollection(
				clientCtx.GetFromAddress().String(),
				marketplaceCreator,
				marketplaceId,
				collectionCreator,
				collectionId,
				&settings,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsImportCollection)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
