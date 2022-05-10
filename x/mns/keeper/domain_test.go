package keeper_test

import (
	"strconv"
	"testing"

	"github.com/LimeChain/mantrachain/x/mns/keeper"
	"github.com/LimeChain/mantrachain/x/mns/types"
	keepertest "github.com/LimeChain/mantrachain/testutil/keeper"
	"github.com/LimeChain/mantrachain/testutil/nullify"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNDomain(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Domain {
	items := make([]types.Domain, n)
	for i := range items {
		items[i].Index = strconv.Itoa(i)
        
		keeper.SetDomain(ctx, items[i])
	}
	return items
}

func TestDomainGet(t *testing.T) {
	keeper, ctx := keepertest.MnsKeeper(t)
	items := createNDomain(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetDomain(ctx,
		    item.Index,
            
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestDomainRemove(t *testing.T) {
	keeper, ctx := keepertest.MnsKeeper(t)
	items := createNDomain(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveDomain(ctx,
		    item.Index,
            
		)
		_, found := keeper.GetDomain(ctx,
		    item.Index,
            
		)
		require.False(t, found)
	}
}

func TestDomainGetAll(t *testing.T) {
	keeper, ctx := keepertest.MnsKeeper(t)
	items := createNDomain(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllDomain(ctx)),
	)
}
