package keeper

import (
	"context"

	"github.com/LimeChain/mantrachain/x/mns/types"
	"github.com/LimeChain/mantrachain/x/mns/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateDomain(goCtx context.Context, msg *types.MsgCreateDomain) (*types.MsgCreateDomainResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ctrl := NewDomainController(ctx, msg.Domain).WithStore(k).WithConfiguration(k.GetParams(ctx))

	err := ctrl.
		MustNotExist().
		ValidDomain().
		Validate()

	if err != nil {
		return nil, err
	}

	// Convert owner address string to sdk.AccAddress
	owner, _ := sdk.AccAddressFromBech32(msg.Creator)

	id := utils.GetDomainIndex(msg.Domain)
	didExecutor := NewDidExecutor(id, owner, msg.PubKeyHex, msg.PubKeyType)

	_, err = didExecutor.SetDid(ctx, k.didKeeper)
	if err != nil {
		return nil, err
	}

	// Create a domain record
	newDomain := types.Domain{
		Index:      id,
		Creator:    owner,
		Domain:     msg.Domain,
		DomainType: types.DomainType(msg.DomainType),
		Did:        didExecutor.GetDidId(),
		Owner:      owner,
	}
	// Write domain information to the store
	k.SetDomain(ctx, newDomain)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyDomain, msg.Domain),
			sdk.NewAttribute(types.AttributeKeyDid, didExecutor.GetDidId()),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyCreator, owner.String()),
		),
	)
	return &types.MsgCreateDomainResponse{}, nil
}
