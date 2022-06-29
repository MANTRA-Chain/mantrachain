package keeper

import (
	"context"

	"github.com/LimeChain/mantrachain/x/mdb/types"
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

	ctrl := NewNftCollectionController(ctx, msg.Metadata, owner).WithStore(k).WithConfiguration(k.GetParams(ctx))

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

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)
	_, err = nftExecutor.SetClass(nfttypes.Class{
		Id:          string(index),
		Name:        msg.Metadata.Name,
		Symbol:      msg.Metadata.Symbol,
		Description: msg.Metadata.Description,
		Uri:         "mdb",
		UriHash:     id,
		Data:        msg.Metadata.Data,
	})
	if err != nil {
		return nil, err
	}

	newNftCollection := types.NftCollection{
		Index:           index,
		Images:          msg.Metadata.Images,
		Url:             msg.Metadata.Url,
		Links:           msg.Metadata.Links,
		Category:        msg.Metadata.Category,
		CreatorEarnings: msg.Metadata.CreatorEarnings,
		DisplayTheme:    msg.Metadata.DisplayTheme,
		Opened:          msg.Metadata.Opened,
		Creator:         owner,
		Owner:           owner,
	}

	k.SetNftCollection(ctx, newNftCollection)

	return &types.MsgCreateNftCollectionResponse{
		Id: id,
	}, nil
}
