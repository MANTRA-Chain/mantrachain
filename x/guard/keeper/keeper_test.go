package keeper_test

import (
	"testing"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/MANTRA-Finance/mantrachain/testutil"
	"github.com/MANTRA-Finance/mantrachain/x/guard/keeper"
	guardtestutil "github.com/MANTRA-Finance/mantrachain/x/guard/testutil"
	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
	cbproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cbttime "github.com/cometbft/cometbft/types/time"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	testAccount    = "mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka"
	testIndex      = []byte{0x01, 0x02, 0x02}
	testPrivileges = []byte{0x04, 0x05, 0x06}
)

type KeeperTestSuite struct {
	suite.Suite

	ctx               sdk.Context
	addrs             []sdk.AccAddress
	guardKeeper       keeper.Keeper
	bankKeeper        *guardtestutil.MockBankKeeper
	nftKeeper         *guardtestutil.MockNFTKeeper
	coinFactoryKeeper *guardtestutil.MockCoinFactoryKeeper
	queryClient       types.QueryClient
	msgServer         types.MsgServer
	defaultPrivileges []byte
	rpKind            types.RequiredPrivilegesKind
	lkIndex           []byte
	params            types.Params
	testAccount       string
	testAdminAccount  string
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	transientStorekey := storetypes.NewTransientStoreKey("transient_test")

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	ctrl := gomock.NewController(s.T())
	logger := log.NewNopLogger()

	storeService := runtime.NewKVStoreService(storeKey)
	authority := authtypes.NewModuleAddress(types.ModuleName).String()

	testCtx := testutil.DefaultContextWithDB(s.T(), storeKey, transientStorekey)
	s.ctx = testCtx.Ctx.WithBlockHeader(cbproto.Header{Time: cbttime.Now()})

	s.addrs = testutil.CreateIncrementalAccounts(3)

	s.bankKeeper = guardtestutil.NewMockBankKeeper(ctrl)
	authzKeeper := guardtestutil.NewMockAuthzKeeper(ctrl)
	s.nftKeeper = guardtestutil.NewMockNFTKeeper(ctrl)
	s.coinFactoryKeeper = guardtestutil.NewMockCoinFactoryKeeper(ctrl)
	tokenKeeper := guardtestutil.NewMockTokenKeeper(ctrl)

	s.guardKeeper = *keeper.NewKeeper(
		cdc,
		storeService,
		logger,
		authority,
		nil,
		authzKeeper,
		s.nftKeeper,
	)
	keeper.SetCoinFactoryKeeper(&s.guardKeeper, s.coinFactoryKeeper)
	keeper.SetTokenKeeper(&s.guardKeeper, tokenKeeper)

	queryHelper := baseapp.NewQueryServerTestHelper(s.ctx, cdc.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, s.guardKeeper)
	queryClient := types.NewQueryClient(queryHelper)
	s.queryClient = queryClient

	s.defaultPrivileges = types.DefaultPrivileges
	s.rpKind = types.RequiredPrivilegesCoin
	s.lkIndex = []byte("factory/mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka/testcoin")

	s.msgServer = keeper.NewMsgServerImpl(&s.guardKeeper)

	s.params = types.NewParams(
		testutil.TestAdminAddress,
		testutil.TestAdminAddress,
		testutil.TestAccountPrivilegesGuardNftCollectionId,
		types.DefaultPrivileges,
		types.DefaultBaseDenom,
	)
	err := s.guardKeeper.SetParams(s.ctx, s.params)
	require.NoError(s.T(), err)
	s.testAdminAccount = "mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka"
	s.testAccount = testAccount
}
