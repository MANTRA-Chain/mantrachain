package amm

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	utils "github.com/MANTRA-Finance/mantrachain/types"
)

func Test_char(t *testing.T) {
	require.Panics(t, func() {
		char(math.LegacyZeroDec())
	})

	for _, tc := range []struct {
		x        math.LegacyDec
		expected int
	}{
		{math.LegacyMustNewDecFromStr("999.99999999999999999"), 20},
		{math.LegacyMustNewDecFromStr("100"), 20},
		{math.LegacyMustNewDecFromStr("99.999999999999999999"), 19},
		{math.LegacyMustNewDecFromStr("10"), 19},
		{math.LegacyMustNewDecFromStr("9.999999999999999999"), 18},
		{math.LegacyMustNewDecFromStr("1"), 18},
		{math.LegacyMustNewDecFromStr("0.999999999999999999"), 17},
		{math.LegacyMustNewDecFromStr("0.1"), 17},
		{math.LegacyMustNewDecFromStr("0.099999999999999999"), 16},
		{math.LegacyMustNewDecFromStr("0.01"), 16},
		{math.LegacyMustNewDecFromStr("0.000000000000000009"), 0},
		{math.LegacyMustNewDecFromStr("0.000000000000000001"), 0},
	} {
		t.Run("", func(t *testing.T) {
			require.Equal(t, tc.expected, char(tc.x))
		})
	}
}

func Test_pow10(t *testing.T) {
	for _, tc := range []struct {
		power    int
		expected math.LegacyDec
	}{
		{18, math.LegacyNewDec(1)},
		{19, math.LegacyNewDec(10)},
		{20, math.LegacyNewDec(100)},
		{17, math.LegacyNewDecWithPrec(1, 1)},
		{16, math.LegacyNewDecWithPrec(1, 2)},
	} {
		t.Run("", func(t *testing.T) {
			require.True(math.LegacyDecEq(t, tc.expected, pow10(tc.power)))
		})
	}
}

func Test_isPow10(t *testing.T) {
	for _, tc := range []struct {
		x        math.LegacyDec
		expected bool
	}{
		{utils.ParseDec("100"), true},
		{utils.ParseDec("101"), false},
		{utils.ParseDec("10"), true},
		{utils.ParseDec("1"), true},
		{utils.ParseDec("1.000000000000000001"), false},
		{utils.ParseDec("0.11"), false},
		{utils.ParseDec("0.000000000000000001"), true},
		{utils.ParseDec("10000000000000000000000000001"), false},
		{utils.ParseDec("10000000000000000000000000000"), true},
		{utils.ParseDec("123456789"), false},
	} {
		t.Run("", func(t *testing.T) {
			require.Equal(t, tc.expected, isPow10(tc.x))
		})
	}
}
