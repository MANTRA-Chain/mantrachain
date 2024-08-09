package keeper_test

import (
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	storemetrics "cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	"github.com/MANTRA-Finance/mantrachain/testutil"
	"github.com/MANTRA-Finance/mantrachain/x/txfees/keeper"
	txfeestestutil "github.com/MANTRA-Finance/mantrachain/x/txfees/testutil"
	"github.com/MANTRA-Finance/mantrachain/x/txfees/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmdb "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TxfeesKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey(types.ModuleName)

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	ctrl := gomock.NewController(t)
	logger := log.NewNopLogger()

	storeService := runtime.NewKVStoreService(storeKey)
	authority := authtypes.NewModuleAddress(types.ModuleName).String()

	accountKeeper := txfeestestutil.NewMockAccountKeeper(ctrl)
	liquidityKeeper := txfeestestutil.NewMockLiquidityKeeper(ctrl)
	guardKeeper := txfeestestutil.NewMockGuardKeeper(ctrl)

	k := keeper.NewKeeper(
		cdc,
		storeService,
		logger,
		authority,
		accountKeeper,
		liquidityKeeper,
		guardKeeper,
	)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, logger, storemetrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	guardKeeper.EXPECT().CheckIsAdmin(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	guardKeeper.EXPECT().GetAdmin(gomock.Any()).Return(sdk.MustAccAddressFromBech32(testutil.TestAdminAddress)).AnyTimes()

	return &k, ctx
}