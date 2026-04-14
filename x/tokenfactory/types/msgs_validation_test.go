package types

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func randomAccAddress() string {
	pk := secp256k1.GenPrivKey().PubKey()
	return sdk.AccAddress(pk.Address()).String()
}

func TestMsgForceTransfer_ValidateGroupedTestcases(t *testing.T) {
	sender := randomAccAddress()
	from := randomAccAddress()
	to := randomAccAddress()

	testcases := []struct {
		name    string
		amount  sdk.Coin
		wantErr bool
	}{
		{
			name:    "valid positive amount",
			amount:  sdk.NewCoin("factory/"+sender+"/subdenom", math.NewInt(1)),
			wantErr: false,
		},
		{
			name:    "zero amount rejected",
			amount:  sdk.NewCoin("factory/"+sender+"/subdenom", math.ZeroInt()),
			wantErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := NewMsgForceTransfer(sender, tc.amount, from, to).Validate()
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
