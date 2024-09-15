package post

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	feemarkettypes "github.com/skip-mev/feemarket/x/feemarket/types"
)

// BankKeeper defines the contract needed for supply related APIs.
//
//go:generate mockery --name BankKeeper --filename mock_bank_keeper.go
type BankKeeper interface {
	bankkeeper.Keeper
}

// FeeMarketKeeper defines the expected feemarket keeper.
//
//go:generate mockery --name FeeMarketKeeper --filename mock_feemarket_keeper.go
type FeeMarketKeeper interface {
	GetState(ctx sdk.Context) (feemarkettypes.State, error)
	GetParams(ctx sdk.Context) (feemarkettypes.Params, error)
	SetParams(ctx sdk.Context, params feemarkettypes.Params) error
	SetState(ctx sdk.Context, state feemarkettypes.State) error
	ResolveToDenom(ctx sdk.Context, coin sdk.DecCoin, denom string) (sdk.DecCoin, error)
	GetMinGasPrice(ctx sdk.Context, denom string) (sdk.DecCoin, error)
	GetEnabledHeight(ctx sdk.Context) (int64, error)
}
