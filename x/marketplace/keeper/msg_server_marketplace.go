package keeper

import (
	"context"
	"strings"

	"github.com/LimeChain/mantrachain/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) RegisterMarketplace(goCtx context.Context, msg *types.MsgRegisterMarketplace) (*types.MsgRegisterMarketplaceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator, err := sdk.AccAddressFromBech32(msg.Creator)

	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(msg.Marketplace.Id) == "" {
		return nil, sdkerrors.Wrap(types.ErrInvalidMarketplaceId, "marketplace id should not be empty")
	}

	marketplaceController := NewMarketplaceController(ctx, creator).
		WithMetadata(msg.Marketplace).
		WithStore(k).
		WithConfiguration(k.GetParams(ctx))

	err = marketplaceController.
		MustNotExist().
		ValidMetadata().
		Validate()

	if err != nil {
		return nil, err
	}

	marketplaceIndex := marketplaceController.getIndex()
	marketplaceId := marketplaceController.getId()

	newMarketplace := types.Marketplace{
		Index:       marketplaceIndex,
		Id:          marketplaceId,
		Name:        msg.Marketplace.Name,
		Description: msg.Marketplace.Description,
		Url:         msg.Marketplace.Url,
		Opened:      msg.Marketplace.Opened,
		Options:     msg.Marketplace.Options,
		Attributes:  msg.Marketplace.Attributes,
		Images:      msg.Marketplace.Images,
		Links:       msg.Marketplace.Links,
		Data:        msg.Marketplace.Data,
		Creator:     creator,
		Owner:       creator,
	}

	k.SetMarketplace(ctx, newMarketplace)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgRegisterMarketplace),
			sdk.NewAttribute(types.AttributeKeyMarketplaceId, marketplaceId),
			sdk.NewAttribute(types.AttributeKeySigner, creator.String()),
			sdk.NewAttribute(types.AttributeKeyOwner, creator.String()),
		),
	)

	return &types.MsgRegisterMarketplaceResponse{
		MarketplaceId:      marketplaceId,
		MarketplaceCreator: creator.String(),
	}, nil
}
