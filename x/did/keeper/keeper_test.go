package keeper

import (
	"fmt"
	"testing"

	"cosmossdk.io/log"
	storemetrics "cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	didtestutil "github.com/MANTRA-Finance/mantrachain/x/did/testutil"
	"github.com/MANTRA-Finance/mantrachain/x/did/types"
	tmdb "github.com/cosmos/cosmos-db"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"cosmossdk.io/store"
	cbproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
)

// Keeper test suit enables the keeper package to be tested
type KeeperTestSuite struct {
	suite.Suite

	ctx         sdk.Context
	keeper      Keeper
	queryClient types.QueryClient
}

// SetupTest creates a test suite to test the did
func (suite *KeeperTestSuite) SetupTest() {
	storeKey := storetypes.NewKVStoreKey(types.ModuleName)

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	ctrl := gomock.NewController(suite.T())
	logger := log.NewNopLogger()

	storeService := runtime.NewKVStoreService(storeKey)
	authority := authtypes.NewModuleAddress(types.ModuleName).String()

	guardKeeper := didtestutil.NewMockGuardKeeper(ctrl)

	k := NewKeeper(
		cdc,
		storeService,
		logger,
		authority,
		guardKeeper,
	)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, logger, storemetrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(suite.T(), stateStore.LoadLatestVersion())

	ctx := sdk.NewContext(stateStore, cbproto.Header{ChainID: "foochainid"}, true, logger)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, cdc.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, k)
	queryClient := types.NewQueryClient(queryHelper)

	suite.ctx, suite.keeper, suite.queryClient = ctx, k, queryClient

	guardKeeper.EXPECT().CheckIsAdmin(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestGenericKeeperSetAndGet() {
	testCases := []struct {
		msg     string
		didFn   func() types.DidDocument
		expPass bool
	}{
		{
			"data stored successfully",
			func() types.DidDocument {
				dd, _ := types.NewDidDocument(
					"did:mantrachain:subject",
				)
				return dd
			},
			true,
		},
	}
	for _, tc := range testCases {
		dd := tc.didFn()
		suite.keeper.Set(suite.ctx,
			[]byte(dd.Id),
			[]byte{0x01},
			dd,
			suite.keeper.Marshal,
		)
		suite.keeper.Set(suite.ctx,
			[]byte(dd.Id+"1"),
			[]byte{0x01},
			dd,
			suite.keeper.Marshal,
		)
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			if tc.expPass {
				_, found := suite.keeper.Get(
					suite.ctx,
					[]byte(dd.Id),
					[]byte{0x01},
					suite.keeper.UnmarshalDidDocument,
				)
				suite.Require().True(found)

				iterator := suite.keeper.GetAll(
					suite.ctx,
					[]byte{0x01},
				)
				defer iterator.Close()

				var array []interface{}
				for ; iterator.Valid(); iterator.Next() {
					array = append(array, iterator.Value())
				}
				suite.Require().Equal(2, len(array))
			} else {
				// TODO write failure cases
				suite.Require().False(tc.expPass)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGenericKeeperDelete() {
	testCases := []struct {
		msg     string
		didFn   func() types.DidDocument
		expPass bool
	}{
		{
			"data stored successfully",
			func() types.DidDocument {
				dd, _ := types.NewDidDocument(
					"did:mantrachain:subject",
				)
				return dd
			},
			true,
		},
	}
	for _, tc := range testCases {
		dd := tc.didFn()
		suite.keeper.Set(suite.ctx,
			[]byte(dd.Id),
			[]byte{0x01},
			dd,
			suite.keeper.Marshal,
		)
		suite.keeper.Set(suite.ctx,
			[]byte(dd.Id+"1"),
			[]byte{0x01},
			dd,
			suite.keeper.Marshal,
		)
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			if tc.expPass {
				suite.keeper.Delete(
					suite.ctx,
					[]byte(dd.Id),
					[]byte{0x01},
				)

				_, found := suite.keeper.Get(
					suite.ctx,
					[]byte(dd.Id),
					[]byte{0x01},
					suite.keeper.UnmarshalDidDocument,
				)
				suite.Require().False(found)

			} else {
				// TODO write failure cases
				suite.Require().False(tc.expPass)
			}
		})
	}
}
