package keeper

import (
	"context"

	"github.com/LimeChain/mantrachain/x/mns/types"
	"github.com/LimeChain/mantrachain/x/mns/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateDomainName(goCtx context.Context, msg *types.MsgCreateDomainName) (*types.MsgCreateDomainNameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	domainCtrl := NewDomainController(ctx, msg.Domain).WithStore(k)

	// Convert owner address string to sdk.AccAddress
	owner, _ := sdk.AccAddressFromBech32(msg.Creator)

	err := domainCtrl.
		MustExist().
		NotExpired().
		HasOwnerIfClosed(owner).
		Validate()

	if err != nil {
		return nil, err
	}

	domainNameCtrl := NewDomainNameController(ctx, msg.Domain, msg.DomainName).WithStore(k).WithConfiguration(k.GetParams(ctx))

	err = domainNameCtrl.
		MustNotExist().
		ValidDomainName().
		Validate()

	if err != nil {
		return nil, err
	}

	id := utils.GetDomainNameIndex(msg.Domain, msg.DomainName)
	didExecutor := NewDidExecutor(id, owner, msg.PubKeyHex, msg.PubKeyType)

	_, err = didExecutor.SetDid(ctx, k.didKeeper)
	if err != nil {
		return nil, err
	}

	newDomainName := types.DomainName{
		Index:      id,
		Creator:    owner,
		Domain:     msg.Domain,
		DomainName: msg.DomainName,
		Did:        didExecutor.GetDidId(),
		Owner:      owner,
	}

	k.SetDomainName(
		ctx,
		newDomainName,
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyDomain, msg.Domain),
			sdk.NewAttribute(types.AttributeKeyDomainName, msg.DomainName),
			sdk.NewAttribute(types.AttributeKeyDid, didExecutor.GetDidId()),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyCreator, owner.String()),
		),
	)
	return &types.MsgCreateDomainNameResponse{}, nil
}
