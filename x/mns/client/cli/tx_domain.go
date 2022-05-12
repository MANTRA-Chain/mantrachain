package cli

import (
	"strconv"

	"github.com/LimeChain/mantrachain/x/mns/types"
	"github.com/LimeChain/mantrachain/x/mns/utils"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdCreateDomain() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-domain [domain] [type]",
		Short: "Broadcast message create-domain",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argDomain := args[0]
			argDomainType := args[1]

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

			msg := types.NewMsgCreateDomain(
				clientCtx.GetFromAddress().String(),
				argDomain,
				argDomainType,
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
