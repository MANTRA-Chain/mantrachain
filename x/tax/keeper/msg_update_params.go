package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/MANTRA-Chain/mantrachain/x/tax/types"
)

func (k msgServer) UpdateParams(ctx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if _, err := k.addressCodec.StringToBytes(req.Authority); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	if k.GetAuthority() != req.Authority {
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.GetAuthority(), req.Authority)
	}

	newParams := types.NewParams(req.Proportion, req.McaAddress)

	if err := newParams.Validate(); err != nil {
		return nil, err
	}

	if err := k.Params.Set(ctx, newParams); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
