package app

import (
	"context"

	"cosmossdk.io/x/feegrant"
	feegrantkeeper "cosmossdk.io/x/feegrant/keeper"
)

type sanctionFeegrantKeeperAdapter struct {
	query feegrant.QueryServer
	msg   feegrant.MsgServer
}

func newSanctionFeegrantKeeperAdapter(k feegrantkeeper.Keeper) *sanctionFeegrantKeeperAdapter {
	return &sanctionFeegrantKeeperAdapter{
		query: k,
		msg:   feegrantkeeper.NewMsgServerImpl(k),
	}
}

func (a *sanctionFeegrantKeeperAdapter) AllowancesByGranter(ctx context.Context, req *feegrant.QueryAllowancesByGranterRequest) (*feegrant.QueryAllowancesByGranterResponse, error) {
	return a.query.AllowancesByGranter(ctx, req)
}

func (a *sanctionFeegrantKeeperAdapter) RevokeAllowance(ctx context.Context, msg *feegrant.MsgRevokeAllowance) (*feegrant.MsgRevokeAllowanceResponse, error) {
	return a.msg.RevokeAllowance(ctx, msg)
}
