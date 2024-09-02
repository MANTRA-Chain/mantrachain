package keeper_test

import (
	module "github.com/MANTRA-Finance/mantrachain/x/coinfactory/module"
	"github.com/MANTRA-Finance/mantrachain/x/coinfactory/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func (suite *KeeperTestSuite) TestGenesis() {
	genesisState := types.GenesisState{
		FactoryDenoms: []types.GenesisDenom{
			{
				Denom: "factory/mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka/bitcoin",
				AuthorityMetadata: types.DenomAuthorityMetadata{
					Admin: "mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka",
				},
			},
			{
				Denom: "factory/mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka/diff-admin",
				AuthorityMetadata: types.DenomAuthorityMetadata{
					Admin: "mantra10d07y265gmmuvt4z0w9aw880jnsr700j3fep4f",
				},
			},
			{
				Denom: "factory/mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka/litecoin",
				AuthorityMetadata: types.DenomAuthorityMetadata{
					Admin: "mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka",
				},
			},
		},
	}

	suite.SetupTest()
	app := suite.app

	// Test both with bank denom metadata set, and not set.
	for i, denom := range genesisState.FactoryDenoms {
		// hacky, sets bank metadata to exist if i != 0, to cover both cases.
		if i != 0 {
			app.BankKeeper.SetDenomMetaData(suite.ctx, banktypes.Metadata{Base: denom.GetDenom()})
		}
	}

	module.InitGenesis(suite.ctx, app.CoinfactoryKeeper, genesisState)

	coinfactoryModuleAccount := app.AccountKeeper.GetModuleAccount(suite.ctx, types.ModuleName)
	suite.Require().NotNil(coinfactoryModuleAccount)

	exportedGenesis := module.ExportGenesis(suite.ctx, app.CoinfactoryKeeper)
	suite.Require().NotNil(exportedGenesis)
	suite.Require().Equal(genesisState, *exportedGenesis)
}
