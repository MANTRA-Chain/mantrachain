package cli

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/LimeChain/mantrachain/x/mns/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
)

// deriveVMType derive the verification method type from a public key
func deriveVMType(pubKey cryptotypes.PubKey) (vmType string, err error) {
	switch pubKey.(type) {
	case *secp256k1.PubKey:
		vmType = types.EcdsaSecp256k1VerificationKey2019
	default:
		err = types.ErrKeyFormatNotSupported
	}
	return
}

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
			pubKey := info.GetPubKey()
			pubKeyHex := strings.ToUpper(fmt.Sprint("F", hex.EncodeToString(pubKey.Bytes())))

			// understand the vmType
			vmType, err := deriveVMType(pubKey)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateDomainName(
				signer.String(),
				argDomain,
				argDomainName,
				pubKeyHex,
				vmType,
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
