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
		Short: "Broadcast message withdraw-nft-rewards",
		Long: "Withdraw NFT rewards. " +
			"[nft-id] is the NFT id.",
		Example: fmt.Sprintf(
			"$ %s tx vault withdraw-nft-rewards <nft-id> "+
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

func CmdSetStaked() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-staked [nft-id] [block-height] [staked-index] [shares]",
		Short: "Broadcast message set-staked",
		Long: "Withdraw NFT rewards. " +
			"[nft-id] is the NFT id." +
			"[block-height] is the staking chain block height at the time of the delegate." +
			"[staked-index] is the staked element index in the array." +
			"[shares] is the shares amount staked on the validator.",
		Example: fmt.Sprintf(
			"$ %s tx vault set-staked <nft-id> <block-height> "+
				"--from=<from> "+
				"--marketplace-creator=<marketplace-creator> "+
				"--marketplace-id=<marketplace-id> "+
				"--collection-creator=<collection-creator> "+
				"--collection-id=<collection-id> "+
				"--staking-chain=<staking-chain> "+
				"--staking-validator=<staking-validator> "+
				"--chain-id=<chain-id> ",
			version.AppName,
		),
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argNftId := args[0]
			argBlockHeight := args[1]
			argStakedIndex := args[2]
			argShares := args[3]

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

			blockHeight, err := strconv.ParseInt(argBlockHeight, 10, 64)

			if err != nil {
				return err
			}

			stakedIndex, err := strconv.ParseInt(argStakedIndex, 10, 64)

			if err != nil {
				return err
			}

			msg := types.NewMsgSetStaked(
				clientCtx.GetFromAddress().String(),
				marketplaceCreator,
				marketplaceId,
				collectionCreator,
				collectionId,
				argNftId,
				stakingChain,
				stakingValidator,
				blockHeight,
				stakedIndex,
				argShares,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsSetStaked)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
