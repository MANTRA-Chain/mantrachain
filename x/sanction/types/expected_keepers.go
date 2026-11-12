package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
)

// AccountKeeper defines the x/auth keeper methods x/sanction uses to tell
// module accounts apart, so they're never sanctioned.
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	GetModulePermissions() map[string]authtypes.PermissionsForAddress
}

// EVMKeeper defines the x/vm keeper methods x/sanction uses to clear an
// EIP-7702 delegation from a sanctioned account.
type EVMKeeper interface {
	GetCodeHash(ctx sdk.Context, addr common.Address) common.Hash
	GetCode(ctx sdk.Context, codeHash common.Hash) []byte
	DeleteCodeHash(ctx sdk.Context, addr common.Address)
}

// TokenFactoryKeeper defines the x/tokenfactory keeper method x/sanction uses
// to recognise the escrow accounts of open IBC transfer channels.
type TokenFactoryKeeper interface {
	IsEscrowAddress(ctx sdk.Context, address sdk.AccAddress) bool
}
