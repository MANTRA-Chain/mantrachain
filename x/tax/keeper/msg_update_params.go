package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/MANTRA-Chain/mantrachain/x/tax/types"
)

func (k msgServer) UpdateParams(ctx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if _, err := k.addressCodec.StringToBytes(req.Authority); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	if req.Authority != params.McaAddress {
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "invalid sender; expected mcaAddress %s, got %s", params.McaAddress, req.Authority)
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
		if updateParams.McaTax.GT(types.MaxMcaTax) {
			return nil, fmt.Errorf("mca tax %s cannot exceed maximum of %s", updateParams.McaTax, types.MaxMcaTax)
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
