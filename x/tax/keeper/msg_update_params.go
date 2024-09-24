package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/MANTRA-Chain/mantrachain/x/tax/types"
)

func (k msgServer) UpdateParams(ctx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if _, err := k.addressCodec.StringToBytes(req.Admin); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	if req.Admin != params.Admin {
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "invalid admin; expected %s, got %s", params.Admin, req.Admin)
	}

	updateParams, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	if req.McaTax != "" {
		updateParams.McaTax, err = math.LegacyNewDecFromStr(req.McaTax)
		if err != nil {
			return nil, err
		}
		// Check against MaxMcaTax
		if updateParams.McaTax.GT(updateParams.MaxMcaTax) {
			return nil, fmt.Errorf("mca tax cannot exceed maximum of %s", updateParams.MaxMcaTax)
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
