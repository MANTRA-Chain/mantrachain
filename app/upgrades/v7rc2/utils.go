package v7rc2

import (
	math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	UOM     = "uom"
	AMANTRA = "amantra"
	MANTRA  = "mantra"
)

var ScalingFactor = math.NewInt(4_000_000_000_000)

func convertCoinToNewDenom(uomCoin sdk.Coin) sdk.Coin {
	return sdk.NewCoin(AMANTRA, uomCoin.Amount.Mul(ScalingFactor))
}

func convertDecCoinToNewDenom(uomCoin sdk.DecCoin) sdk.DecCoin {
	return sdk.NewDecCoinFromDec(AMANTRA, uomCoin.Amount.MulInt(ScalingFactor))
}

func convertCoinsToNewDenom(coins sdk.Coins) sdk.Coins {
	newCoins := sdk.NewCoins()
	for _, coin := range coins {
		if coin.Denom == UOM {
			newCoins = newCoins.Add(convertCoinToNewDenom(coin))
		} else {
			newCoins = newCoins.Add(coin)
		}
	}
	return newCoins
}

func convertDecCoinsToNewDenom(coins sdk.DecCoins) sdk.DecCoins {
	newCoins := sdk.DecCoins{}
	for _, coin := range coins {
		if coin.Denom == UOM {
			newCoins = newCoins.Add(convertDecCoinToNewDenom(coin))
		} else {
			newCoins = newCoins.Add(coin)
		}
	}
	return newCoins
}
