package utils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func SumCoins(coins []*sdk.DecCoin) (result []*sdk.DecCoin) {
	sum := make(map[string]sdk.Dec)

	for _, coin := range coins {
		if sum[coin.Denom].IsNil() {
			sum[coin.Denom] = sdk.NewDec(0)
		}
		sum[coin.Denom] = sum[coin.Denom].Add(coin.Amount)
	}

	for denom, amount := range sum {
		result = append(result, &sdk.DecCoin{
			Denom:  denom,
			Amount: amount,
		})

	}

	return
}
