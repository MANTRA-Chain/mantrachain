package keeper

import (
	"github.com/LimeChain/mantrachain/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) CollectFees(
	ctx sdk.Context,
	minPrice *sdk.Coin,
	nftsEarningsOnSale []*types.MarketplaceEarning,
	initiallyNftsVaultLockPercentage sdk.Int,
	buyer sdk.AccAddress,
	nftOwner sdk.AccAddress,
	initialSale bool,
	cw20ContractAddress sdk.AccAddress,
) (sdk.Coin, error) {
	var lockCoin sdk.Coin
	var wasmExecutor *WasmExecutor = nil

	if minPrice.IsNil() || minPrice.IsZero() {
		return lockCoin, nil
	}

	currAmount := sdk.NewInt(0)

	// The royalties amount is calculated based on the price of the NFT
	for _, earning := range nftsEarningsOnSale {
		if (!initialSale && types.MarketplaceEarningType(earning.Type) == types.Initially) ||
			(initialSale && types.MarketplaceEarningType(earning.Type) == types.Repetitive) {
			continue
		}

		if !earning.Percentage.IsNil() && !earning.Percentage.IsZero() {
			earningAmount := sdk.NewDecFromInt(earning.Percentage.Mul(minPrice.Amount)).Quo(sdk.NewDec(100))

			if earningAmount.GT(sdk.NewDec(1)) {
				earningCoin := sdk.NewCoin(minPrice.GetDenom(), earningAmount.TruncateInt())
				var err error

				if cw20ContractAddress.Empty() {
					err = k.bk.SendCoins(ctx, buyer, sdk.AccAddress(earning.Address), []sdk.Coin{earningCoin})
				} else {
					if wasmExecutor == nil {
						wasmExecutor = NewWasmExecutor(ctx, k.wasmViewKeeper, k.wasmContractKeeper)
					}
					err = wasmExecutor.Transfer(cw20ContractAddress, buyer, sdk.AccAddress(earning.Address), earningCoin.Amount.String())
				}

				if err != nil {
					return lockCoin, err
				}

				currAmount = currAmount.Add(earningCoin.Amount)
			}
		}
	}

	// The amount supposed to be staked on a validator
	if initialSale &&
		!initiallyNftsVaultLockPercentage.IsNil() &&
		!initiallyNftsVaultLockPercentage.IsZero() {
		lockAmount := sdk.NewDecFromInt(initiallyNftsVaultLockPercentage.Mul(minPrice.Amount)).Quo(sdk.NewDec(100))
		var err error

		if lockAmount.GT(sdk.NewDec(1)) {
			lockCoin = sdk.NewCoin(minPrice.GetDenom(), lockAmount.TruncateInt())

			if err != nil {
				return lockCoin, err
			}

			if !cw20ContractAddress.Empty() {
				if wasmExecutor == nil {
					wasmExecutor = NewWasmExecutor(ctx, k.wasmViewKeeper, k.wasmContractKeeper)
				}
				// Burn the cw20 staking amount which goes for anothe chain delegation
				err = wasmExecutor.Burn(cw20ContractAddress, buyer, lockCoin.Amount.String())

				if err != nil {
					return lockCoin, err
				}
			}

			currAmount = currAmount.Add(lockCoin.Amount)
		}
	}

	// Transfer the remainning amount to the nft owner.
	if minPrice.Amount.GT(currAmount) {
		remainning := minPrice.Amount.Sub(currAmount)
		var err error

		// The remaining amount is transferred to the nft owner
		if cw20ContractAddress.Empty() {
			err = k.bk.SendCoins(ctx, buyer, nftOwner, []sdk.Coin{sdk.NewCoin(minPrice.GetDenom(), remainning)})
		} else {
			if wasmExecutor == nil {
				wasmExecutor = NewWasmExecutor(ctx, k.wasmViewKeeper, k.wasmContractKeeper)
			}
			err = wasmExecutor.Transfer(cw20ContractAddress, buyer, nftOwner, remainning.String())
		}

		if err != nil {
			return lockCoin, err
		}
	}

	return lockCoin, nil
}
