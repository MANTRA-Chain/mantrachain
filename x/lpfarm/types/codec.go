package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	// this line is used by starport scaffolding # 1
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreatePrivatePlan{}, "lpfarm/MsgCreatePrivatePlan", nil)
	cdc.RegisterConcrete(&MsgTerminatePrivatePlan{}, "lpfarm/MsgTerminatePrivatePlan", nil)
	cdc.RegisterConcrete(&MsgFarm{}, "lpfarm/MsgFarm", nil)
	cdc.RegisterConcrete(&MsgUnfarm{}, "lpfarm/MsgUnfarm", nil)
	cdc.RegisterConcrete(&MsgHarvest{}, "lpfarm/MsgHarvest", nil)
	cdc.RegisterConcrete(&FarmingPlanProposal{}, "lpfarm/FarmingPlanProposal", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3

	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},

		&MsgCreatePrivatePlan{},
		&MsgTerminatePrivatePlan{},
		&MsgFarm{},
		&MsgUnfarm{},
		&MsgHarvest{},
	)

	registry.RegisterImplementations(
		(*govv1beta1.Content)(nil),
		&FarmingPlanProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
