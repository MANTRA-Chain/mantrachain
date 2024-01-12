package keeper_test

import (
	"testing"

	"github.com/AumegaChain/aumega/testutil"
	"github.com/AumegaChain/aumega/x/txfees/keeper"
	txfeestestutil "github.com/AumegaChain/aumega/x/txfees/testutil"
	"github.com/AumegaChain/aumega/x/txfees/types"
	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TxfeesKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	paramsSubspace := typesparams.NewSubspace(cdc,
		types.Amino,
		storeKey,
		memStoreKey,
		"TxfeesParams",
	)

	ctrl := gomock.NewController(t)
	accountKeeper := txfeestestutil.NewMockAccountKeeper(ctrl)
	bankKeeper := txfeestestutil.NewMockBankKeeper(ctrl)
	guardKeeper := txfeestestutil.NewMockGuardKeeper(ctrl)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		accountKeeper,
		bankKeeper,
		guardKeeper,
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	guardKeeper.EXPECT().CheckIsAdmin(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	guardKeeper.EXPECT().GetAdmin(gomock.Any()).Return(sdk.MustAccAddressFromBech32(testutil.TestAdminAddress)).AnyTimes()

	return k, ctx
}
