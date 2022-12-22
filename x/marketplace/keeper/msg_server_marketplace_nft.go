package keeper

import (
	"context"
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
	nftCollection, found := tokenExecutor.GetNftCollection(collectionCreator, msg.CollectionId)

	if !found {
		return nil, sdkerrors.Wrap(types.ErrCollectionDoesNotExist, "invalid or non-existent nft collection")
	}

	marketplaceCollection := NewMarketplaceCollectionController(ctx, marketplaceIndex, nftCollection.Index, marketplaceId, msg.CollectionId).
		WithStore(k)

	err = marketplaceCollection.
		MustExist().
		Validate()

	if err != nil {
		return nil, err
	}

	collection := marketplaceCollection.getMarketplaceCollection()

	nft, found := tokenExecutor.GetNft(collectionCreator, msg.CollectionId, msg.NftId)

	if !found {
		return nil, sdkerrors.Wrap(types.ErrNftDoesNotExist, "invalid or non-existent nft")
	}

	marketplaceNft, found := k.GetMarketplaceNft(ctx, marketplaceIndex, collection.Index, nft.Index)

	nftsEarningsOnSale := collection.NftsEarningsOnSale
	nftEarningsOnYieldReward := collection.NftsEarningsOnYieldReward

	var initiallyNftsVaultLockPercentage sdk.Int
	initialSale := false

	if !found || !marketplaceNft.InitiallySold {
		initialSale = true
		initiallyNftsVaultLockPercentage = *collection.InitiallyNftsVaultLockPercentage
	}

	if !found {
		marketplaceNft = types.MarketplaceNft{
			Index:               nft.Index,
			MarketplaceIndex:    marketplaceIndex,
			CollectionIndex:     collection.Index,
			MinPrice:            collection.InitiallyNftCollectionOwnerNftsMinPrice,
			Cw20ContractAddress: collection.Cw20ContractAddress,
			ForSale:             collection.InitiallyNftCollectionOwnerNftsForSale,
			Creator:             creator,
		}
	}

	if !marketplaceNft.ForSale {
		return nil, sdkerrors.Wrap(types.ErrNftNotForSale, "nft is not for sale")
	}

	if !marketplaceNft.InitiallySold {
		marketplaceNft.InitiallySold = true
	}

	marketplaceNft.ForSale = false

	k.SetMarketplaceNft(ctx, marketplaceNft)

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)
	owner := nftExecutor.GetNftOwner(string(nftCollection.Index), string(nft.Index))

	if owner == nil || owner.Empty() {
		return nil, sdkerrors.Wrap(types.ErrUnavailable, "nft owner not found")
	}

	if owner.Equals(creator) {
		return nil, sdkerrors.Wrap(types.ErrInvalidNftBuyer, "nft is owned by the buyer")
	}

	lockCoin, err := k.CollectFees(
		ctx,
		marketplaceNft.MinPrice,
		nftsEarningsOnSale,
		initiallyNftsVaultLockPercentage,
		creator,
		owner,
		initialSale,
		collection.Cw20ContractAddress,
	)

	if !lockCoin.IsZero() {
		vaultExecutor := NewVaultExecutor(ctx, k.vaultKeeper)
		err = vaultExecutor.CreateAndDelegateNftStakeStaked(
			msg.MarketplaceCreator,
			msg.MarketplaceId,
			msg.CollectionCreator,
			msg.CollectionId,
			msg.NftId,
			marketplaceIndex,
			nftCollection.Index,
			nft.Index,
			creator,
			lockCoin,
			msg.StakingChain,
			msg.StakingValidator,
			collection.Cw20ContractAddress,
			nftEarningsOnYieldReward,
		)
	}

	if err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to collect fees and stake")
	}

	// TODO: transfer the current yield reward to the owner if the nft is staked

	err = tokenExecutor.TransferNft(
		k.ac.GetModuleAddress(types.ModuleName),
		owner,
		creator,
		collectionCreator,
		nftCollection.Index,
		nft.Index,
	)

	if err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to transfer nft")
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgBuyNft),
			sdk.NewAttribute(types.AttributeKeyMarketplaceCreator, marketplaceCreator.String()),
			sdk.NewAttribute(types.AttributeKeyMarketplaceId, marketplaceId),
			sdk.NewAttribute(types.AttributeKeyCollectionCreator, nftCollection.Owner.String()),
			sdk.NewAttribute(types.AttributeKeyCollectionId, msg.CollectionId),
			sdk.NewAttribute(types.AttributeKeyNftId, msg.NftId),
			sdk.NewAttribute(types.AttributeKeySigner, creator.String()),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, creator.String()),
			sdk.NewAttribute(types.AttributeKeyStakingChain, msg.StakingChain),
			sdk.NewAttribute(types.AttributeKeyStakingValidator, msg.StakingValidator),
		),
	)

	return &types.MsgBuyNftResponse{
		MarketplaceId:      marketplaceId,
		MarketplaceCreator: marketplaceCreator.String(),
		CollectionId:       msg.CollectionId,
		CollectionCreator:  nftCollection.Owner.String(),
		NftId:              msg.NftId,
		Owner:              owner.String(),
		Receiver:           creator.String(),
		StakingChain:       msg.StakingChain,
		StakingValidator:   msg.StakingValidator,
	}, nil
}
