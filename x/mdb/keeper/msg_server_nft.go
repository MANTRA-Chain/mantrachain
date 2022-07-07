package keeper

import (
	"context"
	"strconv"
	"strings"

	"github.com/LimeChain/mantrachain/x/mdb/types"
	nfttypes "github.com/LimeChain/mantrachain/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k msgServer) MintNfts(goCtx context.Context, msg *types.MsgMintNfts) (*types.MsgMintNftsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator, err := sdk.AccAddressFromBech32(msg.Creator)

	if err != nil {
		return nil, err
	}

	if msg.Receiver == "" {
		msg.Receiver = msg.Creator
	}

	receiver, err := sdk.AccAddressFromBech32(msg.Receiver)

	if err != nil {
		return nil, err
	}

	if len(msg.Nfts.Nfts) == 0 {
		return nil, sdkerrors.Wrapf(types.ErrInvalidNftsCount, "nfts length %d invalid, min 1", len(msg.Nfts.Nfts))
	}

	var collectionCreator sdk.AccAddress

	if msg.CollectionCreator == "" {
		msg.CollectionCreator = msg.Creator
		collectionCreator = creator
	} else {
		collectionCreator, err = sdk.AccAddressFromBech32(msg.CollectionCreator)

		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid collection creator")
		}
	}

	collectionController := NewNftCollectionController(ctx, collectionCreator).
		WithId(msg.CollectionId).
		WithStore(k).
		WithConfiguration(k.GetParams(ctx))

	err = collectionController.
		CreateDefaultIfNotExists().
		Execute()

	if err != nil {
		return nil, err
	}

	err = collectionController.
		MustExist().
		IsOpenedOrOwner(creator).
		Validate()

	if err != nil {
		return nil, err
	}

	collectionIndex := collectionController.getIndex()
	collectionId := collectionController.getId()

	nftController := NewNftController(ctx, collectionIndex).
		WithMetadata(msg.Nfts.Nfts).
		WithStore(k).
		WithConfiguration(k.GetParams(ctx))

	err = nftController.
		FilterNotExist().
		Execute()

	if err != nil {
		return nil, err
	}

	err = nftController.
		ValidMetadata().
		Validate()

	if err != nil {
		return nil, err
	}

	nftsMetadata := nftController.getMetadata()

	if len(nftsMetadata) == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidNftsCount, "existing nfts")
	}

	var newNfts []nfttypes.NFT
	var newNftsMetadata []types.Nft
	var nftsIds []string

	for _, nftMetadata := range nftsMetadata {
		nftIndex := types.GetNftIndex(collectionIndex, nftMetadata.Id)

		newNfts = append(newNfts, nfttypes.NFT{
			ClassId: string(collectionIndex),
			Id:      string(nftIndex),
			Uri:     types.ModuleName,
			UriHash: nftMetadata.Id,
			Data:    nftMetadata.Data,
		})

		newNftsMetadata = append(newNftsMetadata, types.Nft{
			Index:           nftIndex,
			Images:          nftMetadata.Images,
			Url:             nftMetadata.Url,
			Links:           nftMetadata.Links,
			Title:           nftMetadata.Title,
			Description:     nftMetadata.Description,
			Attributes:      nftMetadata.Attributes,
			Resellable:      nftMetadata.Resellable,
			CollectionIndex: collectionIndex,
			CollectionId:    collectionId,
			Creator:         creator,
		})

		nftsIds = append(nftsIds, string(nftMetadata.Id))
	}

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)
	err = nftExecutor.MintNftBatch(newNfts, receiver)
	if err != nil {
		return nil, err
	}

	k.SetNfts(ctx, newNftsMetadata)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgMintNfts),
			sdk.NewAttribute(types.AttributeKeyNftCollection, collectionId),
			sdk.NewAttribute(types.AttributeKeyNfts, strings.Join(nftsIds[:], ",")),
			sdk.NewAttribute(types.AttributeKeyCreator, creator.String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, receiver.String()),
		),
	)

	return &types.MsgMintNftsResponse{
		Ids:     nftsIds,
		Address: receiver.String(),
	}, nil
}

func (k msgServer) BurnNfts(goCtx context.Context, msg *types.MsgBurnNfts) (*types.MsgBurnNftsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator, err := sdk.AccAddressFromBech32(msg.Creator)

	if err != nil {
		return nil, err
	}

	if len(msg.Nfts.NftsIds) == 0 {
		return nil, sdkerrors.Wrapf(types.ErrInvalidNftsCount, "nfts length %d invalid, min 1", len(msg.Nfts.NftsIds))
	}

	var collectionCreator sdk.AccAddress

	if msg.CollectionCreator == "" {
		msg.CollectionCreator = msg.Creator
		collectionCreator = creator
	} else {
		collectionCreator, err = sdk.AccAddressFromBech32(msg.CollectionCreator)

		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid collection creator")
		}
	}

	collectionController := NewNftCollectionController(ctx, collectionCreator).
		WithId(msg.CollectionId).
		WithStore(k).
		WithConfiguration(k.GetParams(ctx))

	err = collectionController.
		MustExist().
		Validate()

	if err != nil {
		return nil, err
	}

	collectionIndex := collectionController.getIndex()
	collectionId := collectionController.getId()

	nftController := NewNftController(ctx, collectionIndex).
		WithIds(msg.Nfts.NftsIds).
		WithStore(k).
		WithConfiguration(k.GetParams(ctx))

	err = nftController.
		FilterNotExist().
		FilterNotOwn(creator).
		Execute()

	if err != nil {
		return nil, err
	}

	nftsIds := nftController.getNftsIds()

	if len(nftsIds) == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidNftsCount, "not existing nfts or not an owner")
	}

	nftsIndexes := nftController.getIndexes()

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)
	err = nftExecutor.BurnNftBatch(string(collectionIndex), nftsIds)
	if err != nil {
		return nil, err
	}

	k.DeleteNfts(ctx, collectionIndex, nftsIndexes)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgBurnNfts),
			sdk.NewAttribute(types.AttributeKeyNftCollection, collectionId),
			sdk.NewAttribute(types.AttributeKeyNfts, strings.Join(nftsIds[:], ",")),
			sdk.NewAttribute(types.AttributeKeyCreator, creator.String()),
		),
	)

	return &types.MsgBurnNftsResponse{
		Ids: nftsIds,
	}, nil
}

func (k msgServer) ApproveNfts(goCtx context.Context, msg *types.MsgApproveNfts) (*types.MsgApproveNftsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator, err := sdk.AccAddressFromBech32(msg.Creator)

	if err != nil {
		return nil, err
	}

	receiver, err := sdk.AccAddressFromBech32(msg.Receiver)

	if err != nil {
		return nil, err
	}

	if len(msg.Nfts.NftsIds) == 0 {
		return nil, sdkerrors.Wrapf(types.ErrInvalidNftsCount, "nfts length %d invalid, min 1", len(msg.Nfts.NftsIds))
	}

	var collectionCreator sdk.AccAddress

	if msg.CollectionCreator == "" {
		msg.CollectionCreator = msg.Creator
		collectionCreator = creator
	} else {
		collectionCreator, err = sdk.AccAddressFromBech32(msg.CollectionCreator)

		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid collection creator")
		}
	}

	collectionController := NewNftCollectionController(ctx, collectionCreator).
		WithId(msg.CollectionId).
		WithStore(k).
		WithConfiguration(k.GetParams(ctx))

	err = collectionController.
		MustExist().
		Validate()

	if err != nil {
		return nil, err
	}

	collectionIndex := collectionController.getIndex()
	collectionId := collectionController.getId()

	nftController := NewNftController(ctx, collectionIndex).
		WithIds(msg.Nfts.NftsIds).
		WithStore(k).
		WithConfiguration(k.GetParams(ctx))

	err = nftController.
		FilterNotExist().
		FilterNotOwn(creator).
		Execute()

	if err != nil {
		return nil, err
	}

	nftsIds := nftController.getNftsIds()

	if len(nftsIds) == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidNftsCount, "not existing nfts or not an owner")
	}

	nftsIndexes := nftController.getIndexes()

	k.SetApprovedNfts(ctx, collectionIndex, nftsIndexes, creator, receiver, msg.Approved)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgApproveNfts),
			sdk.NewAttribute(types.AttributeKeyNftCollection, collectionId),
			sdk.NewAttribute(types.AttributeKeyNfts, strings.Join(nftsIds[:], ",")),
			sdk.NewAttribute(types.AttributeKeyApproved, strconv.FormatBool(msg.Approved)),
			sdk.NewAttribute(types.AttributeKeyCreator, creator.String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, receiver.String()),
		),
	)

	return &types.MsgApproveNftsResponse{
		Ids:      nftsIds,
		Address:  receiver.String(),
		Approved: msg.Approved,
	}, nil
}

func (k msgServer) ApproveAllNfts(goCtx context.Context, msg *types.MsgApproveAllNfts) (*types.MsgApproveAllNftsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator, err := sdk.AccAddressFromBech32(msg.Creator)

	if err != nil {
		return nil, err
	}

	receiver, err := sdk.AccAddressFromBech32(msg.Receiver)

	if err != nil {
		return nil, err
	}

	k.SetApprovedAllNfts(ctx, creator, receiver, msg.Approved)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgApproveAllNfts),
			sdk.NewAttribute(types.AttributeKeyApproved, strconv.FormatBool(msg.Approved)),
			sdk.NewAttribute(types.AttributeKeyCreator, creator.String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, receiver.String()),
		),
	)

	return &types.MsgApproveAllNftsResponse{
		Address:  receiver.String(),
		Approved: msg.Approved,
	}, nil
}

func (k msgServer) MintNft(goCtx context.Context, msg *types.MsgMintNft) (*types.MsgMintNftResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator, err := sdk.AccAddressFromBech32(msg.Creator)

	if err != nil {
		return nil, err
	}

	if msg.Receiver == "" {
		msg.Receiver = msg.Creator
	}

	receiver, err := sdk.AccAddressFromBech32(msg.Receiver)

	if err != nil {
		return nil, err
	}

	if msg.Nft == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidNft, "nft cannot be empty")
	}

	var collectionCreator sdk.AccAddress

	if msg.CollectionCreator == "" {
		msg.CollectionCreator = msg.Creator
		collectionCreator = creator
	} else {
		collectionCreator, err = sdk.AccAddressFromBech32(msg.CollectionCreator)

		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid collection creator")
		}
	}

	collectionController := NewNftCollectionController(ctx, collectionCreator).
		WithId(msg.CollectionId).
		WithStore(k).
		WithConfiguration(k.GetParams(ctx))

	err = collectionController.
		CreateDefaultIfNotExists().
		Execute()

	if err != nil {
		return nil, err
	}

	err = collectionController.
		MustExist().
		IsOpenedOrOwner(creator).
		Validate()

	if err != nil {
		return nil, err
	}

	collectionIndex := collectionController.getIndex()
	collectionId := collectionController.getId()

	nftsMetadata := make([]*types.MsgNftMetadata, 1)
	nftsMetadata = append(nftsMetadata, msg.Nft)

	nftController := NewNftController(ctx, collectionIndex).
		WithMetadata(nftsMetadata).
		WithStore(k).
		WithConfiguration(k.GetParams(ctx))

	err = nftController.
		FilterNotExist().
		Execute()

	if err != nil {
		return nil, err
	}

	err = nftController.
		ValidMetadata().
		Validate()

	if err != nil {
		return nil, err
	}

	nftsMetadata = nftController.getMetadata()

	if len(nftsMetadata) == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidNft, "existing nft")
	}

	nftMetadata := nftsMetadata[0]
	nftIndex := types.GetNftIndex(collectionIndex, nftMetadata.Id)

	newNft := nfttypes.NFT{
		ClassId: string(collectionIndex),
		Id:      string(nftIndex),
		Uri:     types.ModuleName,
		UriHash: nftMetadata.Id,
		Data:    nftMetadata.Data,
	}

	newNftMetadata := types.Nft{
		Index:           nftIndex,
		Images:          nftMetadata.Images,
		Url:             nftMetadata.Url,
		Links:           nftMetadata.Links,
		Title:           nftMetadata.Title,
		Description:     nftMetadata.Description,
		Attributes:      nftMetadata.Attributes,
		Resellable:      nftMetadata.Resellable,
		CollectionIndex: collectionIndex,
		CollectionId:    collectionId,
		Creator:         creator,
	}

	nftId := nftMetadata.Id

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)
	err = nftExecutor.MintNft(newNft, receiver)
	if err != nil {
		return nil, err
	}

	k.SetNft(ctx, newNftMetadata)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgMintNft),
			sdk.NewAttribute(types.AttributeKeyNftCollection, collectionId),
			sdk.NewAttribute(types.AttributeKeyNfts, nftId),
			sdk.NewAttribute(types.AttributeKeyCreator, creator.String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, receiver.String()),
		),
	)

	return &types.MsgMintNftResponse{
		Id:      nftId,
		Address: receiver.String(),
	}, nil
}

func (k msgServer) BurnNft(goCtx context.Context, msg *types.MsgBurnNft) (*types.MsgBurnNftResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator, err := sdk.AccAddressFromBech32(msg.Creator)

	if err != nil {
		return nil, err
	}

	if msg.NftId == "" {
		return nil, sdkerrors.Wrap(types.ErrInvalidNft, "nft id cannot be empty")
	}

	var collectionCreator sdk.AccAddress

	if msg.CollectionCreator == "" {
		msg.CollectionCreator = msg.Creator
		collectionCreator = creator
	} else {
		collectionCreator, err = sdk.AccAddressFromBech32(msg.CollectionCreator)

		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid collection creator")
		}
	}

	collectionController := NewNftCollectionController(ctx, collectionCreator).
		WithId(msg.CollectionId).
		WithStore(k).
		WithConfiguration(k.GetParams(ctx))

	err = collectionController.
		MustExist().
		Validate()

	if err != nil {
		return nil, err
	}

	collectionIndex := collectionController.getIndex()
	collectionId := collectionController.getId()

	nftController := NewNftController(ctx, collectionIndex).
		WithId(msg.NftId).
		WithStore(k).
		WithConfiguration(k.GetParams(ctx))

	err = nftController.
		FilterNotExist().
		FilterNotOwn(creator).
		Execute()

	if err != nil {
		return nil, err
	}

	nftsIds := nftController.getNftsIds()

	if len(nftsIds) == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidNft, "not existing nft or not an owner")
	}

	nftsIndexes := nftController.getIndexes()

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)
	err = nftExecutor.BurnNft(string(collectionIndex), nftsIds[0])
	if err != nil {
		return nil, err
	}

	k.DeleteNft(ctx, collectionIndex, nftsIndexes[0])

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgBurnNft),
			sdk.NewAttribute(types.AttributeKeyNftCollection, collectionId),
			sdk.NewAttribute(types.AttributeKeyNfts, nftsIds[0]),
			sdk.NewAttribute(types.AttributeKeyCreator, creator.String()),
		),
	)

	return &types.MsgBurnNftResponse{
		Id: nftsIds[0],
	}, nil
}

func (k msgServer) ApproveNft(goCtx context.Context, msg *types.MsgApproveNft) (*types.MsgApproveNftResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator, err := sdk.AccAddressFromBech32(msg.Creator)

	if err != nil {
		return nil, err
	}

	receiver, err := sdk.AccAddressFromBech32(msg.Receiver)

	if err != nil {
		return nil, err
	}

	if msg.NftId == "" {
		return nil, sdkerrors.Wrap(types.ErrInvalidNft, "nft id cannot be empty")
	}

	var collectionCreator sdk.AccAddress

	if msg.CollectionCreator == "" {
		msg.CollectionCreator = msg.Creator
		collectionCreator = creator
	} else {
		collectionCreator, err = sdk.AccAddressFromBech32(msg.CollectionCreator)

		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid collection creator")
		}
	}

	collectionController := NewNftCollectionController(ctx, collectionCreator).
		WithId(msg.CollectionId).
		WithStore(k).
		WithConfiguration(k.GetParams(ctx))

	err = collectionController.
		MustExist().
		Validate()

	if err != nil {
		return nil, err
	}

	collectionIndex := collectionController.getIndex()
	collectionId := collectionController.getId()

	nftController := NewNftController(ctx, collectionIndex).
		WithId(msg.NftId).
		WithStore(k).
		WithConfiguration(k.GetParams(ctx))

	err = nftController.
		FilterNotExist().
		FilterNotOwn(creator).
		Execute()

	if err != nil {
		return nil, err
	}

	nftsIds := nftController.getNftsIds()

	if len(nftsIds) == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidNft, "not existing nft or not an owner")
	}

	nftsIndexes := nftController.getIndexes()

	k.SetApprovedNft(ctx, collectionIndex, nftsIndexes[0], creator, receiver, msg.Approved)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgApproveNft),
			sdk.NewAttribute(types.AttributeKeyNftCollection, collectionId),
			sdk.NewAttribute(types.AttributeKeyNfts, nftsIds[0]),
			sdk.NewAttribute(types.AttributeKeyApproved, strconv.FormatBool(msg.Approved)),
			sdk.NewAttribute(types.AttributeKeyCreator, creator.String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, receiver.String()),
		),
	)

	return &types.MsgApproveNftResponse{
		Id:       nftsIds[0],
		Address:  receiver.String(),
		Approved: msg.Approved,
	}, nil
}

func (k msgServer) TransferNft(goCtx context.Context, msg *types.MsgTransferNft) (*types.MsgTransferNftResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator, err := sdk.AccAddressFromBech32(msg.Creator)

	if err != nil {
		return nil, err
	}

	owner, err := sdk.AccAddressFromBech32(msg.Owner)

	if err != nil {
		return nil, err
	}

	receiver, err := sdk.AccAddressFromBech32(msg.Receiver)

	if err != nil {
		return nil, err
	}

	if msg.NftId == "" {
		return nil, sdkerrors.Wrap(types.ErrInvalidNft, "nft id cannot be empty")
	}

	var collectionCreator sdk.AccAddress

	if msg.CollectionCreator == "" {
		msg.CollectionCreator = msg.Creator
		collectionCreator = creator
	} else {
		collectionCreator, err = sdk.AccAddressFromBech32(msg.CollectionCreator)

		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid collection creator")
		}
	}

	collectionController := NewNftCollectionController(ctx, collectionCreator).
		WithId(msg.CollectionId).
		WithStore(k).
		WithConfiguration(k.GetParams(ctx))

	err = collectionController.
		MustExist().
		Validate()

	if err != nil {
		return nil, err
	}

	collectionIndex := collectionController.getIndex()
	collectionId := collectionController.getId()

	nftController := NewNftController(ctx, collectionIndex).
		WithId(msg.NftId).
		WithStore(k).
		WithConfiguration(k.GetParams(ctx))

	err = nftController.
		FilterNotExist().
		FilterNotOwn(owner).
		FilterCannotTransfer(creator).
		Execute()

	if err != nil {
		return nil, err
	}

	nftsIds := nftController.getNftsIds()

	if len(nftsIds) == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidNft, "not existing nft or no transfer permission")
	}

	nftsIndexes := nftController.getIndexes()

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)
	err = nftExecutor.TransferNft(string(collectionIndex), string(nftsIndexes[0]), receiver)
	if err != nil {
		return nil, err
	}

	k.DeleteApprovedNft(ctx, collectionIndex, nftsIndexes[0], receiver)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgTransferNft),
			sdk.NewAttribute(types.AttributeKeyNftCollection, collectionId),
			sdk.NewAttribute(types.AttributeKeyNfts, nftsIds[0]),
			sdk.NewAttribute(types.AttributeKeyCreator, creator.String()),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, receiver.String()),
		),
	)

	return &types.MsgTransferNftResponse{
		Id:      nftsIds[0],
		Address: receiver.String(),
	}, nil
}

func (k msgServer) TransferNfts(goCtx context.Context, msg *types.MsgTransferNfts) (*types.MsgTransferNftsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator, err := sdk.AccAddressFromBech32(msg.Creator)

	if err != nil {
		return nil, err
	}

	owner, err := sdk.AccAddressFromBech32(msg.Owner)

	if err != nil {
		return nil, err
	}

	receiver, err := sdk.AccAddressFromBech32(msg.Receiver)

	if err != nil {
		return nil, err
	}

	if len(msg.Nfts.NftsIds) == 0 {
		return nil, sdkerrors.Wrapf(types.ErrInvalidNftsCount, "nfts length %d invalid, min 1", len(msg.Nfts.NftsIds))
	}

	var collectionCreator sdk.AccAddress

	if msg.CollectionCreator == "" {
		msg.CollectionCreator = msg.Creator
		collectionCreator = creator
	} else {
		collectionCreator, err = sdk.AccAddressFromBech32(msg.CollectionCreator)

		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid collection creator")
		}
	}

	collectionController := NewNftCollectionController(ctx, collectionCreator).
		WithId(msg.CollectionId).
		WithStore(k).
		WithConfiguration(k.GetParams(ctx))

	err = collectionController.
		MustExist().
		Validate()

	if err != nil {
		return nil, err
	}

	collectionIndex := collectionController.getIndex()
	collectionId := collectionController.getId()

	nftController := NewNftController(ctx, collectionIndex).
		WithIds(msg.Nfts.NftsIds).
		WithStore(k).
		WithConfiguration(k.GetParams(ctx))

	err = nftController.
		FilterNotExist().
		FilterNotOwn(owner).
		FilterCannotTransfer(creator).
		Execute()

	if err != nil {
		return nil, err
	}

	nftsIds := nftController.getNftsIds()

	if len(nftsIds) == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidNftsCount, "not existing nfts or no transfer permission")
	}

	nftsIndexes := nftController.getIndexes()

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)
	err = nftExecutor.TransferNftBatch(string(collectionIndex), nftsIds, receiver)
	if err != nil {
		return nil, err
	}

	k.DeleteApprovedNfts(ctx, collectionIndex, nftsIndexes, receiver)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgTransferNfts),
			sdk.NewAttribute(types.AttributeKeyNftCollection, collectionId),
			sdk.NewAttribute(types.AttributeKeyNfts, strings.Join(nftsIds[:], ",")),
			sdk.NewAttribute(types.AttributeKeyCreator, creator.String()),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, receiver.String()),
		),
	)

	return &types.MsgTransferNftsResponse{
		Ids:     nftsIds,
		Address: receiver.String(),
	}, nil
}
