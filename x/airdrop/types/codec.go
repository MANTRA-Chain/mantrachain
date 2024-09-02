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
	// this line is used by starport scaffolding # 3

	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},

		&MsgCreateCampaign{},
		&MsgDeleteCampaign{},
		&MsgPauseCampaign{},
		&MsgUnpauseCampaign{},
		&MsgCampaignClaim{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
