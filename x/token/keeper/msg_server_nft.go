package keeper

import (
	"context"
	"strconv"
	"strings"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"mantrachain/x/token/types"
	"mantrachain/x/token/utils"

	nft "github.com/cosmos/cosmos-sdk/x/nft"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k msgServer) UpdateGuardSoulBondNftImage(goCtx context.Context, msg *types.MsgUpdateGuardSoulBondNftImage) (*types.MsgUpdateGuardSoulBondNftImageResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	conf := k.GetParams(ctx)

	if err := k.gk.CheckIsAdmin(ctx, msg.Creator); err != nil {
		return nil, sdkerrors.Wrap(err, "unauthorized")
	}

	collectionCreator, collectionId := k.gk.GetAccountPrivilegesTokenCollectionCreatorAndCollectionId(ctx)

	collectionCreatorAddr, err := sdk.AccAddressFromBech32(collectionCreator)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid collection creator")
	}

	if strings.TrimSpace(collectionId) == "" {
		return nil, sdkerrors.Wrap(types.ErrInvalidNftCollectionId, "nft collection id should not be empty")
	}

	collectionIndex := types.GetNftCollectionIndex(collectionCreatorAddr, collectionId)
	nftIndex := types.GetNftIndex(collectionIndex, msg.NftId)

	nft, found := k.GetNft(ctx, collectionIndex, nftIndex)

	if !found {
		return nil, sdkerrors.Wrap(types.ErrInvalidNft, "nft not found")
	}

	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	nftOwner := k.nftKeeper.GetOwner(ctx, string(collectionIndex), string(nftIndex))

	if nftOwner.Empty() || !owner.Equals(nftOwner) {
		return nil, sdkerrors.Wrap(types.ErrInvalidNft, "nft owner invalid")
	}

	if msg.Index > 0 && (nft.Images == nil || int(msg.Index) >= len(nft.Images)) {
		return nil, sdkerrors.Wrap(types.ErrInvalidNftImageIndex, "invalid nft image index")
	}

	if msg.Image.Image.Type == "" || int32(len(msg.Image.Image.Type)) > conf.ValidNftMetadataImagesTypeMaxLength {
		return nil, sdkerrors.Wrapf(types.ErrInvalidNftImage, "nft id %s image type empty or too long, max %d, image index %d", msg.NftId, conf.ValidNftMetadataImagesTypeMaxLength, msg.Index)
	}

	if !utils.IsUrl(msg.Image.Image.Url) {
		return nil, sdkerrors.Wrapf(types.ErrInvalidNftImage, "nft id %s image invalid url, image index %d", msg.NftId, msg.Index)
	}

	if msg.Index == 0 && (nft.Images == nil || len(nft.Images) == 0) {
		nft.Images = []*types.TokenImage{
			{
				Type: msg.Image.Image.Type,
				Url:  msg.Image.Image.Url,
			},
		}
	} else {
		nft.Images[msg.Index] = &types.TokenImage{
			Type: msg.Image.Image.Type,
			Url:  msg.Image.Image.Url,
		}
	}

	k.SetNft(ctx, nft)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgUpdateGuardSoulBondNftImage),
			sdk.NewAttribute(types.AttributeKeyNftCollectionCreator, collectionCreator),
			sdk.NewAttribute(types.AttributeKeyNftCollectionId, collectionId),
			sdk.NewAttribute(types.AttributeKeyNftId, msg.NftId),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
		),
	)

	return &types.MsgUpdateGuardSoulBondNftImageResponse{
		NftId:             msg.NftId,
		Owner:             owner.String(),
		CollectionCreator: collectionCreator,
		CollectionId:      collectionId,
	}, nil
}

func (k msgServer) MintNfts(goCtx context.Context, msg *types.MsgMintNfts) (*types.MsgMintNftsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(msg.Receiver) == "" {
		msg.Receiver = msg.Creator
	}

	receiver, err := sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return nil, err
	}

	var collectionCreator sdk.AccAddress

	if strings.TrimSpace(msg.CollectionCreator) == "" && !msg.Strict {
		msg.CollectionCreator = msg.Creator
		collectionCreator = creator
	} else {
		collectionCreator, err = sdk.AccAddressFromBech32(msg.CollectionCreator)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid collection creator")
		}
	}

	if err := k.gk.CheckRestrictedNftsCollection(ctx, collectionCreator.String(), msg.CollectionId, msg.GetCreator()); err != nil {
		return nil, sdkerrors.Wrap(err, "unauthorized")
	}

	collectionController := NewNftCollectionController(ctx, collectionCreator, msg.Strict).
		WithId(msg.CollectionId).
		WithStore(k)

	if err := collectionController.CreateDefaultIfNotExists().Execute(); err != nil {
		return nil, err
	}

	collectionController.MustExist()

	restrictedCollection := k.HasRestrictedNftsCollection(
		ctx,
		collectionController.getIndex(),
	)

	if !restrictedCollection {
		collectionController.OpenedOrOwner(creator)
	}

	if msg.Did {
		soulBondNftsCollection := k.HasSoulBondedNftsCollection(
			ctx,
			collectionController.getIndex(),
		)

		if !restrictedCollection || !soulBondNftsCollection {
			return nil, sdkerrors.Wrap(types.ErrInvalidDid, "cannot use did for nfts for collection which is not restricted and/or not for soul-bond nfts")
		}
	}

	if err := collectionController.Validate(); err != nil {
		return nil, err
	}

	collectionIndex := collectionController.getIndex()
	collectionId := collectionController.getId()

	nftController := NewNftController(ctx, collectionIndex).
		WithMetadata(msg.Nfts.Nfts).
		WithStore(k).
		WithConfiguration(k.GetParams(ctx))

	if err := nftController.FilterExist().Execute(); err != nil {
		return nil, err
	}

	if err := nftController.ValidMetadata().Validate(); err != nil {
		return nil, err
	}

	nftsMetadata := nftController.getMetadata()
	if len(nftsMetadata) == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidNftsCount, "existing nfts")
	}

	var newNfts []nft.NFT
	var newNftsMetadata []types.Nft
	var nftsIds []string

	didExecutor := NewDidExecutor(ctx, receiver.String(), k.dk)

	for _, nftMetadata := range nftsMetadata {
		nftIndex := types.GetNftIndex(collectionIndex, nftMetadata.Id)

		newNfts = append(newNfts, nft.NFT{
			ClassId: string(collectionIndex),
			Id:      string(nftIndex),
			Uri:     types.ModuleName,
			UriHash: nftMetadata.Id,
			Data:    nftMetadata.Data,
		})

		var did string

		if msg.Did {
			did, err = didExecutor.CreateDidForNft(nftIndex)
			if err != nil {
				return nil, err
			}
		}

		newNftsMetadata = append(newNftsMetadata, types.Nft{
			Index:             nftIndex,
			Id:                nftMetadata.Id,
			Images:            nftMetadata.Images,
			Url:               nftMetadata.Url,
			Links:             nftMetadata.Links,
			Title:             nftMetadata.Title,
			Description:       nftMetadata.Description,
			Attributes:        nftMetadata.Attributes,
			CollectionIndex:   collectionIndex,
			CollectionId:      collectionId,
			CollectionCreator: collectionCreator,
			Creator:           creator,
			Did:               did,
		})

		nftsIds = append(nftsIds, string(nftMetadata.Id))
	}

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)
	if err := nftExecutor.BatchMintNft(newNfts, receiver); err != nil {
		return nil, err
	}

	k.SetNfts(ctx, newNftsMetadata)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgMintNfts),
			sdk.NewAttribute(types.AttributeKeyNftCollectionCreator, collectionCreator.String()),
			sdk.NewAttribute(types.AttributeKeyNftCollectionId, collectionId),
			sdk.NewAttribute(types.AttributeKeyNftsIds, strings.Join(nftsIds, ",")),
			sdk.NewAttribute(types.AttributeKeySigner, creator.String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, receiver.String()),
		),
	)

	return &types.MsgMintNftsResponse{
		NftsIds:           nftsIds,
		Creator:           creator.String(),
		Receiver:          receiver.String(),
		CollectionCreator: collectionCreator.String(),
		CollectionId:      collectionId,
	}, nil
}

func (k msgServer) BurnNfts(goCtx context.Context, msg *types.MsgBurnNfts) (*types.MsgBurnNftsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	var collectionCreator sdk.AccAddress

	if strings.TrimSpace(msg.CollectionCreator) == "" && !msg.Strict {
		msg.CollectionCreator = msg.Creator
		collectionCreator = owner
	} else {
		collectionCreator, err = sdk.AccAddressFromBech32(msg.CollectionCreator)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid collection creator")
		}
	}

	if err := k.gk.CheckRestrictedNftsCollection(ctx, collectionCreator.String(), msg.CollectionId, msg.GetCreator()); err != nil {
		return nil, sdkerrors.Wrap(err, "unauthorized")
	}

	collectionController := NewNftCollectionController(ctx, collectionCreator, msg.Strict).
		WithId(msg.CollectionId).
		WithStore(k)

	if err := collectionController.MustExist().Validate(); err != nil {
		return nil, err
	}

	collectionIndex := collectionController.getIndex()
	collectionId := collectionController.getId()

	nftController := NewNftController(ctx, collectionIndex).
		WithIds(msg.Nfts.NftsIds).
		WithStore(k).
		WithConfiguration(k.GetParams(ctx))

	if err := nftController.ValidMetadataMaxCount().Validate(); err != nil {
		return nil, err
	}

	nftController.FilterNotExist()
	if !k.HasRestrictedNftsCollection(
		ctx,
		collectionController.getIndex(),
	) {
		nftController.FilterNotOwn(owner)
	}
	if err := nftController.Execute(); err != nil {
		return nil, err
	}

	nftsIds := nftController.getNftsIds()

	if len(nftsIds) == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidNftsCount, "not existing nfts or not an owner")
	}

	nftsIndexes := nftController.getIndexes()

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)
	if err := nftExecutor.BatchBurnNft(string(collectionIndex), nftsIds); err != nil {
		return nil, err
	}

	didExecutor := NewDidExecutor(ctx, "", k.dk)
	for _, nftIndex := range nftsIndexes {
		didExecutor.ForceDeleteDidOfNftIfExists(nftIndex)
	}

	k.DeleteApprovedNfts(ctx, collectionIndex, nftsIndexes)
	k.DeleteNfts(ctx, collectionIndex, nftsIndexes)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgBurnNfts),
			sdk.NewAttribute(types.AttributeKeyNftCollectionCreator, collectionCreator.String()),
			sdk.NewAttribute(types.AttributeKeyNftCollectionId, collectionId),
			sdk.NewAttribute(types.AttributeKeyNftsIds, strings.Join(nftsIds, ",")),
			sdk.NewAttribute(types.AttributeKeySigner, owner.String()),
		),
	)

	return &types.MsgBurnNftsResponse{
		NftsIds:           nftsIds,
		Burner:            owner.String(),
		CollectionCreator: collectionCreator.String(),
		CollectionId:      collectionId,
	}, nil
}

func (k msgServer) ApproveNfts(goCtx context.Context, msg *types.MsgApproveNfts) (*types.MsgApproveNftsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	receiver, err := sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return nil, err
	}

	var collectionCreator sdk.AccAddress

	if strings.TrimSpace(msg.CollectionCreator) == "" && !msg.Strict {
		msg.CollectionCreator = msg.Creator
		collectionCreator = owner
	} else {
		collectionCreator, err = sdk.AccAddressFromBech32(msg.CollectionCreator)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid collection creator")
		}
	}

	if err := k.CheckSoulBondedNftsCollection(ctx, collectionCreator.String(), msg.CollectionId); err != nil {
		return nil, sdkerrors.Wrap(err, "invalid operation on soul-bonded nfts collection")
	}

	if err := k.gk.CheckRestrictedNftsCollection(ctx, collectionCreator.String(), msg.CollectionId, msg.GetCreator()); err != nil {
		return nil, sdkerrors.Wrap(err, "unauthorized")
	}

	collectionController := NewNftCollectionController(ctx, collectionCreator, msg.Strict).
		WithId(msg.CollectionId).
		WithStore(k)

	if err := collectionController.MustExist().Validate(); err != nil {
		return nil, err
	}

	collectionIndex := collectionController.getIndex()
	collectionId := collectionController.getId()

	nftController := NewNftController(ctx, collectionIndex).
		WithIds(msg.Nfts.NftsIds).
		WithStore(k).
		WithConfiguration(k.GetParams(ctx))

	if err := nftController.ValidMetadataMaxCount().Validate(); err != nil {
		return nil, err
	}

	nftController.FilterNotExist()
	if !k.HasRestrictedNftsCollection(
		ctx,
		collectionController.getIndex(),
	) {
		nftController.FilterNotOwn(owner)
	}
	if err := nftController.Execute(); err != nil {
		return nil, err
	}

	nftsIds := nftController.getNftsIds()

	if len(nftsIds) == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidNftsCount, "not existing nfts or not an owner")
	}

	nftsIndexes := nftController.getIndexes()

	k.SetApprovedNfts(ctx, collectionIndex, nftsIndexes, owner, receiver, msg.Approved)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgApproveNfts),
			sdk.NewAttribute(types.AttributeKeyNftCollectionCreator, collectionCreator.String()),
			sdk.NewAttribute(types.AttributeKeyNftCollectionId, collectionId),
			sdk.NewAttribute(types.AttributeKeyNftsIds, strings.Join(nftsIds, ",")),
			sdk.NewAttribute(types.AttributeKeySigner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, receiver.String()),
		),
	)

	return &types.MsgApproveNftsResponse{
		NftsIds:           nftsIds,
		Owner:             owner.String(),
		Receiver:          receiver.String(),
		Approved:          msg.Approved,
		CollectionCreator: collectionCreator.String(),
		CollectionId:      collectionId,
	}, nil
}

func (k msgServer) ApproveAllNfts(goCtx context.Context, msg *types.MsgApproveAllNfts) (*types.MsgApproveAllNftsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	receiver, err := sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return nil, err
	}

	k.SetApprovedAllNfts(ctx, owner, receiver, msg.Approved)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgApproveAllNfts),
			sdk.NewAttribute(types.AttributeKeyApproved, strconv.FormatBool(msg.Approved)),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, receiver.String()),
		),
	)

	return &types.MsgApproveAllNftsResponse{
		Owner:    owner.String(),
		Receiver: receiver.String(),
		Approved: msg.Approved,
	}, nil
}

func (k msgServer) MintNft(goCtx context.Context, msg *types.MsgMintNft) (*types.MsgMintNftResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(msg.Receiver) == "" {
		msg.Receiver = msg.Creator
	}

	receiver, err := sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return nil, err
	}

	var collectionCreator sdk.AccAddress

	if strings.TrimSpace(msg.CollectionCreator) == "" && !msg.Strict {
		msg.CollectionCreator = msg.Creator
		collectionCreator = creator
	} else {
		collectionCreator, err = sdk.AccAddressFromBech32(msg.CollectionCreator)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid collection creator")
		}
	}

	if err := k.gk.CheckRestrictedNftsCollection(ctx, collectionCreator.String(), msg.CollectionId, msg.GetCreator()); err != nil {
		return nil, sdkerrors.Wrap(err, "unauthorized")
	}

	collectionController := NewNftCollectionController(ctx, collectionCreator, msg.Strict).
		WithId(msg.CollectionId).
		WithStore(k)

	if err := collectionController.CreateDefaultIfNotExists().Execute(); err != nil {
		return nil, err
	}

	collectionController.MustExist()

	restrictedCollection := k.HasRestrictedNftsCollection(
		ctx,
		collectionController.getIndex(),
	)

	if !restrictedCollection {
		collectionController.OpenedOrOwner(creator)
	}

	if msg.Did {
		soulBondNftsCollection := k.HasSoulBondedNftsCollection(
			ctx,
			collectionController.getIndex(),
		)

		if !restrictedCollection || !soulBondNftsCollection {
			return nil, sdkerrors.Wrap(types.ErrInvalidDid, "cannot use did for nft for collection which is not restricted and/or not for soul-bond nfts")
		}
	}

	if err := collectionController.Validate(); err != nil {
		return nil, err
	}

	collectionIndex := collectionController.getIndex()
	collectionId := collectionController.getId()

	nftController := NewNftController(ctx, collectionIndex).
		WithMetadata([]*types.MsgNftMetadata{msg.Nft}).
		WithStore(k).
		WithConfiguration(k.GetParams(ctx))

	if err := nftController.FilterExist().Execute(); err != nil {
		return nil, err
	}

	if err := nftController.ValidMetadata().Validate(); err != nil {
		return nil, err
	}

	nftsMetadata := nftController.getMetadata()

	if len(nftsMetadata) == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidNft, "existing or invalid nft")
	}

	nftMetadata := nftsMetadata[0]
	nftIndex := types.GetNftIndex(collectionIndex, nftMetadata.Id)

	newNft := nft.NFT{
		ClassId: string(collectionIndex),
		Id:      string(nftIndex),
		Uri:     types.ModuleName,
		UriHash: nftMetadata.Id,
		Data:    nftMetadata.Data,
	}

	var did string

	if msg.Did {
		didExecutor := NewDidExecutor(ctx, receiver.String(), k.dk)
		did, err = didExecutor.CreateDidForNft(nftIndex)
		if err != nil {
			return nil, err
		}
	}

	newNftMetadata := types.Nft{
		Index:             nftIndex,
		Id:                nftMetadata.Id,
		Images:            nftMetadata.Images,
		Url:               nftMetadata.Url,
		Links:             nftMetadata.Links,
		Title:             nftMetadata.Title,
		Description:       nftMetadata.Description,
		Attributes:        nftMetadata.Attributes,
		CollectionIndex:   collectionIndex,
		CollectionId:      collectionId,
		CollectionCreator: collectionCreator,
		Creator:           creator,
		Did:               did,
	}

	nftId := nftMetadata.Id

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)
	if err := nftExecutor.MintNft(newNft, receiver); err != nil {
		return nil, err
	}

	k.SetNft(ctx, newNftMetadata)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgMintNft),
			sdk.NewAttribute(types.AttributeKeyNftCollectionCreator, collectionCreator.String()),
			sdk.NewAttribute(types.AttributeKeyNftCollectionId, collectionId),
			sdk.NewAttribute(types.AttributeKeyNftId, nftId),
			sdk.NewAttribute(types.AttributeKeySigner, creator.String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, receiver.String()),
		),
	)

	return &types.MsgMintNftResponse{
		NftId:             nftId,
		Creator:           creator.String(),
		Receiver:          receiver.String(),
		CollectionCreator: collectionCreator.String(),
		CollectionId:      collectionId,
	}, nil
}

func (k msgServer) BurnNft(goCtx context.Context, msg *types.MsgBurnNft) (*types.MsgBurnNftResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	var collectionCreator sdk.AccAddress

	if strings.TrimSpace(msg.CollectionCreator) == "" && !msg.Strict {
		msg.CollectionCreator = msg.Creator
		collectionCreator = owner
	} else {
		collectionCreator, err = sdk.AccAddressFromBech32(msg.CollectionCreator)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid collection creator")
		}
	}

	if err := k.gk.CheckRestrictedNftsCollection(ctx, collectionCreator.String(), msg.CollectionId, msg.GetCreator()); err != nil {
		return nil, sdkerrors.Wrap(err, "unauthorized")
	}

	collectionController := NewNftCollectionController(ctx, collectionCreator, msg.Strict).
		WithId(msg.CollectionId).
		WithStore(k)

	if err := collectionController.MustExist().Validate(); err != nil {
		return nil, err
	}

	collectionIndex := collectionController.getIndex()
	collectionId := collectionController.getId()

	nftController := NewNftController(ctx, collectionIndex).
		WithId(msg.NftId).
		WithStore(k).
		WithConfiguration(k.GetParams(ctx))

	nftController.FilterNotExist()
	if !k.HasRestrictedNftsCollection(
		ctx,
		collectionController.getIndex(),
	) {
		nftController.FilterNotOwn(owner)
	}
	if err := nftController.Execute(); err != nil {
		return nil, err
	}

	nftsIds := nftController.getNftsIds()

	if len(nftsIds) == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidNft, "not existing nft or not an owner")
	}

	nftsIndexes := nftController.getIndexes()

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)
	if err := nftExecutor.BurnNft(string(collectionIndex), nftsIds[0]); err != nil {
		return nil, err
	}

	didExecutor := NewDidExecutor(ctx, "", k.dk)
	didExecutor.ForceDeleteDidOfNftIfExists(nftsIndexes[0])

	k.DeleteApprovedNft(ctx, collectionIndex, nftsIndexes[0])
	k.DeleteNft(ctx, collectionIndex, nftsIndexes[0])

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgBurnNft),
			sdk.NewAttribute(types.AttributeKeyNftCollectionCreator, collectionCreator.String()),
			sdk.NewAttribute(types.AttributeKeyNftCollectionId, collectionId),
			sdk.NewAttribute(types.AttributeKeyNftId, nftsIds[0]),
			sdk.NewAttribute(types.AttributeKeySigner, owner.String()),
		),
	)

	return &types.MsgBurnNftResponse{
		NftId:             nftsIds[0],
		Burner:            owner.String(),
		CollectionCreator: collectionCreator.String(),
		CollectionId:      collectionId,
	}, nil
}

func (k msgServer) ApproveNft(goCtx context.Context, msg *types.MsgApproveNft) (*types.MsgApproveNftResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	receiver, err := sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return nil, err
	}

	var collectionCreator sdk.AccAddress

	if strings.TrimSpace(msg.CollectionCreator) == "" && !msg.Strict {
		msg.CollectionCreator = msg.Creator
		collectionCreator = owner
	} else {
		collectionCreator, err = sdk.AccAddressFromBech32(msg.CollectionCreator)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid collection creator")
		}
	}

	if err := k.CheckSoulBondedNftsCollection(ctx, collectionCreator.String(), msg.CollectionId); err != nil {
		return nil, sdkerrors.Wrap(err, "invalid operation on soul-bonded nfts collection")
	}

	if err := k.gk.CheckRestrictedNftsCollection(ctx, collectionCreator.String(), msg.CollectionId, msg.GetCreator()); err != nil {
		return nil, sdkerrors.Wrap(err, "unauthorized")
	}

	collectionController := NewNftCollectionController(ctx, collectionCreator, msg.Strict).
		WithId(msg.CollectionId).
		WithStore(k)

	if err := collectionController.MustExist().Validate(); err != nil {
		return nil, err
	}

	collectionIndex := collectionController.getIndex()
	collectionId := collectionController.getId()

	nftController := NewNftController(ctx, collectionIndex).
		WithId(msg.NftId).
		WithStore(k).
		WithConfiguration(k.GetParams(ctx))

	nftController.FilterNotExist()
	if !k.HasRestrictedNftsCollection(
		ctx,
		collectionController.getIndex(),
	) {
		nftController.FilterNotOwn(owner)
	}
	if err := nftController.Execute(); err != nil {
		return nil, err
	}

	nftsIds := nftController.getNftsIds()

	if len(nftsIds) == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidNft, "not existing nft or not an owner")
	}

	nftsIndexes := nftController.getIndexes()

	k.SetApprovedNft(ctx, collectionIndex, nftsIndexes[0], owner, receiver, msg.Approved)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgApproveNft),
			sdk.NewAttribute(types.AttributeKeyNftCollectionCreator, collectionCreator.String()),
			sdk.NewAttribute(types.AttributeKeyNftCollectionId, collectionId),
			sdk.NewAttribute(types.AttributeKeyNftId, nftsIds[0]),
			sdk.NewAttribute(types.AttributeKeyApproved, strconv.FormatBool(msg.Approved)),
			sdk.NewAttribute(types.AttributeKeySigner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, receiver.String()),
		),
	)

	return &types.MsgApproveNftResponse{
		NftId:             nftsIds[0],
		Owner:             owner.String(),
		Receiver:          receiver.String(),
		Approved:          msg.Approved,
		CollectionCreator: collectionCreator.String(),
		CollectionId:      collectionId,
	}, nil
}

func (k msgServer) TransferNft(goCtx context.Context, msg *types.MsgTransferNft) (*types.MsgTransferNftResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	operator, err := sdk.AccAddressFromBech32(msg.Creator)
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

	var collectionCreator sdk.AccAddress

	if strings.TrimSpace(msg.CollectionCreator) == "" && !msg.Strict {
		msg.CollectionCreator = msg.Creator
		collectionCreator = operator
	} else {
		collectionCreator, err = sdk.AccAddressFromBech32(msg.CollectionCreator)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid collection creator")
		}
	}

	if err := k.CheckSoulBondedNftsCollection(ctx, collectionCreator.String(), msg.CollectionId); err != nil {
		return nil, sdkerrors.Wrap(err, "invalid operation on soul-bonded nfts collection")
	}

	if err := k.gk.CheckRestrictedNftsCollection(ctx, collectionCreator.String(), msg.CollectionId, msg.GetCreator()); err != nil {
		return nil, sdkerrors.Wrap(err, "unauthorized")
	}

	collectionController := NewNftCollectionController(ctx, collectionCreator, msg.Strict).
		WithId(msg.CollectionId).
		WithStore(k)

	if err := collectionController.MustExist().Validate(); err != nil {
		return nil, err
	}

	collectionIndex := collectionController.getIndex()
	collectionId := collectionController.getId()

	nftController := NewNftController(ctx, collectionIndex).
		WithId(msg.NftId).
		WithStore(k).
		WithConfiguration(k.GetParams(ctx))

	nftController.FilterNotExist()
	if !k.HasRestrictedNftsCollection(
		ctx,
		collectionController.getIndex(),
	) {
		nftController.FilterNotOwn(owner).FilterCannotTransfer(operator)
	}
	if err := nftController.Execute(); err != nil {
		return nil, err
	}

	nftsIds := nftController.getNftsIds()

	if len(nftsIds) == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidNft, "not existing nft or no transfer permission")
	}

	nftsIndexes := nftController.getIndexes()

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)
	if err := nftExecutor.TransferNft(string(collectionIndex), string(nftsIndexes[0]), receiver); err != nil {
		return nil, err
	}

	k.DeleteApprovedNft(ctx, collectionIndex, nftsIndexes[0])

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgTransferNft),
			sdk.NewAttribute(types.AttributeKeyNftCollectionCreator, collectionCreator.String()),
			sdk.NewAttribute(types.AttributeKeyNftCollectionId, collectionId),
			sdk.NewAttribute(types.AttributeKeyNftId, nftsIds[0]),
			sdk.NewAttribute(types.AttributeKeySigner, operator.String()),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, receiver.String()),
		),
	)

	return &types.MsgTransferNftResponse{
		NftId:             nftsIds[0],
		Operator:          operator.String(),
		Owner:             owner.String(),
		Receiver:          receiver.String(),
		CollectionCreator: collectionCreator.String(),
		CollectionId:      collectionId,
	}, nil
}

func (k msgServer) TransferNfts(goCtx context.Context, msg *types.MsgTransferNfts) (*types.MsgTransferNftsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	operator, err := sdk.AccAddressFromBech32(msg.Creator)
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

	var collectionCreator sdk.AccAddress

	if strings.TrimSpace(msg.CollectionCreator) == "" && !msg.Strict {
		msg.CollectionCreator = msg.Creator
		collectionCreator = operator
	} else {
		collectionCreator, err = sdk.AccAddressFromBech32(msg.CollectionCreator)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid collection creator")
		}
	}

	if err := k.CheckSoulBondedNftsCollection(ctx, collectionCreator.String(), msg.CollectionId); err != nil {
		return nil, sdkerrors.Wrap(err, "invalid operation on soul-bonded nfts collection")
	}

	if err := k.gk.CheckRestrictedNftsCollection(ctx, collectionCreator.String(), msg.CollectionId, msg.GetCreator()); err != nil {
		return nil, sdkerrors.Wrap(err, "unauthorized")
	}

	collectionController := NewNftCollectionController(ctx, collectionCreator, msg.Strict).
		WithId(msg.CollectionId).
		WithStore(k)

	if err := collectionController.MustExist().Validate(); err != nil {
		return nil, err
	}

	collectionIndex := collectionController.getIndex()
	collectionId := collectionController.getId()

	nftController := NewNftController(ctx, collectionIndex).
		WithIds(msg.Nfts.NftsIds).
		WithStore(k).
		WithConfiguration(k.GetParams(ctx))

	if err := nftController.ValidMetadataMaxCount().Validate(); err != nil {
		return nil, err
	}

	nftController.FilterNotExist()
	if !k.HasRestrictedNftsCollection(
		ctx,
		collectionController.getIndex(),
	) {
		nftController.FilterNotOwn(owner).FilterCannotTransfer(operator)
	}
	if err := nftController.Execute(); err != nil {
		return nil, err
	}

	nftsIds := nftController.getNftsIds()

	if len(nftsIds) == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidNftsCount, "not existing nfts or no transfer permission")
	}

	nftsIndexes := nftController.getIndexes()

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)
	if err := nftExecutor.BatchTransferNft(string(collectionIndex), nftsIds, receiver); err != nil {
		return nil, err
	}

	k.DeleteApprovedNfts(ctx, collectionIndex, nftsIndexes)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgTransferNfts),
			sdk.NewAttribute(types.AttributeKeyNftCollectionCreator, collectionCreator.String()),
			sdk.NewAttribute(types.AttributeKeyNftCollectionId, collectionId),
			sdk.NewAttribute(types.AttributeKeyNftsIds, strings.Join(nftsIds, ",")),
			sdk.NewAttribute(types.AttributeKeySigner, operator.String()),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, receiver.String()),
		),
	)

	return &types.MsgTransferNftsResponse{
		NftsIds:           nftsIds,
		Operator:          operator.String(),
		Owner:             owner.String(),
		Receiver:          receiver.String(),
		CollectionCreator: collectionCreator.String(),
		CollectionId:      collectionId,
	}, nil
}
