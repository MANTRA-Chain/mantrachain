package keeper

import (
	"context"

    "github.com/MANTRA-Chain/mantrachain/x/xfeemarket/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)


func (k msgServer) UpsertFeeDenom(ctx context.Context,  msg *types.MsgUpsertFeeDenom) (*types.MsgUpsertFeeDenomResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

    // TODO: Handle the message

	return &types.MsgUpsertFeeDenomResponse{}, nil
}
