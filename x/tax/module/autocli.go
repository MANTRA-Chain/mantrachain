package tax

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	taxapi "github.com/MANTRA-Chain/mantrachain/api/mantrachain/tax/v1"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: taxapi.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              taxapi.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Use:       "update-params",
					Skip:      false,
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"mca_tax": {
							Usage:        "mca tax for the allocation in decimal",
							DefaultValue: "",
						},
						"mca_address": {
							Usage:        "mca address for the allocation",
							DefaultValue: "",
						},
					},
					Short:   "Update the parameters of the tax module",
					Example: "mantrachaind tx tax update-params --mca_tax 0.4--mca_address mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka",
				},
			},
		},
	}
}
