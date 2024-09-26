package xfeemarket

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	xfeemarketapi "github.com/MANTRA-Chain/mantrachain/api/mantrachain/xfeemarket/v1"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: xfeemarketapi.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              xfeemarketapi.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "UpsertFeeDenom",
					Use:            "upsert-fee-denom [denom] [multiplier]",
					Short:          "Send a upsert-fee-denom tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "denom"}, {ProtoField: "multiplier"}},
				},
				{
					RpcMethod:      "RemoveFeeDenom",
					Use:            "remove-fee-denom [denom]",
					Short:          "Send a remove-fee-denom tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "denom"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
