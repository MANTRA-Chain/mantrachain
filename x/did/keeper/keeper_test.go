package keeper

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"

	didtestutil "github.com/MANTRA-Finance/mantrachain/x/did/testutil"
	"github.com/MANTRA-Finance/mantrachain/x/did/types"
	"github.com/cosmos/cosmos-sdk/codec"
	ct "github.com/cosmos/cosmos-sdk/codec/types"
	server "github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	dbm "github.com/tendermint/tm-db"
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
	keyDidDocument := sdk.NewKVStoreKey(types.StoreKey)
	memKeyDidDocument := sdk.NewKVStoreKey(types.MemStoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyDidDocument, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(memKeyDidDocument, sdk.StoreTypeIAVL, db)
	_ = ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, tmproto.Header{ChainID: "foochainid"}, true, server.ZeroLogWrapper{log.Logger})

	interfaceRegistry := ct.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	ctrl := gomock.NewController(suite.T())
	guardKeeper := didtestutil.NewMockGuardKeeper(ctrl)

	k := NewKeeper(
		marshaler,
		keyDidDocument,
		memKeyDidDocument,
		guardKeeper,
	)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, interfaceRegistry)
	types.RegisterQueryServer(queryHelper, k)
	queryClient := types.NewQueryClient(queryHelper)

	suite.ctx, suite.keeper, suite.queryClient = ctx, *k, queryClient

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
