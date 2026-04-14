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

func malformedFactoryDenoms(sender string) []string {
	return []string{
		"factory/" + sender + "/",
		"factory/" + sender + "//",
		"factory/" + sender + "/a//b",
		"factory/" + sender + "/..",
	}
}

func TestMsgValidate_GroupedTestcases(t *testing.T) {
	tests := []struct {
		name     string
		validate func(sender, denom string) error
	}{
		{
			name: "change admin",
			validate: func(sender, denom string) error {
				newAdmin := randomAccAddress()
				return NewMsgChangeAdmin(sender, denom, newAdmin).Validate()
			},
		},
		{
			name: "set before send hook",
			validate: func(sender, denom string) error {
				contract := randomAccAddress()
				return NewMsgSetBeforeSendHook(sender, denom, contract).Validate()
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sender := randomAccAddress()
			testcases := []struct {
				name    string
				denom   string
				wantErr bool
			}{
				{name: "valid denom", denom: "factory/" + sender + "/subdenom/path", wantErr: false},
			}

			for _, denom := range malformedFactoryDenoms(sender) {
				testcases = append(testcases, struct {
					name    string
					denom   string
					wantErr bool
				}{
					name:    "invalid denom: " + denom,
					denom:   denom,
					wantErr: true,
				})
			}

			for _, testcase := range testcases {
				t.Run(testcase.name, func(t *testing.T) {
					err := tc.validate(sender, testcase.denom)
					if testcase.wantErr {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
				})
			}
		})
	}
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
