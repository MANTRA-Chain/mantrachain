package app

import (
	"encoding/base64"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	evmosencoding "github.com/cosmos/evm/encoding"
	"github.com/stretchr/testify/require"
)

func TestDecodeLegacyTx(t *testing.T) {
	txs := map[string]string{
		"/feemarket.feemarket.v1.MsgParams": "CrICCq8CCiAvY29zbW9zLmdvdi52MS5Nc2dTdWJtaXRQcm9wb3NhbBKKAgq+AQohL2ZlZW1hcmtldC5mZWVtYXJrZXQudjEuTXNnUGFyYW1zEpgBCmcKATASEzEwMDAwMDAwMDAwMDAwMDAwMDAaATAiATAqETEwMDAwMDAwMDAwMDAwMDAwMhIxMjUwMDAwMDAwMDAwMDAwMDA6EjEyNTAwMDAwMDAwMDAwMDAwMEDA0eEjSAFSA3VvbVgBEi1tYW50cmExMGQwN3kyNjVnbW11dnQ0ejB3OWF3ODgwam5zcjcwMGozZmVwNGYSCAoDdW9tEgExGi1tYW50cmExZWNtZjdka2E2MzU0cTJtanE3YWprNXpyOTB6aG1qenRlZ2owZTcqBXRpdGxlMgdzdW1tYXJ5EmUKTgpGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQLoSJJS0cs48ErL6k1q/9VNwaL0K89lcvDMuZInqTdbhRIECgIIARITCg0KA3VvbRIGMjUwMDAwEJChDxpAo5hkwyggn4yZYZarwbA1KorHHGhyhjrPcvIxrUfT4mlO1YZMPYBbVJkhpcHchUWZKBBjTK6frLDoLzdttlbbKA==",
	}

	for typeUrl, tx := range txs {
		bz, err := base64.StdEncoding.DecodeString(tx)
		require.NoError(t, err)

		encodingConfig := evmosencoding.MakeConfig(MANTRAChainID)
		mb := module.NewBasicManager(gov.AppModuleBasic{})
		mb.RegisterLegacyAminoCodec(encodingConfig.Amino)
		mb.RegisterInterfaces(encodingConfig.InterfaceRegistry)

		_, err = encodingConfig.TxConfig.TxDecoder()(bz)
		require.Error(t, err, "")
		require.Contains(t, err.Error(), "unable to resolve type URL "+typeUrl)

		RegisterLegacyCodec(encodingConfig.Amino)
		RegisterLegacyInterfaces(encodingConfig.InterfaceRegistry)
		banktypes.RegisterLegacyAminoCodec(encodingConfig.Amino)
		banktypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)

		_, err = encodingConfig.TxConfig.TxDecoder()(bz)
		require.NoError(t, err)
	}
}
