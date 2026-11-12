package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

// EVMKeeper defines the x/vm keeper methods x/sanction uses to clear an
// EIP-7702 delegation from a sanctioned account.
type EVMKeeper interface {
	GetCodeHash(ctx sdk.Context, addr common.Address) common.Hash
	GetCode(ctx sdk.Context, codeHash common.Hash) []byte
	DeleteCodeHash(ctx sdk.Context, addr common.Address)
}
