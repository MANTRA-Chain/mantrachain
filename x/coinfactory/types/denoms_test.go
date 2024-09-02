package types_test

import (
	"testing"

	"github.com/MANTRA-Finance/mantrachain/x/coinfactory/types"
	"github.com/stretchr/testify/require"
)

func TestDeconstructDenom(t *testing.T) {
	for _, tc := range []struct {
		desc             string
		denom            string
		expectedSubdenom string
		err              error
	}{
		{
			desc:  "empty is invalid",
			denom: "",
			err:   types.ErrInvalidDenom,
		},
		{
			desc:             "normal",
			denom:            "factory/cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw/bitcoin",
			expectedSubdenom: "bitcoin",
		},
		{
			desc:             "multiple slashes in subdenom",
			denom:            "factory/cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw/bitcoin/1",
			expectedSubdenom: "bitcoin/1",
		},
		{
			desc:             "no subdenom",
			denom:            "factory/cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw/",
			expectedSubdenom: "",
		},
		{
			desc:  "incorrect prefix",
			denom: "ibc/cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw/bitcoin",
			err:   types.ErrInvalidDenom,
		},
		{
			desc:             "subdenom of only slashes",
			denom:            "factory/cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw/////",
			expectedSubdenom: "////",
		},
		{
			desc:  "too long name",
			denom: "factory/cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw/adsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsf",
			err:   types.ErrInvalidDenom,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			expectedCreator := "cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw"
			creator, subdenom, err := types.DeconstructDenom(tc.denom)
			if tc.err != nil {
				require.ErrorContains(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, expectedCreator, creator)
				require.Equal(t, tc.expectedSubdenom, subdenom)
			}
		})
	}
}

func TestGetTokenDenom(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		creator  string
		subdenom string
		valid    bool
	}{
		{
			desc:     "normal",
			creator:  "cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw",
			subdenom: "bitcoin",
			valid:    true,
		},
		{
			desc:     "multiple slashes in subdenom",
			creator:  "cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw",
			subdenom: "bitcoin/1",
			valid:    true,
		},
		{
			desc:     "no subdenom",
			creator:  "cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw",
			subdenom: "",
			valid:    true,
		},
		{
			desc:     "subdenom of only slashes",
			creator:  "cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw",
			subdenom: "/////",
			valid:    true,
		},
		{
			desc:     "too long name",
			creator:  "cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw",
			subdenom: "adsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsf",
			valid:    false,
		},
		{
			desc:     "subdenom is exactly max length",
			creator:  "cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw",
			subdenom: "bitcoinfsadfsdfeadfsafwefsefsefsdfsdafasefsf",
			valid:    true,
		},
		{
			desc:     "creator is exactly max length",
			creator:  "cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvwjhgjhgkhjklhkjhkjhgjhgjgjghelu",
			subdenom: "bitcoin",
			valid:    true,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := types.GetTokenDenom(tc.creator, tc.subdenom)
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
