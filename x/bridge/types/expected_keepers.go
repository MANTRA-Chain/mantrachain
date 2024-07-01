package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	InputOutputCoins(ctx context.Context, input banktypes.Input, outputs []banktypes.Output) error
	// Methods imported from bank should be defined here
}

type GuardKeeper interface {
	CheckHasAuthz(ctx sdk.Context, address string, authz string) error
}
