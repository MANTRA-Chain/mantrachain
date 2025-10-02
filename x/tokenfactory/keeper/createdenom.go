package keeper

import (
	"crypto/sha256"
	"fmt"

	"cosmossdk.io/errors"
	"github.com/MANTRA-Chain/mantrachain/v5/x/tokenfactory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	erc20types "github.com/cosmos/evm/x/erc20/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

// ConvertToBaseToken converts a fee amount in a whitelisted fee token to the base fee token amount
func (k Keeper) CreateDenom(ctx sdk.Context, creatorAddr, subdenom string) (newTokenDenom string, err error) {
	denom, err := k.validateCreateDenom(ctx, creatorAddr, subdenom)
	if err != nil {
		return "", err
	}

	err = k.chargeForCreateDenom(ctx, creatorAddr)
	if err != nil {
		return "", err
	}

	err = k.createDenomAfterValidation(ctx, creatorAddr, denom)
	return denom, err
}

// Runs CreateDenom logic after the charge and all denom validation has been handled.
// Made into a second function for genesis initialization.
func (k Keeper) createDenomAfterValidation(ctx sdk.Context, creatorAddr, denom string) (err error) {
	_, exists := k.bankKeeper.GetDenomMetaData(ctx, denom)
	if !exists {
		denomMetaData := banktypes.Metadata{
			DenomUnits: []*banktypes.DenomUnit{{
				Denom:    denom,
				Exponent: 0,
			}},
			Base:    denom,
			Name:    denom,
			Symbol:  denom,
			Display: denom,
		}

		k.bankKeeper.SetDenomMetaData(ctx, denomMetaData)
	}

	authorityMetadata := types.DenomAuthorityMetadata{
		Admin: creatorAddr,
	}
	err = k.setAuthorityMetadata(ctx, denom, authorityMetadata)
	if err != nil {
		return err
	}

	// create erc20 contractAddr and set token pair
	denomHash := sha256.Sum256([]byte(denom))
	ethContractAddr := ethcommon.BytesToAddress(denomHash[:])
	if k.erc20Keeper.IsERC20Registered(ctx, ethContractAddr) {
		return types.ErrInvalidDenom.Wrapf(
			"denom results in already registered ethContractAddr: %v, use a different subdenom and try again", ethContractAddr.Hex())
	}
	pair := erc20types.NewTokenPair(ethContractAddr, denom, erc20types.OWNER_MODULE)
	err = k.erc20Keeper.SetToken(ctx, pair)
	if err != nil {
		return err
	}

	err = k.erc20Keeper.EnableDynamicPrecompile(ctx, pair.GetERC20Contract())
	if err != nil {
		return err
	}

	k.addDenomFromCreator(ctx, creatorAddr, denom)
	return nil
}

// UpdateDenomWithERC20 registers erc20 precompile address for existing tokenfactory token
// to be used for migration in upgradehandler v5.0.0-rc3
func (k Keeper) UpdateDenomWithERC20(ctx sdk.Context, denom string) (err error) {
	denomMetaData, exists := k.bankKeeper.GetDenomMetaData(ctx, denom)
	if !exists {
		return types.ErrInvalidDenom.Wrapf("denom does not exist: %v", denom)
	}
	if denomMetaData.Symbol == "" {
		denomMetaData.Symbol = denom
	}
	if denomMetaData.Display == "" {
		denomMetaData.Display = denom
	}
	if len(denomMetaData.DenomUnits) == 0 {
		denomMetaData.DenomUnits = []*banktypes.DenomUnit{{
			Denom:    denom,
			Exponent: 0,
		}}
	}

	k.bankKeeper.SetDenomMetaData(ctx, denomMetaData)

	// create erc20 contractAddr and set token pair
	denomHash := sha256.Sum256([]byte(denom))
	ethContractAddr := ethcommon.BytesToAddress(denomHash[:])
	if k.erc20Keeper.IsERC20Registered(ctx, ethContractAddr) {
		// skip registering if hash address already registered
		k.Logger(ctx).Error("denom results in already registered ethContractAddr: %s", ethContractAddr.Hex())
		return nil
	}
	pair := erc20types.NewTokenPair(ethContractAddr, denom, erc20types.OWNER_MODULE)
	err = k.erc20Keeper.SetToken(ctx, pair)
	if err != nil {
		return err
	}
	err = k.erc20Keeper.EnableDynamicPrecompile(ctx, pair.GetERC20Contract())
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) validateCreateDenom(ctx sdk.Context, creatorAddr, subdenom string) (newTokenDenom string, err error) {
	// Temporary check until IBC bug is sorted out
	if k.bankKeeper.HasSupply(ctx, subdenom) {
		return "", fmt.Errorf("temporary error until IBC bug is sorted out, " +
			"can't create subdenoms that are the same as a native denom")
	}

	denom, err := types.GetTokenDenom(creatorAddr, subdenom)
	if err != nil {
		return "", err
	}

	_, found := k.bankKeeper.GetDenomMetaData(ctx, denom)
	if found {
		return "", types.ErrDenomExists
	}

	return denom, nil
}

func (k Keeper) chargeForCreateDenom(ctx sdk.Context, creatorAddr string) (err error) {
	params := k.GetParams(ctx)

	// ORIGINAL: if DenomCreationFee is non-zero, transfer the tokens from the creator
	// account to community pool
	// MODIFIED: if DenomCreationFee is non-zero, transfer the tokens from the creator
	// account to feeCollectorAddr
	if len(params.DenomCreationFee) != 0 {
		accAddr, err := sdk.AccAddressFromBech32(creatorAddr)
		if err != nil {
			return err
		}
		// Instead of funding community pool we send funds to fee collector addr
		// if err := k.communityPoolKeeper.FundCommunityPool(ctx, params.DenomCreationFee, accAddr); err != nil {
		//	return err
		//}

		feeCollectorAddr, err := sdk.AccAddressFromBech32(params.FeeCollectorAddress)
		if err != nil {
			return errors.Wrapf(err, "wrong fee collector address: %v", err)
		}

		err = k.bankKeeper.SendCoins(
			ctx,
			accAddr, feeCollectorAddr,
			params.DenomCreationFee,
		)
		if err != nil {
			return errors.Wrap(err, "unable to send coins to fee collector")
		}
	}

	// if DenomCreationGasConsume is non-zero, consume the gas
	if params.DenomCreationGasConsume != 0 {
		ctx.GasMeter().ConsumeGas(params.DenomCreationGasConsume, "consume denom creation gas")
	}

	return nil
}
