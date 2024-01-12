package keeper_test

import (
	"testing"

	cbproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cbttime "github.com/cometbft/cometbft/types/time"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/AumegaChain/aumega/app/params"
	"github.com/AumegaChain/aumega/x/guard/keeper"
	guardtestutil "github.com/AumegaChain/aumega/x/guard/testutil"
	"github.com/AumegaChain/aumega/x/guard/types"

	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/AumegaChain/aumega/testutil"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	testAccount    = "cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw"
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
	encCfg            params.EncodingConfig
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
	key := sdk.NewKVStoreKey(types.StoreKey)
	tkey := sdk.NewTransientStoreKey("transient_test")
	moduleAccAddr := make(map[string]bool)

	testCtx := testutil.DefaultContextWithDB(s.T(), key, tkey)
	s.ctx = testCtx.Ctx.WithBlockHeader(cbproto.Header{Time: cbttime.Now()})
	s.encCfg = testutil.MakeTestEncodingConfig()

	s.addrs = testutil.CreateIncrementalAccounts(3)

	ctrl := gomock.NewController(s.T())
	accountKeeper := guardtestutil.NewMockAccountKeeper(ctrl)
	s.bankKeeper = guardtestutil.NewMockBankKeeper(ctrl)
	authzKeeper := guardtestutil.NewMockAuthzKeeper(ctrl)
	tokenKeeper := guardtestutil.NewMockTokenKeeper(ctrl)
	s.nftKeeper = guardtestutil.NewMockNFTKeeper(ctrl)
	s.coinFactoryKeeper = guardtestutil.NewMockCoinFactoryKeeper(ctrl)

	s.guardKeeper = keeper.NewKeeper(
		s.encCfg.Marshaler,
		key,
		paramstypes.NewSubspace(s.encCfg.Marshaler, s.encCfg.Amino, key, tkey, "testsubspace").WithKeyTable(types.ParamKeyTable()),
		moduleAccAddr,
		nil,
		accountKeeper,
		s.bankKeeper,
		authzKeeper,
		tokenKeeper,
		s.nftKeeper,
		s.coinFactoryKeeper,
	)

	queryHelper := baseapp.NewQueryServerTestHelper(s.ctx, s.encCfg.InterfaceRegistry)
	types.RegisterQueryServer(queryHelper, s.guardKeeper)
	queryClient := types.NewQueryClient(queryHelper)
	s.queryClient = queryClient

	s.defaultPrivileges = types.DefaultPrivileges
	s.rpKind = types.RequiredPrivilegesCoin
	s.lkIndex = []byte("factory/cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw/testcoin")

	s.msgServer = keeper.NewMsgServerImpl(&s.guardKeeper)

	s.params = types.NewParams(
		testutil.TestAdminAddress,
		testutil.TestAdminAddress,
		testutil.TestAccountPrivilegesGuardNftCollectionId,
		types.DefaultPrivileges,
		types.DefaultBaseDenom,
	)
	s.guardKeeper.SetParams(s.ctx, s.params)
	s.testAdminAccount = "cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw"
	s.testAccount = testAccount
}
