package keeper_test

import (
	utils "github.com/MANTRA-Finance/mantrachain/types"
	module "github.com/MANTRA-Finance/mantrachain/x/lpfarm/module"
	"github.com/MANTRA-Finance/mantrachain/x/lpfarm/types"
)

func (s *KeeperTestSuite) TestImportExportGenesis() {
	s.createPair("denom1", "denom2")
	s.createPool(1, utils.ParseCoins("1000_000000denom1,1000_000000denom2"))
	s.createPrivatePlan([]types.RewardAllocation{
		types.NewPairRewardAllocation(1, utils.ParseCoins("100_000000stake")),
	}, utils.ParseCoins("10000_000000stake"))

	farmerAddr := utils.TestAddress(0)
	s.farm(farmerAddr, utils.ParseCoin("1_000000pool1"))
	s.nextBlock()
	s.harvest(farmerAddr, "pool1")
	s.nextBlock()

	genState := module.ExportGenesis(s.ctx, s.keeper)
	bz := s.app.AppCodec().MustMarshalJSON(genState)

	s.SetupTest()
	var genState2 types.GenesisState
	s.app.AppCodec().MustUnmarshalJSON(bz, &genState2)
	module.InitGenesis(s.ctx, s.keeper, genState2)
	genState3 := module.ExportGenesis(s.ctx, s.keeper)
	s.Require().Equal(*genState, *genState3)
}

func (s *KeeperTestSuite) TestImportExportGenesisEmpty() {
	genState := module.ExportGenesis(s.ctx, s.keeper)

	var genState2 types.GenesisState
	bz := s.app.AppCodec().MustMarshalJSON(genState)
	s.app.AppCodec().MustUnmarshalJSON(bz, &genState2)
	module.InitGenesis(s.ctx, s.keeper, genState2)

	genState3 := module.ExportGenesis(s.ctx, s.keeper)
	s.Require().Equal(*genState, genState2)
	s.Require().Equal(genState2, *genState3)
}
