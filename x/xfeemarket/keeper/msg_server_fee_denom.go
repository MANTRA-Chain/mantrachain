package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/MANTRA-Chain/mantrachain/x/xfeemarket/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) UpsertFeeDenom(ctx context.Context, msg *types.MsgUpsertFeeDenom) (*types.MsgUpsertFeeDenomResponse, error) {
	if err := msg.Validate(); err != nil {
		return nil, err
	}

	authority := k.GetAuthority()
	if authority != msg.Authority {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid authority; expected %s, got %s", authority, msg.Authority)
	}

	if !k.Keeper.bankkeeper.HasDenomMetaData(ctx, msg.Denom) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "denom %s does not exist", msg.Denom)
	}

	if err := k.DenomMultipliers.Set(ctx, msg.Denom, sdk.DecProto{Dec: msg.Multiplier}); err != nil {
		return nil, err
	}

	return &types.MsgUpsertFeeDenomResponse{}, nil
}

func (k msgServer) RemoveFeeDenom(ctx context.Context, msg *types.MsgRemoveFeeDenom) (*types.MsgRemoveFeeDenomResponse, error) {
	if err := msg.Validate(); err != nil {
		return nil, err
	}

	authority := k.GetAuthority()
	if authority != msg.Authority {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid authority; expected %s, got %s", authority, msg.Authority)
	}

	hasDenom, err := k.DenomMultipliers.Has(ctx, msg.Denom)
	if err != nil {
		return nil, err
	}
	if !hasDenom {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "denom %s does not exist", msg.Denom)
	}

	if err := k.DenomMultipliers.Remove(ctx, msg.Denom); err != nil {
		return nil, err
	}

	return &types.MsgRemoveFeeDenomResponse{}, nil
}
