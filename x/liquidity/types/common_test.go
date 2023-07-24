package types_test

import (
	"github.com/cometbft/cometbft/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var testAddr = sdk.AccAddress(crypto.AddressHash([]byte("test")))

func newInt(i int64) sdk.Int {
	return sdk.NewInt(i)
}
