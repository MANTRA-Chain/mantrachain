package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"
	"github.com/MANTRA-Chain/mantrachain/v8/x/sanction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

var _ banktypes.SendRestrictionFn = Keeper{}.SendRestrictionFn

// SendRestrictionFn is a bank send restriction that stops coins leaving a
// blacklisted address. The ante handler only stops a blacklisted account
// signing txs; this also covers sends that never pass through it, such as
// ICA host txs, authz executed by a contract, and ERC20 precompile
// transferFrom on an earlier allowance.
//
// Module accounts can't be sanctioned, so this never blocks one.
//
// The lookup isn't charged gas, so transfer gas costs stay as they were.
func (k Keeper) SendRestrictionFn(ctx context.Context, fromAddr, toAddr sdk.AccAddress, _ sdk.Coins) (sdk.AccAddress, error) {
	from := fromAddr.String()
	sdkCtx := sdk.UnwrapSDKContext(ctx).WithGasMeter(storetypes.NewInfiniteGasMeter())
	has, err := k.BlacklistAccounts.Has(sdkCtx, from)
	if err != nil {
		return nil, err
	}
	if has {
		return nil, errorsmod.Wrapf(types.ErrAccountBlacklisted, "%s is blacklisted and cannot send funds", from)
	}
	return toAddr, nil
}
