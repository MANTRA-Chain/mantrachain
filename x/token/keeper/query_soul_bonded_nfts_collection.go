package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/LimeChain/mantrachain/x/token/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) SoulBondedNftsCollectionAll(goCtx context.Context, req *types.QueryAllSoulBondedNftsCollectionRequest) (*types.QueryAllSoulBondedNftsCollectionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var soulBondedNftsCollections []types.SoulBondedNftsCollection
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	soulBondedNftsCollectionStore := prefix.NewStore(store, types.KeyPrefix(types.SoulBondedNftsCollectionKeyPrefix))

	pageRes, err := query.Paginate(soulBondedNftsCollectionStore, req.Pagination, func(key []byte, value []byte) error {
		var soulBondedNftsCollection types.SoulBondedNftsCollection
		if err := k.cdc.Unmarshal(value, &soulBondedNftsCollection); err != nil {
			return err
		}

		soulBondedNftsCollections = append(soulBondedNftsCollections, soulBondedNftsCollection)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllSoulBondedNftsCollectionResponse{SoulBondedNftsCollection: soulBondedNftsCollections, Pagination: pageRes}, nil
}

func (k Keeper) SoulBondedNftsCollection(goCtx context.Context, req *types.QueryGetSoulBondedNftsCollectionRequest) (*types.QueryGetSoulBondedNftsCollectionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetSoulBondedNftsCollection(
	    ctx,
	    req.Index,
        )
	if !found {
	    return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetSoulBondedNftsCollectionResponse{SoulBondedNftsCollection: val}, nil
}