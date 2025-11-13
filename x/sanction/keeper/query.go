package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"github.com/MANTRA-Chain/mantrachain/v7/x/sanction/types"
	"github.com/cosmos/cosmos-sdk/types/query"
)

var _ types.QueryServer = queryServer{}

// NewQueryServerImpl returns an implementation of the QueryServer interface
// for the provided Keeper.
func NewQueryServerImpl(k Keeper) types.QueryServer {
	return queryServer{k}
}

type queryServer struct {
	k Keeper
}

func (q queryServer) Blacklist(
	ctx context.Context,
	req *types.QueryBlacklistRequest,
) (*types.QueryBlacklistResponse, error) {
	blacklist := []string{}
	_, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.BlacklistAccounts,
		req.Pagination,
		func(key string, _ collections.NoValue) (bool, error) {
			blacklist = append(blacklist, key)
			return true, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return &types.QueryBlacklistResponse{
		BlacklistedAccounts: blacklist,
		Pagination:          pageRes,
	}, nil
}
