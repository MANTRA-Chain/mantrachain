package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreatePair{}, "liquidity/MsgCreatePair", nil)
	cdc.RegisterConcrete(&MsgUpdatePairSwapFee{}, "liquidity/MsgUpdatePairSwapFee", nil)
	cdc.RegisterConcrete(&MsgCreatePool{}, "liquidity/MsgCreatePool", nil)
	cdc.RegisterConcrete(&MsgCreateRangedPool{}, "liquidity/MsgCreateRangedPool", nil)
	cdc.RegisterConcrete(&MsgDeposit{}, "liquidity/MsgDeposit", nil)
	cdc.RegisterConcrete(&MsgWithdraw{}, "liquidity/MsgWithdraw", nil)
	cdc.RegisterConcrete(&MsgLimitOrder{}, "liquidity/MsgLimitOrder", nil)
	cdc.RegisterConcrete(&MsgMarketOrder{}, "liquidity/MsgMarketOrder", nil)
	cdc.RegisterConcrete(&MsgMMOrder{}, "liquidity/MsgMMOrder", nil)
	cdc.RegisterConcrete(&MsgCancelOrder{}, "liquidity/MsgCancelOrder", nil)
	cdc.RegisterConcrete(&MsgCancelAllOrders{}, "liquidity/MsgCancelAllOrders", nil)
	cdc.RegisterConcrete(&MsgCancelMMOrder{}, "liquidity/MsgCancelMMOrder", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3

	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},

		&MsgCreatePair{},
		&MsgUpdatePairSwapFee{},
		&MsgCreatePool{},
		&MsgCreateRangedPool{},
		&MsgDeposit{},
		&MsgWithdraw{},
		&MsgLimitOrder{},
		&MsgMarketOrder{},
		&MsgMMOrder{},
		&MsgCancelOrder{},
		&MsgCancelAllOrders{},
		&MsgCancelMMOrder{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
