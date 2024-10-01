package app

import (
	errorsmod "cosmossdk.io/errors"
	xfeemarketpost "github.com/MANTRA-Chain/mantrachain/x/xfeemarket/post"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// PostHandlerOptions are the options required for constructing a FeeMarket PostHandler.
type PostHandlerOptions struct {
	BankKeeper      xfeemarketpost.BankKeeper
	FeeMarketKeeper xfeemarketpost.FeeMarketKeeper
}

// NewPostHandler returns a PostHandler chain with the fee deduct decorator.
func NewPostHandler(options PostHandlerOptions) (sdk.PostHandler, error) {
	if options.BankKeeper == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "bank keeper is required for post builder")
	}

	if options.FeeMarketKeeper == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "feemarket keeper is required for post builder")
	}

	postDecorators := []sdk.PostDecorator{
		xfeemarketpost.NewFeeMarketDeductDecorator(
			options.BankKeeper,
			options.FeeMarketKeeper,
		),
	}

	return sdk.ChainPostDecorators(postDecorators...), nil
}
