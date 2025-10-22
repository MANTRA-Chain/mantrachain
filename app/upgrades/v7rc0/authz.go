package v7rc0

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func migrateAuthz(ctx sdk.Context, authzKeeper authzkeeper.Keeper) {
	authzKeeper.IterateGrants(ctx, func(granterAddr, granteeAddr sdk.AccAddress, grant authz.Grant) bool {
		authorization, err := grant.GetAuthorization()
		if err != nil {
			ctx.Logger().Error("failed to get authorization from grant", "error", err)
			return false // continue
		}

		if sendAuth, ok := authorization.(*banktypes.SendAuthorization); ok {
			newSpendLimit := convertCoinsToNewDenom(sendAuth.SpendLimit)
			sendAuth.SpendLimit = newSpendLimit
			err = authzKeeper.SaveGrant(ctx, granteeAddr, granterAddr, sendAuth, grant.Expiration)
			if err != nil {
				ctx.Logger().Error("failed to save grant", "error", err)
				return false // continue
			}
		}

		return false
	})
}
