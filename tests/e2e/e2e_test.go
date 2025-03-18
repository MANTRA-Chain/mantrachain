package e2e

import "fmt"

var (
	runBankTest                   = true
	runEncodeTest                 = true
	runEvidenceTest               = true
	runGovTest                    = true
	runIBCTest                    = true
	runSlashingTest               = true
	runStakingAndDistributionTest = true
	runVestingTest                = true
	runRestInterfacesTest         = true
	runRateLimitTest              = true
	runTokenfactoryTest           = true
	runSanctionTest               = true
)

func (s *IntegrationTestSuite) TestRestInterfaces() {
	if !runRestInterfacesTest {
		s.T().Skip()
	}
	s.testRestInterfaces()
}

func (s *IntegrationTestSuite) TestBank() {
	if !runBankTest {
		s.T().Skip()
	}
	s.testBankTokenTransfer()
}

func (s *IntegrationTestSuite) TestEncode() {
	if !runEncodeTest {
		s.T().Skip()
	}
	s.testEncode()
	s.testDecode()
}

func (s *IntegrationTestSuite) TestEvidence() {
	if !runEvidenceTest {
		s.T().Skip()
	}
	s.testEvidence()
}

func (s *IntegrationTestSuite) TestGov() {
	if !runGovTest {
		s.T().Skip()
	}

	s.GovCancelSoftwareUpgrade()
	s.GovCommunityPoolSpend()

	// TODO: uncomment in future if CCV is enabled
	// s.testSetBlocksPerEpoch()
	// s.ExpeditedProposalRejected()
	s.GovSoftwareUpgradeExpedited()
}

func (s *IntegrationTestSuite) TestIBC() {
	if !runIBCTest {
		s.T().Skip()
	}

	s.testIBCTokenTransfer()
	// TODO: uncomment in future if we add PFM
	// s.testMultihopIBCTokenTransfer()
	// s.testFailedMultihopIBCTokenTransfer()
	s.testICARegisterAccountAndSendTx()
}

func (s *IntegrationTestSuite) TestSlashing() {
	if !runSlashingTest {
		s.T().Skip()
	}
	chainAPI := fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))
	s.testSlashing(chainAPI)
}

// todo add fee test with wrong denom order
func (s *IntegrationTestSuite) TestStakingAndDistribution() {
	if !runStakingAndDistributionTest {
		s.T().Skip()
	}
	s.testStaking()
	s.testDistribution()
}

func (s *IntegrationTestSuite) TestVesting() {
	if !runVestingTest {
		s.T().Skip()
	}
	chainAAPI := fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))
	s.testDelayedVestingAccount(chainAAPI)
	s.testContinuousVestingAccount(chainAAPI)
	s.testPeriodicVestingAccount(chainAAPI)
}

func (s *IntegrationTestSuite) TestRateLimit() {
	if !runRateLimitTest {
		s.T().Skip()
	}
	s.testAddRateLimits()
	s.testIBCTransfer(true)
	s.testUpdateRateLimit()
	s.testIBCTransfer(false)
	s.testResetRateLimit()
	s.testRemoveRateLimit()
}

func (s *IntegrationTestSuite) TestTokenfactory() {
	if !runTokenfactoryTest {
		s.T().Skip()
	}
	s.testTokenfactoryCreate()
	s.testTokenfactoryMint()
	s.testTokenfactoryBurn()
}

func (s *IntegrationTestSuite) TestSanction() {
	if !runSanctionTest {
		s.T().Skip()
	}
	s.testAddToBlacklist()
	s.testRemoveFromBlacklist()
}
