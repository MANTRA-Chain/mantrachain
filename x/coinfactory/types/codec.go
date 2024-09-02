package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateDenom{}, "coinfactory/create-denom", nil)
	cdc.RegisterConcrete(&MsgMint{}, "coinfactory/mint", nil)
	cdc.RegisterConcrete(&MsgBurn{}, "coinfactory/burn", nil)
	cdc.RegisterConcrete(&MsgForceTransfer{}, "coinfactory/force-transfer", nil)
	cdc.RegisterConcrete(&MsgChangeAdmin{}, "coinfactory/change-admin", nil)
	cdc.RegisterConcrete(&MsgSetDenomMetadata{}, "coinfactory/set-denom-metadata", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3

	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},

		&MsgCreateDenom{},
		&MsgMint{},
		&MsgBurn{},
		&MsgForceTransfer{},
		&MsgChangeAdmin{},
		&MsgSetDenomMetadata{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
