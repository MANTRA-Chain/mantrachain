package keeper_test

import (
	"cosmossdk.io/math"
	marketmaker "github.com/MANTRA-Finance/mantrachain/x/marketmaker/module"
	"github.com/MANTRA-Finance/mantrachain/x/marketmaker/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/stretchr/testify/suite"
)

func (suite *KeeperTestSuite) TestDefaultGenesis() {
	genState := *types.DefaultGenesis()

	marketmaker.InitGenesis(suite.ctx, suite.keeper, genState)
	got := marketmaker.ExportGenesis(suite.ctx, suite.keeper)
	suite.Require().Equal(genState, *got)
}

func (suite *KeeperTestSuite) TestImportExportGenesisEmpty() {
	k, ctx := suite.keeper, suite.ctx
	genState := marketmaker.ExportGenesis(ctx, k)

	var genState2 types.GenesisState
	bz := suite.app.AppCodec().MustMarshalJSON(genState)
	suite.app.AppCodec().MustUnmarshalJSON(bz, &genState2)
	marketmaker.InitGenesis(ctx, k, genState2)

	genState3 := marketmaker.ExportGenesis(ctx, k)
	suite.Require().Equal(*genState, genState2, *genState3)
}

func (suite *KeeperTestSuite) TestInitGenesis() {
	ctx := suite.ctx
	k := suite.keeper
	mmAddr := suite.addrs[0]
	mmAddr2 := suite.addrs[1]

	// set incentive budget
	params := k.GetParams(ctx)
	params.IncentiveBudgetAddress = suite.addrs[5].String()
	err := k.SetParams(ctx, params)
	suite.Require().NoError(err)

	// apply market maker
	err = k.ApplyMarketMaker(ctx, mmAddr, []uint64{1, 2, 3, 4, 5, 6})
	suite.NoError(err)
	err = k.ApplyMarketMaker(ctx, mmAddr2, []uint64{2, 3, 4, 5, 6, 7})
	suite.NoError(err)

	// include market maker
	proposal := types.NewMarketMakerProposal("title", "description",
		[]types.MarketMakerHandle{
			{Address: mmAddr.String(), PairId: 1},
			{Address: mmAddr2.String(), PairId: 3},
		},
		nil, nil, nil)
	suite.handleProposal(proposal)

	// distribute incentive
	incentiveAmount := math.NewInt(500000000)
	incentiveCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, incentiveAmount))
	proposal = types.NewMarketMakerProposal("title", "description", nil, nil, nil,
		[]types.IncentiveDistribution{
			{
				Address: mmAddr.String(),
				PairId:  1,
				Amount:  incentiveCoins,
			},
			{
				Address: mmAddr2.String(),
				PairId:  3,
				Amount:  incentiveCoins,
			},
		})
	suite.handleProposal(proposal)

	mms := k.GetAllMarketMakers(ctx)
	suite.Require().Len(mms, 12)

	incentives := k.GetAllIncentives(ctx)
	suite.Require().Len(incentives, 2)

	var genState *types.GenesisState
	suite.Require().NotPanics(func() {
		genState = marketmaker.ExportGenesis(suite.ctx, suite.keeper)
	})

	err = genState.Validate()
	suite.Require().NoError(err)

	suite.Require().NotPanics(func() {
		marketmaker.InitGenesis(suite.ctx, suite.keeper, *genState)
	})
	suite.Require().Equal(genState, marketmaker.ExportGenesis(suite.ctx, k))

	mms = suite.keeper.GetAllMarketMakers(ctx)
	suite.Require().Len(mms, 12)

	incentives = suite.keeper.GetAllIncentives(ctx)
	suite.Require().Len(incentives, 2)
}
