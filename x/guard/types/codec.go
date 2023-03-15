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
	cdc.RegisterConcrete(&MsgUpdateRequiredPrivileges{}, "guard/UpdateRequiredPrivileges", nil)
	cdc.RegisterConcrete(&MsgUpdateRequiredPrivilegesBatch{}, "guard/UpdateRequiredPrivilegesBatch", nil)
	cdc.RegisterConcrete(&MsgUpdateRequiredPrivilegesGroupedBatch{}, "guard/UpdateRequiredPrivilegesGroupedBatch", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateAccountPrivileges{},
		&MsgUpdateAccountPrivilegesBatch{},
		&MsgUpdateAccountPrivilegesGroupedBatch{},
		&MsgUpdateRequiredPrivileges{},
		&MsgUpdateRequiredPrivilegesBatch{},
		&MsgUpdateRequiredPrivilegesGroupedBatch{},
		&MsgUpdateGuardTransferCoins{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
