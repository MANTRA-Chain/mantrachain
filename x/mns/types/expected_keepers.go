package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DidKeeper interface {
	SetNewDidDocument(ctx sdk.Context, id string, signer sdk.Address, pubKeyHex string, pubKeyType string) (string, error)
}
