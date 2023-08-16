package keeper_test

import (
	"testing"

	cbproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cbttime "github.com/cometbft/cometbft/types/time"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/MANTRA-Finance/mantrachain/app/params"
	"github.com/MANTRA-Finance/mantrachain/x/guard/keeper"
	guardtestutil "github.com/MANTRA-Finance/mantrachain/x/guard/testutil"
	"github.com/MANTRA-Finance/mantrachain/x/guard/types"

	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/MANTRA-Finance/mantrachain/testutil"

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
	encCfg            params.EncodingConfig
	queryClient       types.QueryClient
	msgServer         types.MsgServer
	defaultPrivileges []byte
	rpKind            types.RequiredPrivilegesKind
	lkIndex           []byte
	params            types.Params
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
	nftKeeper := guardtestutil.NewMockNFTKeeper(ctrl)
	coinFactoryKeeper := guardtestutil.NewMockCoinFactoryKeeper(ctrl)

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
		nftKeeper,
		coinFactoryKeeper,
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
	)
	s.guardKeeper.SetParams(s.ctx, s.params)
}
