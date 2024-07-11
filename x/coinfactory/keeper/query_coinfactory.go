package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/MANTRA-Finance/mantrachain/x/coinfactory/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) DenomAuthorityMetadata(ctx context.Context, req *types.QueryDenomAuthorityMetadataRequest) (*types.QueryDenomAuthorityMetadataResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	authorityMetadata, err := k.GetAuthorityMetadata(sdkCtx, req.GetDenom())
	if err != nil {
		return nil, err
	}

	return &types.QueryDenomAuthorityMetadataResponse{AuthorityMetadata: authorityMetadata}, nil
}

func (k Keeper) DenomAuthorityMetadata2(ctx context.Context, req *types.QueryDenomAuthorityMetadata2Request) (*types.QueryDenomAuthorityMetadata2Response, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	denom, err := types.GetTokenDenom(req.GetCreator(), req.GetSubdenom())
	if err != nil {
		return nil, err
	}

	authorityMetadata, err := k.GetAuthorityMetadata(sdkCtx, denom)
	if err != nil {
		return nil, err
	}

	return &types.QueryDenomAuthorityMetadata2Response{AuthorityMetadata: authorityMetadata}, nil
}

func (k Keeper) DenomsFromCreator(ctx context.Context, req *types.QueryDenomsFromCreatorRequest) (*types.QueryDenomsFromCreatorResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	denoms := k.getDenomsFromCreator(sdkCtx, req.GetCreator())
	return &types.QueryDenomsFromCreatorResponse{Denoms: denoms}, nil
}

func (k Keeper) Balance(ctx context.Context, req *types.QueryBalanceRequest) (*types.QueryBalanceResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	denom, err := types.GetTokenDenom(req.GetCreator(), req.GetSubdenom())
	if err != nil {
		return nil, err
	}

	address, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	balance := k.bankKeeper.GetBalance(sdkCtx, address, denom)

	return &types.QueryBalanceResponse{Balance: &balance}, nil
}

func (k Keeper) SupplyOf(ctx context.Context, req *types.QuerySupplyOfRequest) (*types.QuerySupplyOfResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	denom, err := types.GetTokenDenom(req.GetCreator(), req.GetSubdenom())
	if err != nil {
		return nil, err
	}

	amount := k.bankKeeper.GetSupply(sdkCtx, denom)

	return &types.QuerySupplyOfResponse{Amount: amount}, nil
}

func (k Keeper) DenomMetadata(ctx context.Context, req *types.QueryDenomMetadataRequest) (*types.QueryDenomMetadataResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	denom, err := types.GetTokenDenom(req.GetCreator(), req.GetSubdenom())
	if err != nil {
		return nil, err
	}

	metadata, found := k.bankKeeper.GetDenomMetaData(sdkCtx, denom)
	if !found {
		return nil, types.ErrDenomDoesNotExist
	}

	return &types.QueryDenomMetadataResponse{Metadata: metadata}, nil
}
