package keeper

import (
	"context"
	"strings"

	"github.com/LimeChain/mantrachain/x/mns/types"
	"github.com/LimeChain/mantrachain/x/mns/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateDomain(goCtx context.Context, msg *types.MsgCreateDomain) (*types.MsgCreateDomainResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if strings.TrimSpace(msg.Domain) == "" {
		return nil, sdkerrors.Wrap(types.ErrInvalidDomain, "domain should not be empty")
	}

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

	index := types.GetDomainIndex(msg.Domain)
	indexHex := utils.GetIndexHex(index)

	didExecutor := NewDidExecutor(ctx, owner, msg.PubKeyHex, msg.PubKeyType, k.didKeeper)
	_, err = didExecutor.SetDid(indexHex)
	if err != nil {
		return nil, err
	}

	// Create a domain record
	newDomain := types.Domain{
		Index:      string(index),
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
	return &types.MsgCreateDomainResponse{
		Creator:    owner.String(),
		Domain:     msg.Domain,
		DomainType: msg.DomainType,
	}, nil
}
