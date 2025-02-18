package types

import (
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddBlacklistAccount{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRemoveBlacklistAccount{},
	)
	// this line is used by starport scaffolding # 3
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
