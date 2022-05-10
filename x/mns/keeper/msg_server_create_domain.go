package keeper

import (
	"context"

	"github.com/LimeChain/mantrachain/x/mns/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateDomain(goCtx context.Context, msg *types.MsgCreateDomain) (*types.MsgCreateDomainResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// Try getting a domain from the store
	_, isFound := k.GetDomain(ctx, msg.Domain)
	// Convert owner address string to sdk.AccAddress
	owner, _ := sdk.AccAddressFromBech32(msg.Creator)
	// If a domain is found in store
	if isFound {
		// Throw an error
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Domain already exists")
	}
	// Create a domain record
	newDomain := types.Domain{
		Index:      msg.Domain,
		Domain:     msg.Domain,
		DomainType: msg.DomainType,
		Owner:      owner.String(),
	}
	// Write domain information to the store
	k.SetDomain(ctx, newDomain)

	return &types.MsgCreateDomainResponse{}, nil
}
