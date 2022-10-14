package keeper

import (
	"github.com/LimeChain/mantrachain/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO: Refactor
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
	cw20ContractAddress sdk.AccAddress,
	delegate bool,
	stakingChain string,
	stakingValidator string,
) (bool, error) {
	var isStaked bool = false
	var wasmExecutor *WasmExecutor = nil

	if minPrice.IsNil() || minPrice.IsZero() {
		return isStaked, nil
	}

	currAmount := sdk.NewInt(0)

	// The rooyalties amount is calculated based on the price of the NFT
	for _, earning := range nftsEarningsOnSale {
		if !initialSale && types.MarketplaceEarningType(earning.Type) == types.Initially {
			continue
		}

		if !earning.Percentage.IsNil() && !earning.Percentage.IsZero() {
			earningAmount := earning.Percentage.Mul(minPrice.Amount).ToDec().Quo(sdk.NewDec(100))

			if earningAmount.GT(sdk.NewDec(1)) {
				earningCoin := sdk.NewCoin(minPrice.GetDenom(), earningAmount.TruncateInt())
				var err error

				if !cw20ContractAddress.Empty() {
					if wasmExecutor == nil {
						wasmExecutor = NewWasmExecutor(ctx, k.wasmViewKeeper, k.wasmContractKeeper)
					}
					err = wasmExecutor.Transfer(cw20ContractAddress, buyer, sdk.AccAddress(earning.Address), earningCoin.Amount.Abs().Uint64())
				} else {
					err = k.bk.SendCoins(ctx, buyer, sdk.AccAddress(earning.Address), []sdk.Coin{earningCoin})
				}

				if err != nil {
					return isStaked, err
				}

				currAmount = currAmount.Add(earningCoin.Amount)
			}
		}
	}

	// The amount supposed to be staked on a validator
	if !nftsVaultLockPercentage.IsNil() && !nftsVaultLockPercentage.IsZero() {
		lockAmount := nftsVaultLockPercentage.Mul(minPrice.Amount).ToDec().Quo(sdk.NewDec(100))
		var err error

		if lockAmount.GT(sdk.NewDec(1)) {
			lockCoin := sdk.NewCoin(minPrice.GetDenom(), lockAmount.TruncateInt())

			vaultExecutor := NewVaultExecutor(ctx, k.vaultKeeper)
			isStaked, err = vaultExecutor.UpsertNftStake(marketplaceIndex, collectionIndex, nftIndex, buyer, lockCoin, delegate, stakingChain, stakingValidator)

			if err != nil {
				return isStaked, err
			}

			currAmount = currAmount.Add(lockCoin.Amount)
		}
	}

	// Transfer the remainning amount to the nft owner.
	if minPrice.Amount.GT(currAmount) {
		remainning := minPrice.Amount.Sub(currAmount)
		var err error

		if !cw20ContractAddress.Empty() {
			if wasmExecutor == nil {
				wasmExecutor = NewWasmExecutor(ctx, k.wasmViewKeeper, k.wasmContractKeeper)
			}
			err = wasmExecutor.Transfer(cw20ContractAddress, buyer, nftOwner, remainning.Abs().Uint64())
		} else {
			err = k.bk.SendCoins(ctx, buyer, nftOwner, []sdk.Coin{sdk.NewCoin(minPrice.GetDenom(), remainning)})
		}

		if err != nil {
			return isStaked, err
		}
	}

	return isStaked, nil
}
