package cli

import (
	"github.com/LimeChain/mantrachain/x/mns/types"
	"github.com/LimeChain/mantrachain/x/mns/utils"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdCreateDomainName() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-domain-name [domain] [domain-name]",
		Short: "Create a new domainName",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get value arguments
			argDomain := args[0]
			argDomainName := args[1]

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

			pubKey, err := info.GetPubKey()
			if err != nil {
				return err
			}
			pubKeyHex := utils.GetPubKeyHex(pubKey)
			pubKeyType, err := utils.DerivePubKeyType(pubKey)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateDomainName(
				signer.String(),
				argDomain,
				argDomainName,
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
