package app

import (
	feemarkettypes "github.com/MANTRA-Chain/mantrachain/v5/app/legacy/feemarket/types"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
)

func RegisterLegacyCodec(cdc *codec.LegacyAmino) {
	feemarkettypes.RegisterCodec(cdc)
}

func RegisterLegacyInterfaces(registry cdctypes.InterfaceRegistry) {
	feemarkettypes.RegisterInterfaces(registry)
}
