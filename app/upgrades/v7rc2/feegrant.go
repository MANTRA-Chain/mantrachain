package v7rc2

import (
	errorsmod "cosmossdk.io/errors"
	feegrant "cosmossdk.io/x/feegrant"
	feegrantkeeper "cosmossdk.io/x/feegrant/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func migrateFeeGrant(ctx sdk.Context, feeGrantKeeper feegrantkeeper.Keeper) {
	if err := feeGrantKeeper.IterateAllFeeAllowances(ctx, func(grant feegrant.Grant) bool {
		allowance, err := grant.GetGrant()
		if err != nil {
			ctx.Logger().Error("failed to get grant", "grant", grant, "error", err)
			return false // continue
		}

		var newAllowance feegrant.FeeAllowanceI
		newAllowance, err = migrateFeeAllowance(allowance)
		if err != nil {
			ctx.Logger().Error("failed to migrate fee allowance", "grant", grant, "error", err)
			return false // continue
		}

		granterAddr := sdk.MustAccAddressFromBech32(grant.Granter)
		granteeAddr := sdk.MustAccAddressFromBech32(grant.Grantee)

		if err := feeGrantKeeper.UpdateAllowance(ctx, granterAddr, granteeAddr, newAllowance); err != nil {
			ctx.Logger().Error("failed to update allowance", "granter", grant.Granter, "grantee", grant.Grantee, "error", err)
		}

		return false
	}); err != nil {
		ctx.Logger().Error("failed to iterate fee allowances", "error", err)
	}
}

func migrateFeeAllowance(allowance feegrant.FeeAllowanceI) (feegrant.FeeAllowanceI, error) {
	switch allowance := allowance.(type) {
	case *feegrant.BasicAllowance:
		return migrateBasicAllowance(allowance), nil
	case *feegrant.PeriodicAllowance:
		return migratePeriodicAllowance(allowance), nil
	case *feegrant.AllowedMsgAllowance:
		nested, err := allowance.GetAllowance()
		if err != nil {
			return nil, err
		}
		migratedNested, err := migrateFeeAllowance(nested)
		if err != nil {
			return nil, err
		}
		return feegrant.NewAllowedMsgAllowance(migratedNested, allowance.AllowedMessages)
	default:
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidType, "unsupported fee allowance type: %T", allowance)
	}
}

func migrateBasicAllowance(allowance *feegrant.BasicAllowance) *feegrant.BasicAllowance {
	allowance.SpendLimit = convertCoinsToNewDenom(allowance.SpendLimit)
	return allowance
}

func migratePeriodicAllowance(allowance *feegrant.PeriodicAllowance) *feegrant.PeriodicAllowance {
	allowance.Basic.SpendLimit = convertCoinsToNewDenom(allowance.Basic.SpendLimit)
	allowance.PeriodSpendLimit = convertCoinsToNewDenom(allowance.PeriodSpendLimit)
	allowance.PeriodCanSpend = convertCoinsToNewDenom(allowance.PeriodCanSpend)
	return allowance
}
