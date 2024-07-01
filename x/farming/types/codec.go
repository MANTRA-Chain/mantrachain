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
	cdc.RegisterInterface((*PlanI)(nil), nil)
	cdc.RegisterConcrete(&MsgCreateFixedAmountPlan{}, "farming/MsgCreateFixedAmountPlan", nil)
	cdc.RegisterConcrete(&MsgCreateRatioPlan{}, "farming/MsgCreateRatioPlan", nil)
	cdc.RegisterConcrete(&MsgStake{}, "farming/MsgStake", nil)
	cdc.RegisterConcrete(&MsgUnstake{}, "farming/MsgUnstake", nil)
	cdc.RegisterConcrete(&MsgHarvest{}, "farming/MsgHarvest", nil)
	cdc.RegisterConcrete(&MsgRemovePlan{}, "farming/MsgRemovePlan", nil)
	cdc.RegisterConcrete(&FixedAmountPlan{}, "farming/FixedAmountPlan", nil)
	cdc.RegisterConcrete(&RatioPlan{}, "farming/RatioPlan", nil)
	cdc.RegisterConcrete(&PublicPlanProposal{}, "farming/PublicPlanProposal", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3

	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},

		&MsgCreateFixedAmountPlan{},
		&MsgCreateRatioPlan{},
		&MsgStake{},
		&MsgUnstake{},
		&MsgHarvest{},
		&MsgRemovePlan{},
	)

	registry.RegisterImplementations(
		(*govv1beta1.Content)(nil),
		&PublicPlanProposal{},
	)

	registry.RegisterInterface(
		"cosmos.farming.v1beta1.PlanI",
		(*PlanI)(nil),
		&FixedAmountPlan{},
		&RatioPlan{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
