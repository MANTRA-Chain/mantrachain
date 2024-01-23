package keeper_test

import (
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/MANTRA-Finance/mantrachain/x/coinfactory/types"
)

func (suite *KeeperTestSuite) TestGenesis() {
	genesisState := types.GenesisState{
		FactoryDenoms: []types.GenesisDenom{
			{
				Denom: "factory/cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw/bitcoin",
				AuthorityMetadata: types.DenomAuthorityMetadata{
					Admin: "cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw",
				},
			},
			{
				Denom: "factory/cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw/diff-admin",
				AuthorityMetadata: types.DenomAuthorityMetadata{
					Admin: "cosmos15ejrsrfts5jfd8vekdje4t3t56nvflry92uegz",
				},
			},
			{
				Denom: "factory/cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw/litecoin",
				AuthorityMetadata: types.DenomAuthorityMetadata{
					Admin: "cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw",
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

	app.CoinFactoryKeeper.InitGenesis(suite.ctx, genesisState)

	coinfactoryModuleAccount := app.AccountKeeper.GetAccount(suite.ctx, app.AccountKeeper.GetModuleAddress(types.ModuleName))
	suite.Require().NotNil(coinfactoryModuleAccount)

	exportedGenesis := app.CoinFactoryKeeper.ExportGenesis(suite.ctx)
	suite.Require().NotNil(exportedGenesis)
	suite.Require().Equal(genesisState, *exportedGenesis)
}
