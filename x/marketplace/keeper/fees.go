package keeper

import (
	math "math"

	"github.com/LimeChain/mantrachain/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) CollectFeesAndDelegateStake(
	ctx sdk.Context,
	minPrice *sdk.Coin,
	nftsEarningsOnSale []*types.MarketplaceEarning,
	nftsVaultLockPercentage sdk.Int,
	buyer sdk.AccAddress,
	collectionOwner sdk.AccAddress,
	nftOwner sdk.AccAddress,
	marketplaceIndex []byte,
	collectionIndex []byte,
	nftIndex []byte,
	initialSale bool,
) (bool, error) {
	staked := false

	if minPrice.IsNil() || minPrice.IsZero() {
		return staked, nil
	}

	currAmount := sdk.NewInt(0)

	for _, earning := range nftsEarningsOnSale {
		if !initialSale && types.MarketplaceEarningType(earning.Type) == types.Initially {
			continue
		}

		if !earning.Percentage.IsNil() && !earning.Percentage.IsZero() {
			earningAmount := earning.Percentage.Mul(minPrice.Amount).ToDec().MustFloat64() / 100

			if earningAmount > 1 {
				earningCoin := sdk.NewCoin(minPrice.GetDenom(), sdk.NewDec(int64(math.Floor(earningAmount))).TruncateInt())
				err := k.bk.SendCoins(ctx, buyer, sdk.AccAddress(earning.Address), []sdk.Coin{earningCoin})

				if err != nil {
					return staked, err
				}

				currAmount = currAmount.Add(earningCoin.Amount)
			}
		}
	}

	if !nftsVaultLockPercentage.IsNil() && !nftsVaultLockPercentage.IsZero() {
		lockAmount := nftsVaultLockPercentage.Mul(minPrice.Amount).ToDec().MustFloat64() / 100

		if lockAmount > 1 {
			lockCoin := sdk.NewCoin(minPrice.GetDenom(), sdk.NewDec(int64(math.Floor(lockAmount))).TruncateInt())

			vaultExecutor := NewVaultExecutor(ctx, k.vaultKeeper)
			err := vaultExecutor.UpsertNftStake(marketplaceIndex, collectionIndex, nftIndex, buyer, lockCoin, true)

			if err != nil {
				return staked, err
			}

			currAmount = currAmount.Add(lockCoin.Amount)

			staked = true
		}
	}

	// Transfer the unpaid amount to the nft owner.
	if minPrice.Amount.GT(currAmount) {
		err := k.bk.SendCoins(ctx, buyer, nftOwner, []sdk.Coin{sdk.NewCoin(minPrice.GetDenom(), minPrice.Amount.Sub(currAmount))})

		if err != nil {
			return staked, err
		}
	}

	return staked, nil
}
