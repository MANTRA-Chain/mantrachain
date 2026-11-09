package types

import (
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// MsgAddBlacklistAccounts and MsgRemoveBlacklistAccounts have no handler.
	// They stay registered so past proposals that contain them still decode.
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddBlacklistAccounts{},
		&MsgRemoveBlacklistAccounts{},
		&MsgUpdateParams{},
		&MsgUpdateSanctioners{},
		&MsgSanctionAddresses{},
		&MsgUnsanctionAddresses{},
		&MsgRevokeSanctioner{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
