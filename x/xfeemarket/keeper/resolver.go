package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	feemarkettypes "github.com/skip-mev/feemarket/x/feemarket/types"
)

var _ feemarkettypes.DenomResolver = &Keeper{}

// ConvertToDenom returns "coin.Amount denom" for all coins that are not the denom.
func (k Keeper) ConvertToDenom(ctx sdk.Context, coin sdk.DecCoin, denom string) (sdk.DecCoin, error) {
	if coin.Denom == denom {
		return coin, nil
	}

	multiplier, err := k.DenomMultipliers.Get(ctx, denom)
	if err != nil {
		return sdk.DecCoin{}, err
	}
	amount := coin.Amount.Mul(multiplier.Dec)

	return sdk.NewDecCoinFromDec(denom, amount), nil
}

func (k Keeper) ExtraDenoms(ctx sdk.Context) ([]string, error) {
	iter, err := k.DenomMultipliers.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer iter.Close()
	extraDenoms, err := iter.Keys()
	if err != nil {
		return nil, err
	}
	return extraDenoms, nil
}
