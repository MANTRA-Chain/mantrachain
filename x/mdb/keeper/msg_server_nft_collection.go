package keeper

import (
	"context"

	"github.com/LimeChain/mantrachain/x/mdb/types"
	nfttypes "github.com/LimeChain/mantrachain/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateNftCollection(goCtx context.Context, msg *types.MsgCreateNftCollection) (*types.MsgCreateNftCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := sdk.AccAddressFromBech32(msg.Creator)

	if err != nil {
		return nil, err
	}

	if msg.Collection.Id == "" {
		return nil, sdkerrors.Wrap(types.ErrInvalidNftCollectionId, "collection id should not be empty")
	}

	collectionController := NewNftCollectionController(ctx, owner).
		WithMetadata(msg.Collection).
		WithStore(k).
		WithConfiguration(k.GetParams(ctx))

	err = collectionController.
		MustNotExist().
		MustNotBeDefault().
		ValidCollectionMetadata().
		Validate()

	if err != nil {
		return nil, err
	}

	collectionIndex := collectionController.getIndex()
	collectionId := collectionController.getId()

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)
	err = nftExecutor.SetClass(nfttypes.Class{
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
		Images:   msg.Collection.Images,
		Url:      msg.Collection.Url,
		Links:    msg.Collection.Links,
		Category: msg.Collection.Category,
		Options:  msg.Collection.Options,
		Opened:   msg.Collection.Opened,
		Creator:  owner,
		Owner:    owner,
	}

	k.SetNftCollection(ctx, newNftCollection)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgCreateNftCollection),
			sdk.NewAttribute(types.AttributeKeyNftCollection, collectionId),
			sdk.NewAttribute(types.AttributeKeyCreator, owner.String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, owner.String()),
		),
	)

	return &types.MsgCreateNftCollectionResponse{
		Id: collectionId,
	}, nil
}
