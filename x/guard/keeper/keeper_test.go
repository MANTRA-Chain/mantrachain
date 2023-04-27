package keeper_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	tmtproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmttime "github.com/tendermint/tendermint/types/time"

	"github.com/MANTRA-Finance/mantrachain/x/guard"
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
	testutil.IBCConnectionTestSuite

	ctx               sdk.Context
	addrs             []sdk.AccAddress
	guardKeeper       keeper.Keeper
	bankKeeper        *guardtestutil.MockBankKeeper
	encCfg            testutil.TestEncodingConfig
	queryClient       types.QueryClient
	msgServer         types.MsgServer
	defaultPrivileges []byte
	rpKind            types.RequiredPrivilegesKind
	lkKind            types.LockedKind
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
	s.ctx = testCtx.Ctx.WithBlockHeader(tmtproto.Header{Time: tmttime.Now()})
	s.encCfg = testutil.MakeTestEncodingConfig(guard.AppModuleBasic{})

	s.addrs = testutil.CreateIncrementalAccounts(3)

	ctrl := gomock.NewController(s.T())
	accountKeeper := guardtestutil.NewMockAccountKeeper(ctrl)
	s.bankKeeper = guardtestutil.NewMockBankKeeper(ctrl)
	authzKeeper := guardtestutil.NewMockAuthzKeeper(ctrl)
	tokenKeeper := guardtestutil.NewMockTokenKeeper(ctrl)
	nftKeeper := guardtestutil.NewMockNFTKeeper(ctrl)
	coinFactoryKeeper := guardtestutil.NewMockCoinFactoryKeeper(ctrl)

	s.guardKeeper = keeper.NewKeeper(
		s.encCfg.Codec,
		key,
		paramstypes.NewSubspace(s.encCfg.Codec, s.encCfg.Amino, key, tkey, "testsubspace").WithKeyTable(types.ParamKeyTable()),
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
	s.lkKind = types.LockedCoin
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
