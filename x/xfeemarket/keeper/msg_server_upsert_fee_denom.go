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

	if err := k.DenomMultipliers.Set(ctx, msg.Denom, sdk.DecProto{Dec: msg.Multiplier}); err != nil {
		return nil, err
	}

	return &types.MsgUpsertFeeDenomResponse{}, nil
}
