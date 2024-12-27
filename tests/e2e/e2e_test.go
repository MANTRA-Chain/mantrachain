package e2e

import "fmt"

var (
	runBankTest                   = true
	runEncodeTest                 = true
	runEvidenceTest               = true
	runGovTest                    = true
	runIBCTest                    = false // TODO: enable after IBC test is fixed
	runSlashingTest               = true
	runStakingAndDistributionTest = true
	runVestingTest                = true
	runRestInterfacesTest         = true
	runRateLimitTest              = false // TODO: enable after IBC is fixed
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

	// TODO: add back when CCV is enabled
	// s.testSetBlocksPerEpoch()
	// s.ExpeditedProposalRejected()
	s.GovSoftwareUpgradeExpedited()
}

func (s *IntegrationTestSuite) TestIBC() {
	if !runIBCTest {
		s.T().Skip()
	}

	s.testIBCTokenTransfer()
	s.testMultihopIBCTokenTransfer()
	s.testFailedMultihopIBCTokenTransfer()
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

// func (s *IntegrationTestSuite) TestRateLimit() {
// 	if !runRateLimitTest {
// 		s.T().Skip()
// 	}
// 	s.testAddRateLimits()
// 	s.testIBCTransfer(true)
// 	s.testUpdateRateLimit()
// 	s.testIBCTransfer(false)
// 	s.testResetRateLimit()
// 	s.testRemoveRateLimit()
// }
