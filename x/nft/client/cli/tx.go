package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/LimeChain/mantrachain/x/nft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	nftTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "nft transactions subcommands",
		Long:                       "Provides the most common nft logic for upper-level applications, compatible with Ethereum's erc721 contract",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	// NFT Send is disabled since we don't want to allow users to transfer soul bounded NFTs
	// nftTxCmd.AddCommand(
	// 	NewCmdSend(),
	// )

	return nftTxCmd
}

func NewCmdSend() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [class-id] [nft-id] [receiver] --from [sender]",
		Args:  cobra.ExactArgs(3),
		Short: "transfer ownership of nft",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s send <class-id> <nft-id> <receiver> --from <sender> --chain-id <chain-id>`, version.AppName, types.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgSend{
				ClassId:  args[0],
				Id:       args[1],
				Sender:   clientCtx.GetFromAddress().String(),
				Receiver: args[2],
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
