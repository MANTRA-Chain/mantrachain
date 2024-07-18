package types_test

import (
	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var testAddr = sdk.AccAddress(crypto.AddressHash([]byte("test")))

func newInt(i int64) math.Int {
	return math.NewInt(i)
}
