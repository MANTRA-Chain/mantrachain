package utils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func SumCoins(coins []*sdk.DecCoin, minTreshold sdk.Coin) (result []*sdk.DecCoin) {
	sum := make(map[string]sdk.Dec)

	for _, coin := range coins {
		sum[coin.Denom] = sum[coin.Denom].Add(coin.Amount)
	}

	for denom, amount := range sum {
		if denom != minTreshold.Denom || amount.GT(sdk.NewDec(minTreshold.Amount.Int64())) {
			result = append(result, &sdk.DecCoin{
				Denom:  denom,
				Amount: amount,
			})
		}
	}

	return
}
