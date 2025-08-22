package slinky_test

import (
	enccodec "github.com/cosmos/evm/encoding/codec"

	amino "github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/testutil"
	"github.com/cosmos/cosmos-sdk/types/module"
	sdktestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
)

// MakeConfig creates a new EncodingConfig and returns it
func MakeTestEncodingConfig(modules ...module.AppModuleBasic) sdktestutil.TestEncodingConfig {
	cdc := amino.NewLegacyAmino()
	interfaceRegistry := testutil.CodecOptions{}.NewInterfaceRegistry()
	mb := module.NewBasicManager(modules...)
	mb.RegisterInterfaces(interfaceRegistry)
	codec := amino.NewProtoCodec(interfaceRegistry)
	enccodec.RegisterLegacyAminoCodec(cdc)
	enccodec.RegisterInterfaces(interfaceRegistry)

	return sdktestutil.TestEncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Codec:             codec,
		TxConfig:          tx.NewTxConfig(codec, tx.DefaultSignModes),
		Amino:             cdc,
	}
}
