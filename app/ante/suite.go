package ante

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/client"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	transferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	channelkeeper "github.com/cosmos/ibc-go/v8/modules/core/04-channel/keeper"
	portkeeper "github.com/cosmos/ibc-go/v8/modules/core/05-port/keeper"
	feeabskeeper "github.com/osmosis-labs/fee-abstraction/v8/x/feeabs/keeper"
	feeabstestutil "github.com/osmosis-labs/fee-abstraction/v8/x/feeabs/testutil"
	feeabstypes "github.com/osmosis-labs/fee-abstraction/v8/x/feeabs/types"
	feemarketante "github.com/skip-mev/feemarket/x/feemarket/ante"
	feemarketmocks "github.com/skip-mev/feemarket/x/feemarket/ante/mocks"
	feemarkettypes "github.com/skip-mev/feemarket/x/feemarket/types"
	"github.com/stretchr/testify/require"
	ubermock "go.uber.org/mock/gomock"
)

// TestAccount represents an account used in the tests in x/auth/ante.
type TestAccount struct {
	acc  sdk.AccountI
	priv cryptotypes.PrivKey
}

// AnteTestSuite is a test suite to be used with ante handler tests.
type AnteTestSuite struct {
	ctx             sdk.Context
	clientCtx       client.Context
	txBuilder       client.TxBuilder
	accountKeeper   authkeeper.AccountKeeper
	bankKeeper      *feemarketmocks.BankKeeper
	feeGrantKeeper  *feeabstestutil.MockFeegrantKeeper
	stakingKeeper   *feeabstestutil.MockStakingKeeper
	feeabsKeeper    feeabskeeper.Keeper
	feemarketKeeper feemarketante.FeeMarketKeeper
	channelKeeper   *feeabstestutil.MockChannelKeeper
	portKeeper      *feeabstestutil.MockPortKeeper
	scopedKeeper    *feeabstestutil.MockScopedKeeper
	encCfg          moduletestutil.TestEncodingConfig
}

// SetupTest setups a new test, with new app, context, and anteHandler.
func SetupTestSuite(t *testing.T, isCheckTx bool) *AnteTestSuite {
	t.Helper()
	suite := &AnteTestSuite{}
	ctrl := ubermock.NewController(t)

	govAuthority := authtypes.NewModuleAddress("gov").String()

	// Setup mock keepers
	suite.bankKeeper = feemarketmocks.NewBankKeeper(t)
	suite.stakingKeeper = feeabstestutil.NewMockStakingKeeper(ctrl)
	suite.feeGrantKeeper = feeabstestutil.NewMockFeegrantKeeper(ctrl)
	suite.channelKeeper = feeabstestutil.NewMockChannelKeeper(ctrl)
	suite.portKeeper = feeabstestutil.NewMockPortKeeper(ctrl)
	suite.scopedKeeper = feeabstestutil.NewMockScopedKeeper(ctrl)

	// setup necessary params for Account Keeper
	key := storetypes.NewKVStoreKey(feeabstypes.StoreKey)
	authKey := storetypes.NewKVStoreKey(authtypes.StoreKey)
	subspace := paramtypes.NewSubspace(nil, nil, nil, nil, "feeabs")
	subspace = subspace.WithKeyTable(feeabstypes.ParamKeyTable())
	maccPerms := map[string][]string{
		"fee_collector":          nil,
		"mint":                   {"minter"},
		"bonded_tokens_pool":     {"burner", "staking"},
		"not_bonded_tokens_pool": {"burner", "staking"},
		"multiPerm":              {"burner", "minter", "staking"},
		"random":                 {"random"},
		"feeabs":                 nil,
	}

	// setup context for Account Keeper
	testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))
	testCtx.CMS.MountStoreWithDB(authKey, storetypes.StoreTypeIAVL, testCtx.DB)
	testCtx.CMS.MountStoreWithDB(storetypes.NewTransientStoreKey("transient_test2"), storetypes.StoreTypeTransient, testCtx.DB)
	err := testCtx.CMS.LoadLatestVersion()
	require.NoError(t, err)
	suite.ctx = testCtx.Ctx.WithIsCheckTx(isCheckTx).WithBlockHeight(1) // app.BaseApp.NewContext(isCheckTx, tmproto.Header{}).WithBlockHeight(1)

	suite.encCfg = moduletestutil.MakeTestEncodingConfig(auth.AppModuleBasic{})
	suite.encCfg.Amino.RegisterConcrete(&testdata.TestMsg{}, "testdata.TestMsg", nil)
	testdata.RegisterInterfaces(suite.encCfg.InterfaceRegistry)
	suite.accountKeeper = authkeeper.NewAccountKeeper(
		suite.encCfg.Codec, runtime.NewKVStoreService(authKey), authtypes.ProtoBaseAccount, maccPerms,
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()), sdk.Bech32MainPrefix, govAuthority,
	)
	suite.accountKeeper.SetModuleAccount(suite.ctx, authtypes.NewEmptyModuleAccount(feeabstypes.ModuleName))
	// Setup feeabs keeper
	suite.feeabsKeeper = feeabskeeper.NewKeeper(suite.encCfg.Codec, key, subspace, suite.stakingKeeper, suite.accountKeeper, keeper.BaseKeeper{}, transferkeeper.Keeper{}, channelkeeper.Keeper{}, &portkeeper.Keeper{}, capabilitykeeper.ScopedKeeper{}, govAuthority)
	suite.clientCtx = client.Context{}.
		WithTxConfig(suite.encCfg.TxConfig)
	require.NoError(t, err)

	// setup txBuilder
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	// setup feemarket
	feemarketParams := feemarkettypes.DefaultParams()
	feemarketParams.FeeDenom = "uom"
	//	suite.feemarketKeeper = feemarketkeeper.NewKeeper(suite.encCfg.Codec, key, suite.accountKeeper, &DenomResolverImpl{
	//		FeeabsKeeper:  suite.feeabsKeeper,
	//		StakingKeeper: suite.stakingKeeper,
	//	}, govAuthority)

	// err = suite.feemarketKeeper.SetState(suite.ctx, feemarkettypes.DefaultState())
	// require.NoError(t, err)
	return suite
}

func (suite *AnteTestSuite) CreateTestAccounts(numAccs int) []TestAccount {
	var accounts []TestAccount

	for i := 0; i < numAccs; i++ {
		priv, _, addr := testdata.KeyTestPubAddr()
		acc := suite.accountKeeper.NewAccountWithAddress(suite.ctx, addr)
		err := acc.SetAccountNumber(uint64(i + 100))
		if err != nil {
			panic(err)
		}
		suite.accountKeeper.SetAccount(suite.ctx, acc)
		accounts = append(accounts, TestAccount{acc, priv})
	}

	return accounts
}
