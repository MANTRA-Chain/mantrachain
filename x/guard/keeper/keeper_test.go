package keeper_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/libs/log"
	tmtproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmttime "github.com/tendermint/tendermint/types/time"

	"github.com/MANTRA-Finance/mantrachain/x/guard"
	"github.com/MANTRA-Finance/mantrachain/x/guard/keeper"
	guardtestutil "github.com/MANTRA-Finance/mantrachain/x/guard/testutil"
	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/MANTRA-Finance/mantrachain/testutil"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type KeeperTestSuite struct {
	testutil.IBCConnectionTestSuite

	ctx               sdk.Context
	addrs             []sdk.AccAddress
	baseApp           *baseapp.BaseApp
	guardKeeper       keeper.Keeper
	encCfg            testutil.TestEncodingConfig
	queryClient       types.QueryClient
	msgServer         types.MsgServer
	defaultPrivileges []byte
	kind              types.RequiredPrivilegesKind
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	key := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := sdk.NewMemoryStoreKeys(types.StoreKey)
	moduleAccAddr := make(map[string]bool)
	testCtx := testutil.DefaultContextWithDB(s.T(), key, sdk.NewTransientStoreKey("transient_test"))
	s.ctx = testCtx.Ctx.WithBlockHeader(tmtproto.Header{Time: tmttime.Now()})
	s.encCfg = testutil.MakeTestEncodingConfig(guard.AppModuleBasic{})
	paramSpace := paramskeeper.NewSubspace(s.encCfg.Codec, s.encCfg.Amino, key, memStoreKey[types.MemStoreKey], types.ModuleName)

	s.baseApp = baseapp.NewBaseApp(
		types.ModuleName,
		log.NewNopLogger(),
		testCtx.DB,
		s.encCfg.TxConfig.TxDecoder(),
	)
	s.baseApp.SetCMS(testCtx.CMS)
	s.baseApp.SetInterfaceRegistry(s.encCfg.InterfaceRegistry)

	s.addrs = testutil.CreateIncrementalAccounts(3)

	ctrl := gomock.NewController(s.T())
	accountKeeper := guardtestutil.NewMockAccountKeeper(ctrl)
	bankKeeper := guardtestutil.NewMockBankKeeper(ctrl)
	authzKeeper := guardtestutil.NewMockAuthzKeeper(ctrl)
	tokenKeeper := guardtestutil.NewMockTokenKeeper(ctrl)
	nftKeeper := guardtestutil.NewMockNFTKeeper(ctrl)
	coinFactoryKeeper := guardtestutil.NewMockCoinFactoryKeeper(ctrl)

	s.guardKeeper = keeper.NewKeeper(
		s.encCfg.Codec,
		key,
		memStoreKey[types.MemStoreKey],
		paramSpace,
		moduleAccAddr,
		s.baseApp.MsgServiceRouter(),
		accountKeeper,
		bankKeeper,
		authzKeeper,
		tokenKeeper,
		nftKeeper,
		coinFactoryKeeper,
	)

	queryHelper := baseapp.NewQueryServerTestHelper(s.ctx, s.encCfg.InterfaceRegistry)
	types.RegisterQueryServer(queryHelper, s.guardKeeper)
	queryClient := types.NewQueryClient(queryHelper)
	s.queryClient = queryClient

	s.defaultPrivileges = []byte{0x01}
	s.kind = types.RequiredPrivilegesCoin

	s.msgServer = keeper.NewMsgServerImpl(&s.guardKeeper)
}
