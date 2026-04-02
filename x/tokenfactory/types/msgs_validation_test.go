package types

import (
	"testing"

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
