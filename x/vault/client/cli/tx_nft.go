package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/LimeChain/mantrachain/x/vault/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdWithdrawNftRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-nft-rewards [nft-id]",
		Short: "Broadcast message withdraw-nft-reward",
		Long: "Withdraw NFT rewards. " +
			"[nft-id] is the NFT id.",
		Example: fmt.Sprintf(
			"$ %s tx marketplace withdraw-nft-reward <nft-id> "+
				"--from=<from> "+
				"--receiver=<receiver> "+
				"--marketplace-creator=<marketplace-creator> "+
				"--marketplace-id=<marketplace-id> "+
				"--collection-creator=<collection-creator> "+
				"--collection-id=<collection-id> "+
				"--staking-chain=<staking-chain> "+
				"--staking-validator=<staking-validator> "+
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

			stakingChain, err := cmd.Flags().GetString(FlagStakingChain)
			if err != nil {
				return err
			}

			stakingValidator, err := cmd.Flags().GetString(FlagStakingValidator)
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

			msg := types.NewMsgWithdrawNftRewards(
				clientCtx.GetFromAddress().String(),
				marketplaceCreator,
				marketplaceId,
				collectionCreator,
				collectionId,
				argNftId,
				receiver,
				stakingChain,
				stakingValidator,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsWithdrawNftRewards)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
