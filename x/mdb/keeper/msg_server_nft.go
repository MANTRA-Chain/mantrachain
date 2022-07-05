package keeper

import (
	"context"

	"github.com/LimeChain/mantrachain/x/mdb/types"
	// "github.com/LimeChain/mantrachain/x/mdb/utils"
	nfttypes "github.com/LimeChain/mantrachain/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k msgServer) MintNfts(goCtx context.Context, msg *types.MsgMintNfts) (*types.MsgMintNftsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Convert owner address string to sdk.AccAddress
	owner, err := sdk.AccAddressFromBech32(msg.Creator)

	if err != nil {
		return nil, err
	}

	if len(msg.Nfts.Nfts) == 0 {
		return nil, sdkerrors.Wrapf(types.ErrInvalidNftsCount, "nfts length %d invalid, min 1", len(msg.Nfts.Nfts))
	}

	var collectionCreator sdk.AccAddress

	if msg.CollectionCreator == "" {
		msg.CollectionCreator = msg.Creator
		collectionCreator = owner
	} else {
		collectionCreator, err = sdk.AccAddressFromBech32(msg.CollectionCreator)

		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid collection creator")
		}
	}

	collCtrl := NewNftCollectionController(ctx, &types.MsgCreateNftCollectionMetadata{
		Id: msg.CollectionId,
	}, collectionCreator).WithStore(k).WithConfiguration(k.GetParams(ctx))

	err = collCtrl.
		CreateDefaultIfNotExists().
		Execute()

	if err != nil {
		return nil, err
	}

	err = collCtrl.
		MustExist().
		CanMintNfts(owner).
		Validate()

	if err != nil {
		return nil, err
	}

	collIndex := collCtrl.getIndex()
	collId := collCtrl.getId()

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid collection creator")
	}

	ctrl := NewNftController(ctx, collIndex, msg.Nfts.Nfts).WithStore(k).WithConfiguration(k.GetParams(ctx))

	err = ctrl.
		FilterNotExist().
		Execute()

	if err != nil {
		return nil, err
	}

	err = ctrl.
		ValidNftMetadata().
		Validate()

	if err != nil {
		return nil, err
	}

	var nfts []nfttypes.NFT
	var ids []string
	filtered := ctrl.getFiltered()

	for _, nftMetadata := range filtered {
		index := types.GetNftIndex(collIndex, nftMetadata.Id)
		nfts = append(nfts, nfttypes.NFT{
			ClassId: string(collIndex),
			Id:      string(index),
			Uri:     types.ModuleName,
			UriHash: nftMetadata.Id,
			Data:    nftMetadata.Data,
		})
	}

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)
	_, err = nftExecutor.MintNftBatch(nfts, owner)
	if err != nil {
		return nil, err
	}

	// Uncomment after added delete did update did owner methods
	// didExecutor := NewDidExecutor(ctx, owner, msg.PubKeyHex, msg.PubKeyType, k.didKeeper)

	// TODO: use async iterator
	// TODO: use the upper range filtered
	for _, nftMetadata := range filtered {
		index := types.GetNftIndex(collIndex, nftMetadata.Id)
		// indexHex := utils.GetIndexHex(index)

		// _, err = didExecutor.SetDid(indexHex)
		// if err != nil {
		// 	return nil, err
		// }

		newNft := types.Nft{
			Index: index,
			// Did:             didExecutor.GetDidId(),
			Images:          nftMetadata.Images,
			Url:             nftMetadata.Url,
			Links:           nftMetadata.Links,
			Title:           nftMetadata.Title,
			Description:     nftMetadata.Description,
			Attributes:      nftMetadata.Attributes,
			Resellable:      nftMetadata.Resellable,
			CollectionIndex: collIndex,
			CollectionId:    collId,
			Creator:         owner,
		}

		k.SetNft(ctx, newNft)

		ids = append(ids, string(nftMetadata.Id))
	}

	// TODO: emit event

	return &types.MsgMintNftsResponse{
		Ids: ids,
	}, nil
}

func (k msgServer) BurnNfts(goCtx context.Context, msg *types.MsgBurnNfts) (*types.MsgBurnNftsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Convert owner address string to sdk.AccAddress
	owner, err := sdk.AccAddressFromBech32(msg.Creator)

	if err != nil {
		return nil, err
	}

	if len(msg.Nfts.NftsIds) == 0 {
		return nil, sdkerrors.Wrapf(types.ErrInvalidNftsBurnCount, "nfts length %d invalid, min 1", len(msg.Nfts.NftsIds))
	}

	var collectionCreator sdk.AccAddress

	if msg.CollectionCreator == "" {
		msg.CollectionCreator = msg.Creator
		collectionCreator = owner
	} else {
		collectionCreator, err = sdk.AccAddressFromBech32(msg.CollectionCreator)

		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid collection creator")
		}
	}

	collCtrl := NewNftCollectionController(ctx, &types.MsgCreateNftCollectionMetadata{
		Id: msg.CollectionId,
	}, collectionCreator).WithStore(k).WithConfiguration(k.GetParams(ctx))

	err = collCtrl.
		MustExist().
		CanBurnNfts(owner).
		Validate()

	if err != nil {
		return nil, err
	}

	collIndex := collCtrl.getIndex()

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid collection creator")
	}

	var nfts []*types.MsgNftMetadata

	for _, id := range msg.Nfts.NftsIds {
		nfts = append(nfts, &types.MsgNftMetadata{
			Id: id,
		})
	}

	ctrl := NewNftController(ctx, collIndex, nfts).WithStore(k).WithConfiguration(k.GetParams(ctx))

	err = ctrl.
		FilterNotExist().
		FilterNotOwn(owner).
		Execute()

	if err != nil {
		return nil, err
	}

	var nftsIds []string
	var ids []string
	filtered := ctrl.getFiltered()

	for _, nftMetadata := range filtered {
		index := types.GetNftIndex(collIndex, nftMetadata.Id)
		nftsIds = append(nftsIds, string(index))
	}

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)
	_, err = nftExecutor.BurnNftBatch(string(collIndex), nftsIds)
	if err != nil {
		return nil, err
	}

	// TODO: use async iterator
	for _, id := range nftsIds {
		index := types.GetNftIndex(collIndex, id)

		// TODO: Add delete did method

		k.DeleteNft(ctx, collIndex, index)

		ids = append(ids, id)
	}

	// TODO: emit event

	return &types.MsgBurnNftsResponse{
		Ids: ids,
	}, nil
}
