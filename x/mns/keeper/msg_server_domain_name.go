package keeper

import (
	"context"

	"github.com/LimeChain/mantrachain/x/mns/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateDomainName(goCtx context.Context, msg *types.MsgCreateDomainName) (*types.MsgCreateDomainNameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Convert owner address string to sdk.AccAddress
	owner, _ := sdk.AccAddressFromBech32(msg.Creator)
	// Check if the value already exists
	domain, isDomainFound := k.GetDomain(ctx, msg.Domain)
	// If a domain is found in store
	if !isDomainFound {
		// Throw an error
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Domain not set")
	}
	_, isDomainNameFound := k.GetDomainName(
		ctx,
		msg.Domain,
		msg.DomainName,
	)
	if isDomainNameFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Domain name already set")
	}

	var domainName = types.DomainName{
		Creator:    msg.Creator,
		Index:      msg.DomainName + "@" + msg.Domain,
		Domain:     domain.Index,
		DomainName: msg.DomainName,
		Owner:      owner.String(),
	}

	k.SetDomainName(
		ctx,
		domainName,
	)
	return &types.MsgCreateDomainNameResponse{}, nil
}
