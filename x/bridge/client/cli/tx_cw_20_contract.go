package cli

import (
	"github.com/spf13/cast"

	"github.com/LimeChain/mantrachain/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdCreateCw20Contract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-cw-20-contract [store-id] [ver] [path]",
		Short: "Create Cw20Contract",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argStoreId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			argVer := args[1]
			argPath := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateCw20Contract(clientCtx.GetFromAddress().String(), argStoreId, argVer, argPath)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdUpdateCw20Contract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-cw-20-contract [store-id] [ver] [path]",
		Short: "Update Cw20Contract",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argStoreId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			argVer := args[1]
			argPath := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateCw20Contract(clientCtx.GetFromAddress().String(), argStoreId, argVer, argPath)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdDeleteCw20Contract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-cw-20-contract",
		Short: "Delete Cw20Contract",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteCw20Contract(clientCtx.GetFromAddress().String())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
