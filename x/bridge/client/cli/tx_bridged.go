package cli

import (
	"fmt"
	"strings"

	"github.com/MANTRA-Finance/mantrachain/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

func CmdCreateMultiBridged() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-multi-bridged [from_key_or_address] [amount] [to_address_1,to_address_2,...] [amount_1,amount_2,...] [eth-tx-hash_1,eth-tx-hash_2,...]",
		Short: "Create multi bridged",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cmd.Flags().Set(flags.FlagFrom, args[0])
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			toAddressesStr := strings.TrimSpace(args[2])
			if len(toAddressesStr) == 0 {
				return fmt.Errorf("toAddresses cannot be empty")
			}

			toAddressStrs := strings.Split(toAddressesStr, ",")

			coinStrs := strings.Split(args[3], ",")
			coins := make([]sdk.Coin, 0, len(coinStrs))
			for _, coinStr := range coinStrs {
				newCoin, err := sdk.ParseCoinsNormalized(coinStr)
				if err != nil {
					return err
				}
				coins = append(coins, newCoin[0])
			}

			if len(coins) != len(toAddressStrs) {
				return fmt.Errorf("number of coins and toAddresses should be equal")
			}

			ethTxHashesStr := strings.TrimSpace(args[4])
			if len(ethTxHashesStr) == 0 {
				return fmt.Errorf("ethTxHashes cannot be empty")
			}

			ethTxHashes := strings.Split(ethTxHashesStr, ",")
			if len(ethTxHashes) != len(toAddressStrs) {
				return fmt.Errorf("number of ethTxHashes and toAddresses should be equal")
			}

			outputs := make([]types.Output, len(toAddressStrs))
			for i, toAddress := range toAddressStrs {
				outputs[i] = types.Output{
					Address: toAddress,
					Coins:   sdk.NewCoins(coins[i]),
				}
			}

			msg := types.NewMsgCreateMultiBridged(
				types.Input{
					Address: clientCtx.FromAddress.String(),
					Coins:   sdk.NewCoins(amount),
				},
				outputs,
				ethTxHashes,
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
