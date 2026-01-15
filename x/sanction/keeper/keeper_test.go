package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"cosmossdk.io/core/header"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"

	sanctionkeeper "github.com/MANTRA-Chain/mantrachain/v8/x/sanction/keeper"
	"github.com/MANTRA-Chain/mantrachain/v8/x/sanction/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx            sdk.Context
	sanctionKeeper sanctionkeeper.Keeper

	msgServer   types.MsgServer
	queryServer types.QueryServer
}

func (suite *KeeperTestSuite) SetupTest() {
	key := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(key)
	testCtx := testutil.DefaultContextWithDB(suite.T(), key, storetypes.NewTransientStoreKey("transient_test"))
	ctx := testCtx.Ctx.WithHeaderInfo(header.Info{Time: time.Now()})
	encCfg := moduletestutil.MakeTestEncodingConfig()

	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()
	sanctionKeeper := sanctionkeeper.NewKeeper(
		encCfg.Codec,
		storeService,
		log.NewNopLogger(),
		authority,
	)
	suite.ctx = ctx
	suite.sanctionKeeper = sanctionKeeper

	err := suite.sanctionKeeper.Params.Set(ctx, types.Params{})
	suite.Require().NoError(err)

	types.RegisterInterfaces(encCfg.InterfaceRegistry)
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, encCfg.InterfaceRegistry)
	queryServer := sanctionkeeper.NewQueryServerImpl(sanctionKeeper)
	types.RegisterQueryServer(queryHelper, queryServer)
	suite.msgServer = sanctionkeeper.NewMsgServerImpl(sanctionKeeper)
	suite.queryServer = queryServer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}