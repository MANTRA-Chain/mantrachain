package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateCampaign{}, "airdrop/CreateCampaign", nil)
	cdc.RegisterConcrete(&MsgDeleteCampaign{}, "airdrop/DeleteCampaign", nil)
	cdc.RegisterConcrete(&MsgPauseCampaign{}, "airdrop/PauseCampaign", nil)
	cdc.RegisterConcrete(&MsgUnpauseCampaign{}, "airdrop/UnpauseCampaign", nil)
	cdc.RegisterConcrete(&MsgCampaignClaim{}, "airdrop/CampaignClaim", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateCampaign{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgDeleteCampaign{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgPauseCampaign{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUnpauseCampaign{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCampaignClaim{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
