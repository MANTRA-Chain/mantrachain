package keeper

import (
	"context"
	"strings"

	"github.com/LimeChain/mantrachain/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) ImportNftCollection(goCtx context.Context, msg *types.MsgImportNftCollection) (*types.MsgImportNftCollectionResponse, error) {
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
	nftCollection, found := tokenExecutor.GetNftCollection(collectionCreator, msg.CollectionId)

	if !found {
		return nil, sdkerrors.Wrap(types.ErrCollectionDoesNotExist, "invalid or non-existent nft collection")
	}

	if !owner.Equals(nftCollection.Owner) {
		return nil, sdkerrors.Wrap(types.ErrUnauthorized, "not a nft collection owner")
	}

	approved := tokenExecutor.HasApprovedAllNfts(owner, k.GetAddress(ctx), true)

	if !approved {
		return nil, sdkerrors.Wrap(types.ErrNftsNotApproved, "nfts not approved for marketplace")
	}

	marketplaceCollectionController := NewMarketplaceCollectionController(ctx, marketplaceIndex, nftCollection.Index, marketplaceId, msg.CollectionId).
		WithCollection(msg.Collection).
		WithStore(k).
		WithConfiguration(k.GetParams(ctx))

	err = marketplaceCollectionController.
		MustNotExist().
		ValidCollection().
		Validate()

	if err != nil {
		return nil, err
	}

	nftCollectionIndex := marketplaceCollectionController.getIndex()

	parsed, err := sdk.ParseCoinNormalized(msg.Collection.InitiallyNftCollectionOwnerNftsMinPrice)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidInitiallyNftMinPrice, "initially nft min price is invalid")
	}

	var cw20ContractAddress sdk.AccAddress

	if strings.TrimSpace(msg.Collection.Cw20ContractAddress) != "" {
		cw20ContractAddress, err = sdk.AccAddressFromBech32(msg.Collection.Cw20ContractAddress)

		if err != nil {
			return nil, err
		}
	}

	newMarketplaceCollection := types.MarketplaceCollection{
		Index:                                   nftCollectionIndex,
		MarketplaceIndex:                        marketplaceIndex,
		CollectionCreator:                       collectionCreator.String(),
		CollectionId:                            msg.CollectionId,
		InitiallyNftCollectionOwnerNftsForSale:  msg.Collection.InitiallyNftCollectionOwnerNftsForSale,
		InitiallyNftCollectionOwnerNftsMinPrice: &parsed,
		Cw20ContractAddress:                     cw20ContractAddress,
		NftsEarningsOnSale:                      msg.Collection.NftsEarningsOnSale,
		NftsEarningsOnYieldReward:               msg.Collection.NftsEarningsOnYieldReward,
		InitiallyNftsVaultLockPercentage:        msg.Collection.InitiallyNftsVaultLockPercentage,
		Creator:                                 owner,
	}

	k.SetMarketplaceCollection(ctx, newMarketplaceCollection)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgImportNftCollection),
			sdk.NewAttribute(types.AttributeKeyMarketplaceCreator, marketplaceCreator.String()),
			sdk.NewAttribute(types.AttributeKeyMarketplaceId, marketplaceId),
			sdk.NewAttribute(types.AttributeKeyCollectionCreator, nftCollection.Owner.String()),
			sdk.NewAttribute(types.AttributeKeyCollectionId, msg.CollectionId),
			sdk.NewAttribute(types.AttributeKeySigner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
		),
	)

	return &types.MsgImportNftCollectionResponse{
		MarketplaceId:    marketplaceId,
		MarketplaceOwner: marketplaceCreator.String(),
		CollectionId:     msg.CollectionId,
		CollectionOwner:  nftCollection.Owner.String(),
	}, nil
}
