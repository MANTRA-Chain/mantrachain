package cli

import (
    "strconv"
	
	"github.com/spf13/cobra"
    "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/LimeChain/mantrachain/x/mns/types"
)

var _ = strconv.Itoa(0)

func CmdCreateDomainName() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-domain-name [domain] [name]",
		Short: "Broadcast message create-domain-name",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
      		 argDomain := args[0]
             argName := args[1]
            
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateDomainName(
				clientCtx.GetFromAddress().String(),
				argDomain,
				argName,
				
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