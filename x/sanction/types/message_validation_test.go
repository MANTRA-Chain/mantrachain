package types

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func randomSanctionAddress() string {
	pk := secp256k1.GenPrivKey().PubKey()
	return sdk.AccAddress(pk.Address()).String()
}

func randomSanctionAddresses(n int) []string {
	addrs := make([]string, n)
	for i := 0; i < n; i++ {
		addrs[i] = randomSanctionAddress()
	}
	return addrs
}

type blacklistValidateTestcase struct {
	name      string
	authority string
	accounts  []string
	wantErr   bool
}

func baseBlacklistValidateTestcases() []blacklistValidateTestcase {
	return []blacklistValidateTestcase{
		{
			name:      "valid authority and valid accounts",
			authority: randomSanctionAddress(),
			accounts:  []string{randomSanctionAddress()},
			wantErr:   false,
		},
		{
			name:      "invalid authority rejected",
			authority: "not_bech32",
			accounts:  []string{randomSanctionAddress()},
			wantErr:   true,
		},
		{
			name:      "empty authority rejected",
			authority: "",
			accounts:  []string{randomSanctionAddress()},
			wantErr:   true,
		},
		{
			name:      "empty accounts rejected",
			authority: randomSanctionAddress(),
			accounts:  []string{},
			wantErr:   true,
		},
		{
			name:      "invalid account rejected",
			authority: randomSanctionAddress(),
			accounts:  []string{"not_bech32"},
			wantErr:   true,
		},
	}
}

func runBlacklistValidateTestcases(t *testing.T, testcases []blacklistValidateTestcase, validate func(authority string, accounts []string) error) {
	t.Helper()

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := validate(tc.authority, tc.accounts)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgAddBlacklistAccounts(t *testing.T) {
	testcases := append(baseBlacklistValidateTestcases(), blacklistValidateTestcase{
		name:      "more than 100 accounts rejected",
		authority: randomSanctionAddress(),
		accounts:  randomSanctionAddresses(101),
		wantErr:   true,
	})

	runBlacklistValidateTestcases(t, testcases, func(authority string, accounts []string) error {
		return NewMsgAddBlacklistAccounts(authority, accounts).Validate()
	})
}

func TestMsgRemoveBlacklistAccounts(t *testing.T) {
	runBlacklistValidateTestcases(t, baseBlacklistValidateTestcases(), func(authority string, accounts []string) error {
		return NewMsgRemoveBlacklistAccounts(authority, accounts).Validate()
	})
}
