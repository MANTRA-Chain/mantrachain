package keeper

import (
	"context"

	"github.com/LimeChain/mantrachain/x/mdb/types"
	nfttypes "github.com/LimeChain/mantrachain/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k msgServer) MintNft(goCtx context.Context, msg *types.MsgMintNft) (*types.MsgMintNftResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Convert owner address string to sdk.AccAddress
	owner, err := sdk.AccAddressFromBech32(msg.Creator)

	if err != nil {
		return nil, err
	}

	if len(msg.Nfts.Metadata) == 0 {
		return nil, sdkerrors.Wrapf(types.ErrInvalidNftsLength, "nfts length %d invalid, min 1", len(msg.Nfts.Metadata))
	}

	collectionCreator, err := sdk.AccAddressFromBech32(msg.Nfts.CollectionCreator)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid collection creator")
	}

	collCtrl := NewNftCollectionController(ctx, &types.MsgCreateNftCollectionMetadata{
		Id: msg.Nfts.CollectionId,
	}, collectionCreator).WithStore(k).WithConfiguration(k.GetParams(ctx))

	err = collCtrl.
		CreateDefaultIfNotExists().
		Execute()

	collIndex := collCtrl.getIndex()

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid collection creator")
	}

	ctrl := NewNftController(ctx, collIndex, msg.Nfts.Metadata).WithStore(k).WithConfiguration(k.GetParams(ctx))

	err = ctrl.
		FilterNotExist().
		ValidNftMetadata().
		Validate()

	if err != nil {
		return nil, err
	}

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)

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

	_, err = nftExecutor.MintNftBatch(nfts, owner)
	if err != nil {
		return nil, err
	}

	for _, nftMetadata := range filtered {
		index := types.GetNftIndex(collIndex, nftMetadata.Id)
		newNft := types.Nft{
			Index:           index,
			Images:          nftMetadata.Images,
			Url:             nftMetadata.Url,
			Links:           nftMetadata.Links,
			Title:           nftMetadata.Title,
			Description:     nftMetadata.Description,
			Attributes:      nftMetadata.Attributes,
			Resellable:      nftMetadata.Resellable,
			CollectionIndex: collIndex,
			Owner:           owner,
			Creator:         owner,
		}

		k.SetNft(ctx, newNft)

		ids = append(ids, string(nftMetadata.Id))
	}

	return &types.MsgMintNftResponse{
		Ids: ids,
	}, nil
}
