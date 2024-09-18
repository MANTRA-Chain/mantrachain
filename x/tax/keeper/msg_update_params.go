package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/MANTRA-Chain/mantrachain/x/tax/types"
)

func (k msgServer) UpdateParams(ctx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if _, err := k.addressCodec.StringToBytes(req.Authority); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	if k.GetAuthority() != req.Authority {
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.GetAuthority(), req.Authority)
	}

	updateParams, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	if req.Proportion != "" {
		updateParams.Proportion, err = math.LegacyNewDecFromStr(req.Proportion)
		if err != nil {
			return nil, err
		}
	}

	if req.McaAddress != "" {
		updateParams.McaAddress = req.McaAddress
	}

	if err := updateParams.Validate(); err != nil {
		return nil, err
	}

	if err := k.Params.Set(ctx, updateParams); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
