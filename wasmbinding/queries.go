package wasmbinding

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/MANTRA-Finance/aumega/wasmbinding/bindings"
	coinfactorykeeper "github.com/MANTRA-Finance/aumega/x/coinfactory/keeper"
)

type QueryPlugin struct {
	coinFactoryKeeper *coinfactorykeeper.Keeper
}

// NewQueryPlugin returns a reference to a new QueryPlugin.
func NewQueryPlugin(tfk *coinfactorykeeper.Keeper) *QueryPlugin {
	return &QueryPlugin{
		coinFactoryKeeper: tfk,
	}
}

// GetDenomAdmin is a query to get denom admin.
func (qp QueryPlugin) GetDenomAdmin(ctx sdk.Context, denom string) (*bindings.DenomAdminResponse, error) {
	metadata, err := qp.coinFactoryKeeper.GetAuthorityMetadata(ctx, denom)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin for denom: %s", denom)
	}

	return &bindings.DenomAdminResponse{Admin: metadata.Admin}, nil
}
