package keeper

import (
	"context"

    "github.com/LimeChain/mantrachain/x/mns/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)


func (k msgServer) CreateDomain(goCtx context.Context,  msg *types.MsgCreateDomain) (*types.MsgCreateDomainResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

    // TODO: Handling the message
    _ = ctx

	return &types.MsgCreateDomainResponse{}, nil
}
