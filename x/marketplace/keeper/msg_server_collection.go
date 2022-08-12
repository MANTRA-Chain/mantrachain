package keeper

import (
	"context"
	"strings"

	"github.com/LimeChain/mantrachain/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) ImportCollection(goCtx context.Context, msg *types.MsgImportCollection) (*types.MsgImportCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := sdk.AccAddressFromBech32(msg.Creator)

	if err != nil {
		return nil, err
	}

	marketplaceCreator, err := sdk.AccAddressFromBech32(msg.MarketplaceCreator)

	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(msg.MarketplaceId) == "" {
		return nil, sdkerrors.Wrap(types.ErrInvalidMarketplaceId, "marketplace id should not be empty")
	}

	collectionCreator, err := sdk.AccAddressFromBech32(msg.CollectionCreator)

	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(msg.CollectionId) == "" {
		return nil, sdkerrors.Wrap(types.ErrInvalidCollectionId, "marketplace id should not be empty")
	}

	marketplaceController := NewMarketplaceController(ctx, marketplaceCreator).
		WithId(msg.MarketplaceId).
		WithStore(k)

	err = marketplaceController.
		MustExist().
		IsOpenedOrHasOwner(owner).
		Validate()

	if err != nil {
		return nil, err
	}

	marketplaceIndex := marketplaceController.getIndex()
	marketplaceId := marketplaceController.getId()

	tokenExecutor := NewTokenExecutor(ctx, k.tokenKeeper)
	collection, found := tokenExecutor.GetCollection(collectionCreator, msg.CollectionId)

	if !found {
		return nil, sdkerrors.Wrap(types.ErrCollectionDoesNotExist, "invalid or non-existent nft collection")
	}

	if !owner.Equals(collection.Owner) {
		return nil, sdkerrors.Wrap(types.ErrUnauthorized, "not a nft collection owner")
	}

	collectionSettingsController := NewCollectionSettingsController(ctx, marketplaceIndex, collection.Index, marketplaceId, msg.CollectionId).
		WithSettings(msg.Settings).
		WithStore(k).
		WithConfiguration(k.GetParams(ctx))

	err = collectionSettingsController.
		MustNotExist().
		ValidSettings().
		Validate()

	if err != nil {
		return nil, err
	}

	collectionSettingsIndex := collectionSettingsController.getIndex()

	parsed, err := sdk.ParseCoinNormalized(msg.Settings.InitiallyCollectionOwnerNftsMinPrice)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidInitiallyNftMinPrice, "initially nft min price is invalid")
	}

	newCollectionSettings := types.CollectionSettings{
		Index:                                collectionSettingsIndex,
		MarketplaceIndex:                     marketplaceIndex,
		InitiallyCollectionOwnerNftsForSale:  msg.Settings.InitiallyCollectionOwnerNftsForSale,
		InitiallyCollectionOwnerNftsMinPrice: &parsed,
		NftsEarningsOnSale:                   msg.Settings.NftsEarningsOnSale,
		NftsEarningsOnYieldReward:            msg.Settings.NftsEarningsOnYieldReward,
		InitiallyNftsVaultLockPercentage:     msg.Settings.InitiallyNftsVaultLockPercentage,
		Creator:                              owner,
	}

	k.SetCollectionSettings(ctx, newCollectionSettings)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgImportCollection),
			sdk.NewAttribute(types.AttributeKeyMarketplaceCreator, marketplaceCreator.String()),
			sdk.NewAttribute(types.AttributeKeyMarketplaceId, marketplaceId),
			sdk.NewAttribute(types.AttributeKeyCollectionCreator, collection.Owner.String()),
			sdk.NewAttribute(types.AttributeKeyCollectionId, msg.CollectionId),
			sdk.NewAttribute(types.AttributeKeySigner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
		),
	)

	return &types.MsgImportCollectionResponse{
		MarketplaceId:    marketplaceId,
		MarketplaceOwner: marketplaceCreator.String(),
		CollectionId:     msg.CollectionId,
		CollectionOwner:  collection.Owner.String(),
	}, nil
}
