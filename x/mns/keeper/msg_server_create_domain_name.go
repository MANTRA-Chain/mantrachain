package keeper

import (
	"context"

    "github.com/LimeChain/mantrachain/x/mns/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)


func (k msgServer) CreateDomainName(goCtx context.Context,  msg *types.MsgCreateDomainName) (*types.MsgCreateDomainNameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

    // TODO: Handling the message
    _ = ctx

	return &types.MsgCreateDomainNameResponse{}, nil
}
