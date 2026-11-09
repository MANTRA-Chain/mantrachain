package types

import (
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
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
