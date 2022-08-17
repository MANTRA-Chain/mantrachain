package keeper

import (
	"context"
	"strconv"
	"strings"

	"github.com/LimeChain/mantrachain/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) BuyNft(goCtx context.Context, msg *types.MsgBuyNft) (*types.MsgBuyNftResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator, err := sdk.AccAddressFromBech32(msg.Creator)

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

	collectionSettingsController := NewCollectionSettingsController(ctx, marketplaceIndex, collection.Index, marketplaceId, msg.CollectionId).
		WithStore(k)

	err = collectionSettingsController.
		MustExist().
		Validate()

	if err != nil {
		return nil, err
	}

	collectionSettings := collectionSettingsController.getCollectionSettings()

	nft, found := tokenExecutor.GetNft(collectionCreator, msg.CollectionId, msg.NftId)

	if !found {
		return nil, sdkerrors.Wrap(types.ErrNftDoesNotExist, "invalid or non-existent nft")
	}

	nftSettings, found := k.GetNftSettings(ctx, marketplaceIndex, collection.Index, nft.Index)

	nftsEarningsOnSale := collectionSettings.NftsEarningsOnSale
	var nftsVaultLockPercentage sdk.Int
	initialSale := false

	if !found || !nftSettings.InitiallySold {
		initialSale = true
		nftsVaultLockPercentage = *collectionSettings.InitiallyNftsVaultLockPercentage
	}

	if !found {
		nftSettings = types.NftSettings{
			Index:            nft.Index,
			MarketplaceIndex: marketplaceIndex,
			CollectionIndex:  collection.Index,
			MinPrice:         collectionSettings.InitiallyCollectionOwnerNftsMinPrice,
			ForSale:          collectionSettings.InitiallyCollectionOwnerNftsForSale,
			Creator:          creator,
		}
	}

	if !nftSettings.ForSale {
		return nil, sdkerrors.Wrap(types.ErrNftNotForSale, "nft is not for sale")
	}

	if !nftSettings.InitiallySold {
		nftSettings.InitiallySold = true
	}

	nftSettings.ForSale = false

	k.SetNftSettings(ctx, nftSettings)

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)
	owner := nftExecutor.GetNftOwner(string(collection.Index), string(nft.Index))

	if owner.Equals(creator) {
		return nil, sdkerrors.Wrap(types.ErrInvalidNftBuyer, "nft is owned by the buyer")
	}

	staked, err := k.CollectFeesAndDelegateStake(
		ctx,
		nftSettings.MinPrice,
		nftsEarningsOnSale,
		nftsVaultLockPercentage,
		creator,
		collection.Owner,
		owner,
		marketplaceIndex,
		collection.Index,
		nft.Index,
		initialSale,
	)

	if err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to collect fees and stake")
	}

	err = nftExecutor.TransferNft(string(collection.Index), string(nft.Index), creator)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgBuyNft),
			sdk.NewAttribute(types.AttributeKeyMarketplaceCreator, marketplaceCreator.String()),
			sdk.NewAttribute(types.AttributeKeyMarketplaceId, marketplaceId),
			sdk.NewAttribute(types.AttributeKeyCollectionCreator, collection.Owner.String()),
			sdk.NewAttribute(types.AttributeKeyCollectionId, msg.CollectionId),
			sdk.NewAttribute(types.AttributeKeySigner, creator.String()),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, creator.String()),
			sdk.NewAttribute(types.AttributeKeyStaked, strconv.FormatBool(staked)),
		),
	)

	return &types.MsgBuyNftResponse{
		MarketplaceId:      marketplaceId,
		MarketplaceCreator: marketplaceCreator.String(),
		CollectionId:       msg.CollectionId,
		CollectionCreator:  collection.Owner.String(),
		NftId:              msg.NftId,
		Owner:              owner.String(),
		Receiver:           creator.String(),
		Staked:             staked,
	}, nil
}
