package v7rc0

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func migrateAuthz(ctx sdk.Context, authzKeeper authzkeeper.Keeper) {
	authzKeeper.IterateGrants(ctx, func(granterAddr, granteeAddr sdk.AccAddress, grant authz.Grant) bool {
		authorization, err := grant.GetAuthorization()
		if err != nil {
			ctx.Logger().Error("failed to get authorization from grant", "error", err)
			return false // continue
		}

		var needsUpdate bool
		var updatedAuth authz.Authorization

		switch auth := authorization.(type) {
		case *banktypes.SendAuthorization:
			newSpendLimit := convertCoinsToNewDenom(auth.SpendLimit)
			auth.SpendLimit = newSpendLimit
			updatedAuth = auth
			needsUpdate = true
		case *stakingtypes.StakeAuthorization:
			if auth.MaxTokens != nil && auth.MaxTokens.Denom == UOM {
				newMaxTokens := convertCoinToNewDenom(*auth.MaxTokens)
				auth.MaxTokens = &newMaxTokens
				updatedAuth = auth
				needsUpdate = true
			}
		}

		if needsUpdate {
			err = authzKeeper.SaveGrant(ctx, granteeAddr, granterAddr, updatedAuth, grant.Expiration)
			if err != nil {
				ctx.Logger().Error("failed to save grant", "error", err)
				return false // continue
			}
		}

		return false
	})
}
