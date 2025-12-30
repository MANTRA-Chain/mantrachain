package v7

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	tokenfactorykeeper "github.com/MANTRA-Chain/mantrachain/v7/x/tokenfactory/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	precisebankkeeper "github.com/cosmos/evm/x/precisebank/keeper"
	precisebanktypes "github.com/cosmos/evm/x/precisebank/types"
)

func migrateBank(ctx sdk.Context, bankKeeper bankkeeper.Keeper, tokenFactoryKeeper tokenfactorykeeper.Keeper, accountKeeper authkeeper.AccountKeeper) error {
	// add denom metadata for amantra
	bankKeeper.SetDenomMetaData(ctx, banktypes.Metadata{
		Description: "The native staking token of the Mantrachain.",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    AMANTRA,
				Exponent: 0,
			},
			{
				Denom:    "mantra",
				Exponent: 18,
			},
		},
		Base:    AMANTRA,
		Display: MANTRA,
		Name:    MANTRA,
		Symbol:  MANTRA,
	})

	// update uom metadata description
	uomMetaData, found := bankKeeper.GetDenomMetaData(ctx, UOM)
	if !found {
		return errorsmod.Wrapf(banktypes.ErrDenomMetadataNotFound, "uom metadata not found")
	}
	uomMetaData.Description = "The legacy uom token is deprecated and has been replaced by amantra."
	bankKeeper.SetDenomMetaData(ctx, uomMetaData)

	// migrate all uom balances to amantra at the new scaling factor
	beforeSupply := bankKeeper.GetSupply(ctx, UOM)
	escrowSupply := sdk.NewCoin(UOM, math.ZeroInt())
	var err error
	bankKeeper.IterateAllBalances(ctx, func(addr sdk.AccAddress, coin sdk.Coin) (stop bool) {
		if coin.Denom != UOM {
			return false
		}

		// skip IBC escrow accounts
		if tokenFactoryKeeper.IsEscrowAddress(ctx, addr) {
			escrowSupply = escrowSupply.Add(coin)
			return false
		}

		amantraCoin := convertCoinToNewDenom(coin)
		err = bankKeeper.SendCoinsFromAccountToModule(ctx, addr, UpgradeName, sdk.NewCoins(coin))
		if err != nil {
			return true
		}
		err = bankKeeper.MintCoins(ctx, UpgradeName, sdk.NewCoins(amantraCoin))
		if err != nil {
			return true
		}
		err = bankKeeper.SendCoinsFromModuleToAccount(ctx, UpgradeName, addr, sdk.NewCoins(amantraCoin))
		return err != nil
	})
	if err != nil {
		return err
	}
	// check that the amount in this upgrade account matches the expected supply change
	toBurnAmount := bankKeeper.GetBalance(ctx, accountKeeper.GetModuleAddress(UpgradeName), UOM)
	expectedEscrowSupply := beforeSupply.Sub(toBurnAmount)
	if !escrowSupply.IsEqual(expectedEscrowSupply) {
		return errorsmod.Wrapf(
			banktypes.ErrInputOutputMismatch,
			"escrow Supply mismatch after migration: got %s, actual %s", expectedEscrowSupply.String(), escrowSupply.String())
	}
	// burn the old uom tokens from this upgrade account
	err = bankKeeper.BurnCoins(ctx, UpgradeName, sdk.NewCoins(toBurnAmount))
	if err != nil {
		return err
	}

	return nil
}

func migratePreciseBank(ctx sdk.Context, preciseBankKeeper precisebankkeeper.Keeper, bankKeeper bankkeeper.Keeper, accountKeeper authkeeper.AccountKeeper) (err error) {
	preciseBankBalance := bankKeeper.GetAllBalances(ctx, accountKeeper.GetModuleAddress(precisebanktypes.ModuleName))
	for _, coin := range preciseBankBalance {
		if coin.Denom != UOM {
			return errorsmod.Wrapf(banktypes.ErrInputOutputMismatch, "precise bank module has non-uom balance: %s", coin.String())
		}
	}
	if err = bankKeeper.SendCoinsFromModuleToModule(ctx, precisebanktypes.ModuleName, UpgradeName, preciseBankBalance); err != nil {
		return err
	}
	if err = bankKeeper.BurnCoins(ctx, UpgradeName, preciseBankBalance); err != nil {
		return err
	}
	preciseBankKeeper.IterateFractionalBalances(ctx, func(addr sdk.AccAddress, fractionalAmount math.Int) (stop bool) {
		// get current integer balance
		integerCoin := sdk.NewCoin(AMANTRA, fractionalAmount.Mul(math.NewInt(4)))
		preciseBankKeeper.SetFractionalBalance(ctx, addr, math.ZeroInt())
		err = bankKeeper.MintCoins(ctx, UpgradeName, sdk.NewCoins(integerCoin))
		if err != nil {
			return true
		}
		err = bankKeeper.SendCoinsFromModuleToAccount(ctx, UpgradeName, addr, sdk.NewCoins(integerCoin))
		return err != nil
	})
	if err != nil {
		return err
	}
	return nil
}
