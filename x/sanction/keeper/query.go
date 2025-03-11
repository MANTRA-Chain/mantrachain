package keeper

import (
	"context"

	"github.com/MANTRA-Chain/mantrachain/v4/x/sanction/types"
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
	iter, err := q.k.BlacklistAccounts.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	blacklist, err := iter.Keys()
	if err != nil {
		return nil, err
	}

	return &types.QueryBlacklistResponse{BlacklistedAccounts: blacklist}, nil
}
