package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateAccPerm{}, "guard/CreateAccPerm", nil)
	cdc.RegisterConcrete(&MsgUpdateAccPerm{}, "guard/UpdateAccPerm", nil)
	cdc.RegisterConcrete(&MsgDeleteAccPerm{}, "guard/DeleteAccPerm", nil)
	cdc.RegisterConcrete(&MsgUpdateGuardTransfer{}, "guard/UpdateGuardTransfer", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateAccPerm{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateGuardTransfer{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
