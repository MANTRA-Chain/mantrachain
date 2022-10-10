package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"

	// this line is used by starport scaffolding # 1
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgRegisterBridge{}, "bridge/RegisterBridge", nil)
	cdc.RegisterConcrete(&MsgMint{}, "bridge/Mint", nil)
	cdc.RegisterConcrete(&MsgCreateCw20Contract{}, "bridge/CreateCw20Contract", nil)
cdc.RegisterConcrete(&MsgUpdateCw20Contract{}, "bridge/UpdateCw20Contract", nil)
cdc.RegisterConcrete(&MsgDeleteCw20Contract{}, "bridge/DeleteCw20Contract", nil)
// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRegisterBridge{},
		&MsgMint{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
	&MsgCreateCw20Contract{},
	&MsgUpdateCw20Contract{},
	&MsgDeleteCw20Contract{},
)
// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
