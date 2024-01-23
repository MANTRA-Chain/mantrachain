package keeper

import (
	"context"

	"github.com/MANTRA-Finance/mantrachain/x/rewards/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ProviderAll(goCtx context.Context, req *types.QueryAllProviderRequest) (*types.QueryAllProviderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var providers []types.Provider
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	providerStore := prefix.NewStore(store, types.KeyPrefix(types.ProviderKeyPrefix))

	pageRes, err := query.Paginate(providerStore, req.Pagination, func(key []byte, value []byte) error {
		var provider types.Provider
		if err := k.cdc.Unmarshal(value, &provider); err != nil {
			return err
		}

		providers = append(providers, provider)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllProviderResponse{Provider: providers, Pagination: pageRes}, nil
}

func (k Keeper) ProviderPairs(goCtx context.Context, req *types.QueryGetProviderPairsRequest) (*types.QueryGetProviderPairsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator, err := sdk.AccAddressFromBech32(req.Provider)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	provider, found := k.GetProvider(ctx, creator.String())

	if !found {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var pairsIds []uint64

	for pairId := range provider.PairIdToIdx {
		pairsIds = append(pairsIds, pairId)
	}

	return &types.QueryGetProviderPairsResponse{Provider: creator.String(), PairsIds: pairsIds}, nil
}

func (k Keeper) Provider(goCtx context.Context, req *types.QueryGetProviderRequest) (*types.QueryGetProviderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator, err := sdk.AccAddressFromBech32(req.Provider)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	provider, found := k.GetProvider(ctx, creator.String())

	if !found {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	return &types.QueryGetProviderResponse{Provider: provider}, nil
}
