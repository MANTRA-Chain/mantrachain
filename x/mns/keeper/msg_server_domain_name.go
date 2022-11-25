package keeper

import (
	"context"
	"strings"

	"github.com/LimeChain/mantrachain/x/mns/types"
	"github.com/LimeChain/mantrachain/x/mns/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateDomainName(goCtx context.Context, msg *types.MsgCreateDomainName) (*types.MsgCreateDomainNameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if strings.TrimSpace(msg.Domain) == "" {
		return nil, sdkerrors.Wrap(types.ErrInvalidDomain, "domain should not be empty")
	}

	if strings.TrimSpace(msg.DomainName) == "" {
		return nil, sdkerrors.Wrap(types.ErrInvalidDomainName, "domain name should not be empty")
	}

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

	index := types.GetDomainNameIndex(msg.Domain, msg.DomainName)
	indexHex := utils.GetIndexHex(index)

	didExecutor := NewDidExecutor(ctx, owner, msg.PubKeyHex, msg.PubKeyType, k.didKeeper)
	_, err = didExecutor.SetDid(indexHex)
	if err != nil {
		return nil, err
	}

	newDomainName := types.DomainName{
		Index:      string(index),
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
	return &types.MsgCreateDomainNameResponse{
		Creator:    msg.Creator,
		Domain:     msg.Domain,
		DomainName: msg.DomainName,
	}, nil
}
