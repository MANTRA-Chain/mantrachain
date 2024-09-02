package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateDidDocument{}, "did/create-did", nil)
	cdc.RegisterConcrete(&MsgUpdateDidDocument{}, "did/update-did", nil)
	cdc.RegisterConcrete(&MsgAddVerification{}, "did/add-verification", nil)
	cdc.RegisterConcrete(&MsgSetVerificationRelationships{}, "did/set-verification-relationship", nil)
	cdc.RegisterConcrete(&MsgRevokeVerification{}, "did/revoke-verification", nil)
	cdc.RegisterConcrete(&MsgAddService{}, "did/add-service", nil)
	cdc.RegisterConcrete(&MsgDeleteService{}, "did/delete-service", nil)
	cdc.RegisterConcrete(&MsgAddController{}, "did/add-controller", nil)
	cdc.RegisterConcrete(&MsgDeleteController{}, "did/delete-controller", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3

	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},

		&MsgCreateDidDocument{},
		&MsgUpdateDidDocument{},
		&MsgAddVerification{},
		&MsgSetVerificationRelationships{},
		&MsgRevokeVerification{},
		&MsgAddService{},
		&MsgDeleteService{},
		&MsgAddController{},
		&MsgDeleteController{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
