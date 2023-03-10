package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgUpdateAccountPrivileges{}, "guard/UpdateAccountPrivileges", nil)
	cdc.RegisterConcrete(&MsgUpdateAccountPrivilegesBatch{}, "guard/UpdateAccountPrivilegesBatch", nil)
	cdc.RegisterConcrete(&MsgUpdateAccountPrivilegesGroupedBatch{}, "guard/UpdateAccountPrivilegesGroupedBatch", nil)
	cdc.RegisterConcrete(&MsgUpdateGuardTransferCoins{}, "guard/UpdateGuardTransferCoins", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateAccountPrivileges{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateAccountPrivilegesBatch{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateAccountPrivilegesGroupedBatch{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateGuardTransferCoins{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
