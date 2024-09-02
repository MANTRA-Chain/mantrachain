package keeper

import (
	"context"

	"cosmossdk.io/errors"
	nft "cosmossdk.io/x/nft"
	"github.com/MANTRA-Finance/mantrachain/x/token/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateNftCollection(goCtx context.Context, msg *types.MsgCreateNftCollection) (*types.MsgCreateNftCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator, err := sdk.AccAddressFromBech32(msg.GetCreator())
	if err != nil {
		return nil, err
	}

	if err := k.guardKeeper.CheckNewRestrictedNftsCollection(ctx, msg.Collection.RestrictedNfts, creator.String()); err != nil {
		return nil, errors.Wrap(err, "unauthorized")
	}

	collectionController := NewNftCollectionController(ctx, creator, false).
		WithMetadata(msg.Collection).
		WithStore(k).
		WithConfiguration(k.GetParams(ctx))

	err = collectionController.
		MustNotExist().
		MustNotBeDefault().
		ValidMetadata().
		Validate()
	if err != nil {
		return nil, err
	}

	collectionIndex := collectionController.getIndex()
	collectionId := collectionController.getId()

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)
	err = nftExecutor.SetClass(nft.Class{
		Id:          string(collectionIndex),
		Name:        msg.Collection.Name,
		Symbol:      msg.Collection.Symbol,
		Description: msg.Collection.Description,
		Uri:         types.ModuleName,
		UriHash:     collectionId,
		Data:        msg.Collection.Data,
	})
	if err != nil {
		return nil, err
	}

	newNftCollection := types.NftCollection{
		Index:    collectionIndex,
		Id:       collectionId,
		Images:   msg.Collection.Images,
		Url:      msg.Collection.Url,
		Links:    msg.Collection.Links,
		Category: msg.Collection.Category,
		Options:  msg.Collection.Options,
		Creator:  creator,
	}

	k.SetNftCollection(ctx, newNftCollection)
	k.SetNftCollectionOwner(ctx, collectionIndex, creator)

	if msg.Collection.SoulBondedNfts {
		k.SetSoulBondedNftsCollection(ctx, collectionIndex)
	}
	if msg.Collection.Opened {
		k.SetOpenedNftsCollection(ctx, collectionIndex)
	}
	// It can be only set by an admin or another authorized address
	if msg.Collection.RestrictedNfts {
		k.SetRestrictedNftsCollection(ctx, collectionIndex)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgCreateNftCollection),
			sdk.NewAttribute(types.AttributeKeyNftCollectionId, collectionId),
			sdk.NewAttribute(types.AttributeKeySigner, creator.String()),
			sdk.NewAttribute(types.AttributeKeyOwner, creator.String()),
		),
	)

	return &types.MsgCreateNftCollectionResponse{
		CollectionId:      collectionId,
		CollectionCreator: creator.String(),
	}, nil
}
