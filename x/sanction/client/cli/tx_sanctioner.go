package cli

import (
	"strings"

	"github.com/MANTRA-Chain/mantrachain/v8/x/sanction/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

const flagReason = "reason"

// NewSanctionAddressesCmd broadcasts MsgSanctionAddresses.
func NewSanctionAddressesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sanction-addresses [addresses-separated-by-comma] --reason [reason] [flags]",
		Short: "Add addresses to the sanction list as a sanctioner",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			reason, err := cmd.Flags().GetString(flagReason)
			if err != nil {
				return err
			}
			return broadcastSanctionerMsg(cmd, func(signer string) sdk.Msg {
				return types.NewMsgSanctionAddresses(signer, strings.Split(args[0], ","), reason)
			})
		},
	}
	addReasonFlag(cmd)
	return cmd
}

// NewUnsanctionAddressesCmd broadcasts MsgUnsanctionAddresses.
func NewUnsanctionAddressesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unsanction-addresses [addresses-separated-by-comma] --reason [reason] [flags]",
		Short: "Remove addresses from the sanction list as a full sanctioner",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			reason, err := cmd.Flags().GetString(flagReason)
			if err != nil {
				return err
			}
			return broadcastSanctionerMsg(cmd, func(signer string) sdk.Msg {
				return types.NewMsgUnsanctionAddresses(signer, strings.Split(args[0], ","), reason)
			})
		},
	}
	addReasonFlag(cmd)
	return cmd
}

// NewRevokeSanctionerCmd broadcasts MsgRevokeSanctioner.
func NewRevokeSanctionerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke-sanctioner [address] --reason [reason] [flags]",
		Short: "Remove an add-only sanctioner from the sanctioner list as a full sanctioner",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			reason, err := cmd.Flags().GetString(flagReason)
			if err != nil {
				return err
			}
			return broadcastSanctionerMsg(cmd, func(signer string) sdk.Msg {
				return types.NewMsgRevokeSanctioner(signer, args[0], reason)
			})
		},
	}
	addReasonFlag(cmd)
	return cmd
}

func addReasonFlag(cmd *cobra.Command) {
	cmd.Flags().String(flagReason, "", "Why this change is being made; emitted in the event and stored with the entry")
	_ = cmd.MarkFlagRequired(flagReason)
	flags.AddTxFlagsToCmd(cmd)
}

func broadcastSanctionerMsg(cmd *cobra.Command, build func(signer string) sdk.Msg) error {
	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return err
	}
	txf, err := tx.NewFactoryCLI(clientCtx, cmd.Flags())
	if err != nil {
		return err
	}
	msg := build(clientCtx.GetFromAddress().String())
	return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf.WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever), msg)
}
