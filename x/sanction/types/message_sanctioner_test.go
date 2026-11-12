package types

import (
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func randomSanctionAddress() string {
	pk := secp256k1.GenPrivKey().PubKey()
	return sdk.AccAddress(pk.Address()).String()
}

func TestMsgSanctionerValidate(t *testing.T) {
	signer := randomSanctionAddress()
	target := randomSanctionAddress()

	tests := []struct {
		name    string
		msg     interface{ Validate() error }
		wantErr bool
	}{
		{"sanction valid", NewMsgSanctionAddresses(signer, []string{target}, "rule 7, tx ABC"), false},
		{"sanction empty reason", NewMsgSanctionAddresses(signer, []string{target}, "  "), true},
		{"sanction no addresses", NewMsgSanctionAddresses(signer, nil, "r"), true},
		{"sanction bad address", NewMsgSanctionAddresses(signer, []string{"nope"}, "r"), true},
		{"sanction bad signer", NewMsgSanctionAddresses("nope", []string{target}, "r"), true},
		{"unsanction valid", NewMsgUnsanctionAddresses(signer, []string{target}, "false positive"), false},
		{"sanction whitespace-only reason", NewMsgSanctionAddresses(signer, []string{target}, "\t\n"), true},
		{"sanction one bad address among good", NewMsgSanctionAddresses(signer, []string{target, "nope"}, "r"), true},
		{"sanction mixed-case address", NewMsgSanctionAddresses(signer, []string{strings.ToUpper(target)}, "r"), false},
		{"unsanction empty reason", NewMsgUnsanctionAddresses(signer, []string{target}, ""), true},
		{"unsanction no addresses", NewMsgUnsanctionAddresses(signer, []string{}, "r"), true},
		{"unsanction bad address", NewMsgUnsanctionAddresses(signer, []string{"nope"}, "r"), true},
		{"unsanction bad signer", NewMsgUnsanctionAddresses("", []string{target}, "r"), true},
		{"revoke valid", NewMsgRevokeSanctioner(signer, target, "key leaked"), false},
		{"revoke bad signer", NewMsgRevokeSanctioner("nope", target, "r"), true},
		{"revoke bad target", NewMsgRevokeSanctioner(signer, "nope", "r"), true},
		{"revoke empty reason", NewMsgRevokeSanctioner(signer, target, ""), true},
		{"update params valid", NewMsgUpdateParams(signer, DefaultParams()), false},
		{"update params zero cap", NewMsgUpdateParams(signer, NewParams(0, 10)), true},
		{"update params zero reason length", NewMsgUpdateParams(signer, NewParams(10, 0)), true},
		{"update params bad authority", NewMsgUpdateParams("nope", DefaultParams()), true},
		{"update sanctioners empty list", NewMsgUpdateSanctioners(signer, nil), false},
		{"update sanctioners valid", NewMsgUpdateSanctioners(signer, []Sanctioner{
			{Address: target, Role: SanctionerRole_SANCTIONER_ROLE_FULL},
			{Address: signer, Role: SanctionerRole_SANCTIONER_ROLE_ADD_ONLY},
		}), false},
		{"update sanctioners duplicate", NewMsgUpdateSanctioners(signer, []Sanctioner{
			{Address: target, Role: SanctionerRole_SANCTIONER_ROLE_FULL},
			{Address: target, Role: SanctionerRole_SANCTIONER_ROLE_ADD_ONLY},
		}), true},
		{"update sanctioners duplicate differing case", NewMsgUpdateSanctioners(signer, []Sanctioner{
			{Address: target, Role: SanctionerRole_SANCTIONER_ROLE_FULL},
			{Address: strings.ToUpper(target), Role: SanctionerRole_SANCTIONER_ROLE_ADD_ONLY},
		}), true},
		{"update sanctioners unspecified role", NewMsgUpdateSanctioners(signer, []Sanctioner{{Address: target}}), true},
		{"update sanctioners bad authority", NewMsgUpdateSanctioners("nope", nil), true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.Validate()
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
