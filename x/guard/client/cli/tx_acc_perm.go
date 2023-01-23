package cli

import (
	
	 "strings"
    "github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/LimeChain/mantrachain/x/guard/types"
)

func CmdCreateAccPerm() *cobra.Command {
    cmd := &cobra.Command{
		Use:   "create-acc-perm [cat] [whl-curr]",
		Short: "Create a new acc_perm",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
            // Get indexes
         indexCat := args[0]
        
            // Get value arguments
		 argWhlCurr := strings.Split(args[1], listSeparator)
		
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateAccPerm(
			    clientCtx.GetFromAddress().String(),
			    indexCat,
                argWhlCurr,
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

func CmdUpdateAccPerm() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-acc-perm [cat] [whl-curr]",
		Short: "Update a acc_perm",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
            // Get indexes
         indexCat := args[0]
        
            // Get value arguments
		 argWhlCurr := strings.Split(args[1], listSeparator)
		
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateAccPerm(
			    clientCtx.GetFromAddress().String(),
			    indexCat,
                argWhlCurr,
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

func CmdDeleteAccPerm() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-acc-perm [cat]",
		Short: "Delete a acc_perm",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
             indexCat := args[0]
            
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteAccPerm(
			    clientCtx.GetFromAddress().String(),
			    indexCat,
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