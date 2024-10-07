package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	feemarkettypes "github.com/skip-mev/feemarket/x/feemarket/types"
)

var _ feemarkettypes.DenomResolver = &DefaultFeemarketDenomResolver{}

type DefaultFeemarketDenomResolver struct{}

// ConvertToDenom returns "coin.Amount denom" for all coins that are not the denom.
func (r *DefaultFeemarketDenomResolver) ConvertToDenom(_ sdk.Context, coin sdk.DecCoin, denom string) (sdk.DecCoin, error) {
	if coin.Denom == denom {
		return coin, nil
	}

	return sdk.DecCoin{}, fmt.Errorf("error resolving denom")
}

func (r *DefaultFeemarketDenomResolver) ExtraDenoms(_ sdk.Context) ([]string, error) {
	return []string{}, nil
}
