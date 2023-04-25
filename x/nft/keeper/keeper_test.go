package keeper_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmttime "github.com/tendermint/tendermint/types/time"

	"github.com/MANTRA-Finance/mantrachain/x/nft"
	"github.com/MANTRA-Finance/mantrachain/x/nft/keeper"
	nfttestutil "github.com/MANTRA-Finance/mantrachain/x/nft/testutil"
	"github.com/MANTRA-Finance/mantrachain/x/nft/types"

	"github.com/MANTRA-Finance/mantrachain/testutil"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	testClassID          = "kitty"
	testClassName        = "Crypto Kitty"
	testClassSymbol      = "kitty"
	testClassDescription = "Crypto Kitty"
	testClassURI         = "class uri"
	testClassURIHash     = "ae702cefd6b6a65fe2f991ad6d9969ed"
	testID               = "kitty1"
	testURI              = "kitty uri"
	testURIHash          = "229bfd3c1b431c14a526497873897108"
)

type TestSuite struct {
	testutil.IBCConnectionTestSuite

	ctx           sdk.Context
	addrs         []sdk.AccAddress
	queryClient   types.QueryClient
	nftKeeper     keeper.Keeper
	accountKeeper *nfttestutil.MockAccountKeeper

	encCfg testutil.TestEncodingConfig
}

func (s *TestSuite) SetupTest() {
	// suite setup
	s.addrs = testutil.CreateIncrementalAccounts(3)
	s.encCfg = testutil.MakeTestEncodingConfig(nft.AppModuleBasic{})

	key := sdk.NewKVStoreKey(types.StoreKey)
	testCtx := testutil.DefaultContextWithDB(s.T(), key, sdk.NewTransientStoreKey("transient_test"))
	ctx := testCtx.Ctx.WithBlockHeader(tmproto.Header{Time: tmttime.Now()})

	// gomock initializations
	ctrl := gomock.NewController(s.T())
	accountKeeper := nfttestutil.NewMockAccountKeeper(ctrl)
	bankKeeper := nfttestutil.NewMockBankKeeper(ctrl)
	accountKeeper.EXPECT().GetModuleAddress("nft").Return(s.addrs[0]).AnyTimes()

	s.accountKeeper = accountKeeper

	nftKeeper := keeper.NewKeeper(key, s.encCfg.Codec, accountKeeper, bankKeeper)
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, s.encCfg.InterfaceRegistry)
	types.RegisterQueryServer(queryHelper, nftKeeper)

	s.nftKeeper = nftKeeper
	s.queryClient = types.NewQueryClient(queryHelper)
	s.ctx = ctx
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TestSaveClass() {
	except := types.Class{
		Id:          testClassID,
		Name:        testClassName,
		Symbol:      testClassSymbol,
		Description: testClassDescription,
		Uri:         testClassURI,
		UriHash:     testClassURIHash,
	}
	err := s.nftKeeper.SaveClass(s.ctx, except)
	s.Require().NoError(err)

	actual, has := s.nftKeeper.GetClass(s.ctx, testClassID)
	s.Require().True(has)
	s.Require().EqualValues(except, actual)

	classes := s.nftKeeper.GetClasses(s.ctx)
	s.Require().EqualValues([]*types.Class{&except}, classes)
}

func (s *TestSuite) TestUpdateClass() {
	class := types.Class{
		Id:          testClassID,
		Name:        testClassName,
		Symbol:      testClassSymbol,
		Description: testClassDescription,
		Uri:         testClassURI,
		UriHash:     testClassURIHash,
	}
	err := s.nftKeeper.SaveClass(s.ctx, class)
	s.Require().NoError(err)

	noExistClass := types.Class{
		Id:          "kitty1",
		Name:        testClassName,
		Symbol:      testClassSymbol,
		Description: testClassDescription,
		Uri:         testClassURI,
		UriHash:     testClassURIHash,
	}

	err = s.nftKeeper.UpdateClass(s.ctx, noExistClass)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "nft class does not exist")

	except := types.Class{
		Id:          testClassID,
		Name:        "My crypto Kitty",
		Symbol:      testClassSymbol,
		Description: testClassDescription,
		Uri:         testClassURI,
		UriHash:     testClassURIHash,
	}

	err = s.nftKeeper.UpdateClass(s.ctx, except)
	s.Require().NoError(err)

	actual, has := s.nftKeeper.GetClass(s.ctx, testClassID)
	s.Require().True(has)
	s.Require().EqualValues(except, actual)
}

func (s *TestSuite) TestMint() {
	class := types.Class{
		Id:          testClassID,
		Name:        testClassName,
		Symbol:      testClassSymbol,
		Description: testClassDescription,
		Uri:         testClassURI,
		UriHash:     testClassURIHash,
	}
	err := s.nftKeeper.SaveClass(s.ctx, class)
	s.Require().NoError(err)

	expNFT := types.NFT{
		ClassId: testClassID,
		Id:      testID,
		Uri:     testURI,
	}
	err = s.nftKeeper.Mint(s.ctx, expNFT, s.addrs[0])
	s.Require().NoError(err)

	// test GetNFT
	actNFT, has := s.nftKeeper.GetNFT(s.ctx, testClassID, testID)
	s.Require().True(has)
	s.Require().EqualValues(expNFT, actNFT)

	// test GetOwner
	owner := s.nftKeeper.GetOwner(s.ctx, testClassID, testID)
	s.Require().True(s.addrs[0].Equals(owner))

	// test GetNFTsOfClass
	actNFTs := s.nftKeeper.GetNFTsOfClass(s.ctx, testClassID)
	s.Require().EqualValues([]types.NFT{expNFT}, actNFTs)

	// test GetNFTsOfClassByOwner
	actNFTs = s.nftKeeper.GetNFTsOfClassByOwner(s.ctx, testClassID, s.addrs[0])
	s.Require().EqualValues([]types.NFT{expNFT}, actNFTs)

	// test GetBalance
	balance := s.nftKeeper.GetBalance(s.ctx, testClassID, s.addrs[0])
	s.Require().EqualValues(uint64(1), balance)

	// test GetTotalSupply
	supply := s.nftKeeper.GetTotalSupply(s.ctx, testClassID)
	s.Require().EqualValues(uint64(1), supply)

	expNFT2 := types.NFT{
		ClassId: testClassID,
		Id:      testID + "2",
		Uri:     testURI + "2",
	}
	err = s.nftKeeper.Mint(s.ctx, expNFT2, s.addrs[0])
	s.Require().NoError(err)

	// test GetNFTsOfClassByOwner
	actNFTs = s.nftKeeper.GetNFTsOfClassByOwner(s.ctx, testClassID, s.addrs[0])
	s.Require().EqualValues([]types.NFT{expNFT, expNFT2}, actNFTs)

	// test GetBalance
	balance = s.nftKeeper.GetBalance(s.ctx, testClassID, s.addrs[0])
	s.Require().EqualValues(uint64(2), balance)
}

func (s *TestSuite) TestBurn() {
	except := types.Class{
		Id:          testClassID,
		Name:        testClassName,
		Symbol:      testClassSymbol,
		Description: testClassDescription,
		Uri:         testClassURI,
		UriHash:     testClassURIHash,
	}
	err := s.nftKeeper.SaveClass(s.ctx, except)
	s.Require().NoError(err)

	expNFT := types.NFT{
		ClassId: testClassID,
		Id:      testID,
		Uri:     testURI,
	}
	err = s.nftKeeper.Mint(s.ctx, expNFT, s.addrs[0])
	s.Require().NoError(err)

	err = s.nftKeeper.Burn(s.ctx, testClassID, testID)
	s.Require().NoError(err)

	// test GetNFT
	_, has := s.nftKeeper.GetNFT(s.ctx, testClassID, testID)
	s.Require().False(has)

	// test GetOwner
	owner := s.nftKeeper.GetOwner(s.ctx, testClassID, testID)
	s.Require().Nil(owner)

	// test GetNFTsOfClass
	actNFTs := s.nftKeeper.GetNFTsOfClass(s.ctx, testClassID)
	s.Require().Empty(actNFTs)

	// test GetNFTsOfClassByOwner
	actNFTs = s.nftKeeper.GetNFTsOfClassByOwner(s.ctx, testClassID, s.addrs[0])
	s.Require().Empty(actNFTs)

	// test GetBalance
	balance := s.nftKeeper.GetBalance(s.ctx, testClassID, s.addrs[0])
	s.Require().EqualValues(uint64(0), balance)

	// test GetTotalSupply
	supply := s.nftKeeper.GetTotalSupply(s.ctx, testClassID)
	s.Require().EqualValues(uint64(0), supply)
}

func (s *TestSuite) TestUpdate() {
	class := types.Class{
		Id:          testClassID,
		Name:        testClassName,
		Symbol:      testClassSymbol,
		Description: testClassDescription,
		Uri:         testClassURI,
		UriHash:     testClassURIHash,
	}
	err := s.nftKeeper.SaveClass(s.ctx, class)
	s.Require().NoError(err)

	myNFT := types.NFT{
		ClassId: testClassID,
		Id:      testID,
		Uri:     testURI,
	}
	err = s.nftKeeper.Mint(s.ctx, myNFT, s.addrs[0])
	s.Require().NoError(err)

	expNFT := types.NFT{
		ClassId: testClassID,
		Id:      testID,
		Uri:     "updated",
	}

	err = s.nftKeeper.Update(s.ctx, expNFT)
	s.Require().NoError(err)

	// test GetNFT
	actNFT, has := s.nftKeeper.GetNFT(s.ctx, testClassID, testID)
	s.Require().True(has)
	s.Require().EqualValues(expNFT, actNFT)
}

func (s *TestSuite) TestTransfer() {
	class := types.Class{
		Id:          testClassID,
		Name:        testClassName,
		Symbol:      testClassSymbol,
		Description: testClassDescription,
		Uri:         testClassURI,
		UriHash:     testClassURIHash,
	}
	err := s.nftKeeper.SaveClass(s.ctx, class)
	s.Require().NoError(err)

	expNFT := types.NFT{
		ClassId: testClassID,
		Id:      testID,
		Uri:     testURI,
	}
	err = s.nftKeeper.Mint(s.ctx, expNFT, s.addrs[0])
	s.Require().NoError(err)

	// valid owner
	err = s.nftKeeper.Transfer(s.ctx, testClassID, testID, s.addrs[1])
	s.Require().NoError(err)

	// test GetOwner
	owner := s.nftKeeper.GetOwner(s.ctx, testClassID, testID)
	s.Require().Equal(s.addrs[1], owner)

	balanceAddr0 := s.nftKeeper.GetBalance(s.ctx, testClassID, s.addrs[0])
	s.Require().EqualValues(uint64(0), balanceAddr0)

	balanceAddr1 := s.nftKeeper.GetBalance(s.ctx, testClassID, s.addrs[1])
	s.Require().EqualValues(uint64(1), balanceAddr1)

	// test GetNFTsOfClassByOwner
	actNFTs := s.nftKeeper.GetNFTsOfClassByOwner(s.ctx, testClassID, s.addrs[1])
	s.Require().EqualValues([]types.NFT{expNFT}, actNFTs)
}

func (s *TestSuite) TestExportGenesis() {
	class := types.Class{
		Id:          testClassID,
		Name:        testClassName,
		Symbol:      testClassSymbol,
		Description: testClassDescription,
		Uri:         testClassURI,
		UriHash:     testClassURIHash,
	}
	err := s.nftKeeper.SaveClass(s.ctx, class)
	s.Require().NoError(err)

	expNFT := types.NFT{
		ClassId: testClassID,
		Id:      testID,
		Uri:     testURI,
	}
	err = s.nftKeeper.Mint(s.ctx, expNFT, s.addrs[0])
	s.Require().NoError(err)

	expGenesis := &types.GenesisState{
		Classes: []*types.Class{&class},
		Entries: []*types.Entry{{
			Owner: s.addrs[0].String(),
			Nfts:  []*types.NFT{&expNFT},
		}},
	}
	genesis := s.nftKeeper.ExportGenesis(s.ctx)
	s.Require().Equal(expGenesis, genesis)
}

func (s *TestSuite) TestInitGenesis() {
	expClass := types.Class{
		Id:          testClassID,
		Name:        testClassName,
		Symbol:      testClassSymbol,
		Description: testClassDescription,
		Uri:         testClassURI,
		UriHash:     testClassURIHash,
	}
	expNFT := types.NFT{
		ClassId: testClassID,
		Id:      testID,
		Uri:     testURI,
	}
	expGenesis := &types.GenesisState{
		Classes: []*types.Class{&expClass},
		Entries: []*types.Entry{{
			Owner: s.addrs[0].String(),
			Nfts:  []*types.NFT{&expNFT},
		}},
	}
	s.nftKeeper.InitGenesis(s.ctx, expGenesis)

	actual, has := s.nftKeeper.GetClass(s.ctx, testClassID)
	s.Require().True(has)
	s.Require().EqualValues(expClass, actual)

	// test GetNFT
	actNFT, has := s.nftKeeper.GetNFT(s.ctx, testClassID, testID)
	s.Require().True(has)
	s.Require().EqualValues(expNFT, actNFT)
}
