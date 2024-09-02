package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgApplyMarketMaker{}, "marketmaker/MsgApplyMarketMaker", nil)
	cdc.RegisterConcrete(&MsgClaimIncentives{}, "marketmaker/MsgClaimIncentives", nil)
	cdc.RegisterConcrete(&MarketMakerProposal{}, "marketmaker/MarketMakerProposal", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3

	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},

		&MsgApplyMarketMaker{},
		&MsgClaimIncentives{},
	)

	registry.RegisterImplementations(
		(*govv1beta1.Content)(nil),
		&MarketMakerProposal{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
