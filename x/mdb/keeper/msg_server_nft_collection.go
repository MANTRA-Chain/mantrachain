package keeper

import (
	"context"

	"github.com/LimeChain/mantrachain/x/mdb/types"
	"github.com/LimeChain/mantrachain/x/mdb/utils"
	nfttypes "github.com/LimeChain/mantrachain/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateNftCollection(goCtx context.Context, msg *types.MsgCreateNftCollection) (*types.MsgCreateNftCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Convert owner address string to sdk.AccAddress
	owner, err := sdk.AccAddressFromBech32(msg.Creator)

	if err != nil {
		return nil, err
	}

	ctrl := NewNftCollectionController(ctx, msg.Collection, owner).
		WithStore(k).
		WithConfiguration(k.GetParams(ctx))

	err = ctrl.
		MustNotExist().
		MustNotBeDefault().
		ValidNftCollectionMetadata().
		Validate()

	if err != nil {
		return nil, err
	}

	index := ctrl.getIndex()
	id := ctrl.getId()
	indexHex := utils.GetIndexHex(index)

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)
	_, err = nftExecutor.SetClass(nfttypes.Class{
		Id:          string(index),
		Name:        msg.Collection.Name,
		Symbol:      msg.Collection.Symbol,
		Description: msg.Collection.Description,
		Uri:         types.ModuleName,
		UriHash:     id,
		Data:        msg.Collection.Data,
	})
	if err != nil {
		return nil, err
	}

	didExecutor := NewDidExecutor(ctx, owner, msg.PubKeyHex, msg.PubKeyType, k.didKeeper)
	_, err = didExecutor.SetDid(indexHex)
	if err != nil {
		return nil, err
	}

	newNftCollection := types.NftCollection{
		Index:    index,
		Did:      didExecutor.GetDidId(),
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

	return &types.MsgCreateNftCollectionResponse{
		Id: id,
	}, nil
}
